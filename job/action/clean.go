package action

import (
	"fmt"
	"maps"
	"runtime"
	"strings"

	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/mode"
	"github.com/auho/go-etl/v2/tool/slices"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/storage/database"
	"github.com/auho/go-toolkit/flow/storage/database/destination"
)

var _ Actor = (*Clean)(nil)

type CleanConfig struct {
	NotTruncate  bool
	AddExtraTags bool // tags to deleted data
	BatchSize    int
	Concurrency  int
	Keys         []string // source columns name，priority of use this keys
}

func (cc *CleanConfig) check() {
	if cc.BatchSize <= 0 {
		cc.BatchSize = batchSize
	}

	if cc.Concurrency <= 0 {
		cc.Concurrency = runtime.NumCPU()
	}
}

func WithCleanConfig(cc CleanConfig) func(*Clean) {
	return func(c *Clean) {
		c.config = cc
	}
}

// Clean
// filter
type Clean struct {
	action

	cleanTarget job.CleanResource
	keys        []string
	modes       []mode.UpdateModer

	config CleanConfig

	dataDst    *destination.Destination[storage.MapEntry]
	deletedDst *destination.Destination[storage.MapEntry]

	dataDstLine    int
	deletedDstLine int
}

func NewClean(cr job.CleanResource, modes []mode.UpdateModer, opts ...func(*Clean)) *Clean {
	c := &Clean{}
	c.cleanTarget = cr
	c.modes = modes

	for _, opt := range opts {
		opt(c)
	}

	c.Init()

	c.config.check()

	return c
}

func (c *Clean) GetFields() []string {
	if len(c.config.Keys) > 0 {
		c.keys = append(c.keys, c.config.Keys...)
		for _, m := range c.modes {
			c.keys = append(c.keys, m.GetFields()...)
		}

		c.keys = slices.SliceDropDuplicates(c.keys)

	} else {
		var err error
		c.keys, err = c.cleanTarget.DeletedTarget().GetDB().GetTableColumns(c.cleanTarget.DeletedTarget().TableName())
		if err != nil {
			panic(fmt.Errorf("GetTableColumns error; %w", err))
		}
	}

	return c.keys
}

func (c *Clean) Title() string {
	s := make([]string, 0)
	for _, m := range c.modes {
		s = append(s, m.GetTitle())
	}

	return fmt.Sprintf("Clean[%s, %s] {%s}",
		c.cleanTarget.DataTarget().TableName(),
		c.cleanTarget.DeletedTarget().TableName(),
		strings.Join(s, ", "),
	)
}

func (c *Clean) Prepare() error {
	var err error
	for _, m := range c.modes {
		err = m.Prepare()
		if err != nil {
			return fmt.Errorf("clean action prepare error; %w", err)
		}
	}

	c.dataDst, err = newInsertToDB(c, c.cleanTarget.DataTarget())
	if err != nil {
		return fmt.Errorf("newInsertToDB data error; %w", err)
	}

	c.deletedDst, err = newInsertToDB(c, c.cleanTarget.DeletedTarget())
	if err != nil {
		return fmt.Errorf("newInsertToDB deleted error; %w", err)
	}

	return nil
}

func (c *Clean) PreDo() error {
	err := c.dataDst.Accept()
	if err != nil {
		return fmt.Errorf("data accept error;%w", err)
	}

	err = c.deletedDst.Accept()
	if err != nil {
		return fmt.Errorf("deleted accept error;%w", err)
	}

	c.dataDstLine = c.AddState(fmt.Sprintf("data: %s", strings.Join(c.dataDst.State(), "\n")))
	c.deletedDstLine = c.AddState(fmt.Sprintf("deleted: %s", strings.Join(c.deletedDst.State(), "\n")))

	return nil
}

func (c *Clean) Do(item map[string]any) ([]map[string]any, bool) {
	_needDeleted := false
	for _, m := range c.modes {
		_res := m.Do(item)
		if len(_res) > 0 {
			_needDeleted = true

			if c.config.AddExtraTags {
				maps.Copy(item, _res)
			}

			break
		}
	}

	if _needDeleted {
		c.deletedDst.Receive([]map[string]any{item})
	} else {
		c.dataDst.Receive([]map[string]any{item})
	}

	return nil, true
}

func (c *Clean) PostBatchDo(items []map[string]any) {}

func (c *Clean) Blink() {
	c.SetState(c.dataDstLine, fmt.Sprintf("data: %s", strings.Join(c.dataDst.State(), "\n")))
	c.SetState(c.deletedDstLine, fmt.Sprintf("deleted: %s", strings.Join(c.deletedDst.State(), "\n")))
}

func (c *Clean) PostDo() error {
	c.dataDst.Done()
	c.deletedDst.Done()

	c.dataDst.Finish()
	c.deletedDst.Finish()

	return nil
}

func (c *Clean) Close() error { return nil }

var _ destination.Destinationer[storage.MapEntry] = (*cleanToDB)(nil)

type cleanToDB struct{}

func (i *cleanToDB) Exec(d *destination.Destination[storage.MapEntry], items storage.MapEntries) error {
	return d.DB().BulkInsertFromSliceMap(d.TableName(), items, int(d.PageSize()))
}

func newInsertToDB(c *Clean, target job.Target) (*destination.Destination[storage.MapEntry], error) {
	return destination.NewDestination[storage.MapEntry](&destination.Config{
		IsTruncate:  !c.config.NotTruncate,
		Concurrency: c.config.Concurrency,
		PageSize:    int64(c.config.BatchSize),
		TableName:   target.TableName(),
	}, &cleanToDB{}, func() (*database.DB, error) {
		return database.NewFromSimpleDb(target.GetDB()), nil
	})
}
