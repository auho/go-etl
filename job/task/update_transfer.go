package task

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/transform"
	slices "github.com/auho/go-etl/v2/tool/slicex"
	"github.com/auho/go-toolkit-flow/storage"
	"github.com/auho/go-toolkit-flow/storage/database/destination"
)

var _ itemProducer = (*UpdateTransfer)(nil)

type UpdateTransferConfig struct {
	NotTruncate bool // for update and transfer
	BatchSize   int  // for update and transfer
	Concurrency int  // for update and transfer
}

func (uc *UpdateTransferConfig) check() {
	if uc.BatchSize <= 0 {
		uc.BatchSize = batchSize
	}

	if uc.Concurrency <= 0 {
		uc.Concurrency = runtime.NumCPU()
	}
}

func WithUpdateTransferConfig(cc UpdateTransferConfig) func(update *UpdateTransfer) {
	return func(c *UpdateTransfer) {
		c.config = cc
	}
}

type UpdateTransfer struct {
	producerTask

	source job.Source
	modes  []transform.UpdateOperator

	config UpdateTransferConfig
	dst    *destination.Bulk[storage.MapEntry]
}

func NewUpdateTransfer(source job.Source, target job.Target, modes []transform.UpdateOperator, opts ...func(*UpdateTransfer)) *UpdateTransfer {
	u := &UpdateTransfer{}
	u.source = source
	u.modes = modes
	u.target = target

	for _, opt := range opts {
		opt(u)
	}

	u.config.check()

	return u
}

func (u *UpdateTransfer) GetFields() []string {
	fields := make([]string, 0)
	fields = append(fields, u.source.GetIDName())

	for _, m := range u.modes {
		fields = append(fields, m.GetFields()...)
	}

	columns, err := u.target.GetDB().GetTableColumns(u.target.TableName())
	if err != nil {
		panic(err)
	}

	fields = append(fields, columns...)
	fields = slices.SliceDropDuplicates(fields)

	return fields
}

func (u *UpdateTransfer) Summary() string {
	s := make([]string, 0)
	for _, m := range u.modes {
		s = append(s, m.GetTitle())
	}

	return fmt.Sprintf("UpdateTransfer[%s] {%s}", u.source.TableName(), strings.Join(s, ", "))
}

func (u *UpdateTransfer) Prepare() error {
	var err error
	for _, m := range u.modes {
		err = m.Prepare()
		if err != nil {
			return fmt.Errorf("update action prepare error; %w", err)
		}
	}

	return nil
}

func (u *UpdateTransfer) BeforeRun() error { return nil }

func (u *UpdateTransfer) Exec(item map[string]any) ([]map[string]any, bool, error) {
	_does := make(map[string]any)
	for _, m := range u.modes {
		_do := m.Do(item)
		for k, v := range _do {
			_does[k] = v
		}
	}

	if len(_does) <= 0 {
		return nil, false, nil
	}

	for k, v := range _does {
		item[k] = v
	}

	return []map[string]any{item}, true, nil
}

func (u *UpdateTransfer) AppendState() {}

func (u *UpdateTransfer) AfterRun() error { return nil }

func (u *UpdateTransfer) Close() error {
	for _, m := range u.modes {
		err := m.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
