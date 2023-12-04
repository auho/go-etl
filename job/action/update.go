package action

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/mode"
	"github.com/auho/go-etl/v2/tool/slices"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/storage/database"
	"github.com/auho/go-toolkit/flow/storage/database/destination"
)

var _ Actor = (*Update)(nil)

type UpdateConfig struct {
	NotTruncate bool // for update and transfer
	BatchSize   int  // for update and transfer
	Concurrency int  // for update and transfer
}

func (uc *UpdateConfig) check() {
	if uc.BatchSize <= 0 {
		uc.BatchSize = batchSize
	}

	if uc.Concurrency <= 0 {
		uc.Concurrency = runtime.NumCPU()
	}
}

func WithUpdateConfig(cc UpdateConfig) func(update *Update) {
	return func(c *Update) {
		c.config = cc
	}
}

type Update struct {
	TargetAction

	source     job.Source
	modes      []mode.UpdateModer
	isTransfer bool

	config  UpdateConfig
	dst     *destination.Destination[storage.MapEntry]
	dstLine int
}

func NewUpdateAndTransfer(source job.Source, target job.Target, modes []mode.UpdateModer, opts ...func(*Update)) *Update {
	u := NewUpdate(source, modes, opts...)
	u.target = target
	u.isTransfer = true

	return u
}

func NewUpdate(source job.Source, modes []mode.UpdateModer, opts ...func(*Update)) *Update {
	u := &Update{}
	u.source = source
	u.modes = modes

	for _, opt := range opts {
		opt(u)
	}

	u.Init()

	u.config.check()

	return u
}

func (u *Update) GetFields() []string {
	fields := make([]string, 0)
	fields = append(fields, u.source.GetIdName())

	for _, m := range u.modes {
		fields = append(fields, m.GetFields()...)
	}

	if u.isTransfer {
		columns, err := u.target.GetDB().GetTableColumns(u.target.TableName())
		if err != nil {
			panic(err)
		}

		fields = append(fields, columns...)
	}

	fields = slices.SliceDropDuplicates(fields)

	return fields
}

func (u *Update) Title() string {
	s := make([]string, 0)
	for _, m := range u.modes {
		s = append(s, m.GetTitle())
	}

	return fmt.Sprintf("Update[%s] {%s}", u.source.TableName(), strings.Join(s, ", "))
}

func (u *Update) Prepare() error {
	var err error
	for _, m := range u.modes {
		err = m.Prepare()
		if err != nil {
			return fmt.Errorf("update action prepare error; %w", err)
		}
	}

	if u.isTransfer {
		u.dst, err = newUpdateToDB(u, u.target)
		if err != nil {
			return fmt.Errorf("newUpdateToDB dst error; %w", err)
		}
	}

	return nil
}

func (u *Update) PreDo() error {
	if u.isTransfer {
		err := u.dst.Accept()
		if err != nil {
			return fmt.Errorf("dst accept error;%w", err)
		}

		u.dstLine = u.AddState(fmt.Sprintf("data: %s", strings.Join(u.dst.State(), "\n")))
	}

	return nil
}

func (u *Update) Do(item map[string]any) ([]map[string]any, bool) {
	_does := make(map[string]any)
	for _, m := range u.modes {
		_do := m.Do(item)
		for k, v := range _do {
			_does[k] = v
		}
	}

	if len(_does) <= 0 && u.isTransfer == false {
		return nil, false
	}

	var newItem map[string]any
	if u.isTransfer {
		newItem = item
	} else {
		newItem = make(map[string]any)
		newItem[u.source.GetIdName()] = item[u.source.GetIdName()]
	}

	for k, v := range _does {
		newItem[k] = v
	}

	return []map[string]any{newItem}, true
}

func (u *Update) PostBatchDo(items []map[string]any) {
	var err error
	if u.isTransfer {
		u.dst.Receive(items)
	} else {
		err = u.source.GetDB().BulkUpdateFromSliceMapById(u.source.TableName(), u.source.GetIdName(), items)
	}

	if err != nil {
		panic(err)
	}
}

func (u *Update) Blink() {
	if u.isTransfer {
		u.SetState(u.dstLine, fmt.Sprintf("data: %s", strings.Join(u.dst.State(), "\n")))
	}
}

func (u *Update) PostDo() error {
	if u.isTransfer {
		u.dst.Done()
		u.dst.Finish()
	}

	return nil
}

func (u *Update) Close() error {
	for _, m := range u.modes {
		err := m.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

var _ destination.Destinationer[storage.MapEntry] = (*updateToDB)(nil)

type updateToDB struct{}

func (i *updateToDB) Exec(d *destination.Destination[storage.MapEntry], items storage.MapEntries) error {
	return d.DB().BulkInsertFromSliceMap(d.TableName(), items, int(d.PageSize()))
}

func newUpdateToDB(u *Update, target job.Target) (*destination.Destination[storage.MapEntry], error) {
	return destination.NewDestination[storage.MapEntry](&destination.Config{
		IsTruncate:  !u.config.NotTruncate,
		Concurrency: u.config.Concurrency,
		PageSize:    int64(u.config.BatchSize),
		TableName:   target.TableName(),
	}, &updateToDB{}, func() (*database.DB, error) {
		return database.NewFromSimpleDb(target.GetDB()), nil
	})
}
