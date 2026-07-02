package task

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-etl/v3/job/transform"
	"github.com/auho/go-etl/v3/tool/slicex"
	"github.com/auho/go-toolkit-flow/v3/storage"
	"github.com/auho/go-toolkit-flow/v3/storage/database/destination"
)

var _ itemProducer = (*UpdateTransfer)(nil)

type UpdateTransferConfig struct {
	baseConfig
	SkipTruncate bool // for update and transfer
}

func WithUpdateTransferConfig(cc UpdateTransferConfig) func(update *UpdateTransfer) {
	return func(c *UpdateTransfer) {
		c.config = cc
	}
}

type UpdateTransfer struct {
	producerTask

	modes  []transform.UpdateOperator
	source job.Table
	target job.Table
	dst    *destination.Bulk[storage.MapEntry]
	config UpdateTransferConfig
}

func NewUpdateTransfer(source job.Table, target job.Table, modes []transform.UpdateOperator, opts ...func(*UpdateTransfer)) *UpdateTransfer {
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

func (u *UpdateTransfer) Fields() ([]string, error) {
	fields := make([]string, 0)
	fields = append(fields, u.source.IDName())

	for _, m := range u.modes {
		fields = append(fields, m.Fields()...)
	}

	columns, err := u.target.GetDB().GetTableColumns(u.target.TableName())
	if err != nil {
		return nil, fmt.Errorf("GetTableColumns: %w", err)
	}

	fields = append(fields, columns...)
	fields = slicex.SliceDropDuplicates(fields)

	return fields, nil
}

func (u *UpdateTransfer) Summary() string {
	s := make([]string, 0)
	for _, m := range u.modes {
		s = append(s, m.Title())
	}

	return fmt.Sprintf("UpdateTransfer[%s] {%s}", u.source.TableName(), strings.Join(s, ", "))
}

func (u *UpdateTransfer) Prepare() error {
	var err error
	for _, m := range u.modes {
		err = m.Prepare()
		if err != nil {
			return fmt.Errorf("prepare: %w", err)
		}
	}

	return nil
}

func (u *UpdateTransfer) BeforeRun() error { return nil }

func (u *UpdateTransfer) Exec(item map[string]any) ([]map[string]any, bool, error) {
	does := make(map[string]any)
	for _, m := range u.modes {
		do, err := m.Apply(item)
		if err != nil {
			return nil, false, fmt.Errorf("apply: %w", err)
		}
		for k, v := range do {
			does[k] = v
		}
	}

	if len(does) <= 0 {
		return nil, false, nil
	}

	for k, v := range does {
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

func (u *UpdateTransfer) destinations() ([]storage.Destination[storage.MapEntry], error) {
	var ds []storage.Destination[storage.MapEntry]
	dest, err := destination.NewBulkInsertMapWithGorm(
		destination.BulkConfig{
			IsTruncate:  !u.config.SkipTruncate,
			Concurrency: u.config.Concurrency,
			PageSize:    int64(u.config.BatchSize),
		},
		destination.WriteConfig{
			TableName: u.target.TableName(),
		},
		u.target.GetDB().GormDB(),
	)
	if err != nil {
		return nil, fmt.Errorf("NewBulkInsertMapWithGorm: %w", err)
	}

	ds = append(ds, dest)
	return ds, nil
}
