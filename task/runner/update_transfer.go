package runner

import (
	"context"
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/task"
	"github.com/auho/go-etl/v3/task/transform"
	"github.com/auho/go-etl/v3/tool/slicex"
	"github.com/auho/go-toolkit-flow/v3/storage"
	"github.com/auho/go-toolkit-flow/v3/storage/database/destination"
)

var _ itemProducer = (*UpdateTransfer)(nil)

// UpdateTransferConfig configures an UpdateTransfer task.
type UpdateTransferConfig struct {
	baseConfig
	SkipTruncate bool // for update and transfer
}

// WithUpdateTransferConfig returns an option that sets the UpdateTransfer config.
func WithUpdateTransferConfig(cc UpdateTransferConfig) func(update *UpdateTransfer) {
	return func(c *UpdateTransfer) {
		c.config = cc
	}
}

// UpdateTransfer is a producer that updates each source row in place via operators
// and copies the updated row to the target table.
type UpdateTransfer struct {
	producerTask

	operators []transform.UpdateOperator
	source    task.Table
	target    task.Table
	dst       *destination.Bulk[storage.MapEntry]
	config    UpdateTransferConfig
}

// NewUpdateTransfer creates an UpdateTransfer producer from source to target.
func NewUpdateTransfer(source task.Table, target task.Table, operators []transform.UpdateOperator, opts ...func(*UpdateTransfer)) *UpdateTransfer {
	u := &UpdateTransfer{}
	u.source = source
	u.operators = operators
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

	for _, op := range u.operators {
		fields = append(fields, op.Fields()...)
	}

	columns, err := u.target.GetDB().GetTableColumns(context.TODO(), u.target.TableName())
	if err != nil {
		return nil, fmt.Errorf("GetTableColumns: %w", err)
	}

	fields = append(fields, columns...)
	fields = slicex.SliceDropDuplicates(fields)

	return fields, nil
}

func (u *UpdateTransfer) Summary() string {
	s := make([]string, 0)
	for _, op := range u.operators {
		s = append(s, op.Title())
	}

	return fmt.Sprintf("UpdateTransfer[%s] {%s}", u.source.TableName(), strings.Join(s, ", "))
}

func (u *UpdateTransfer) Prepare() error {
	var err error
	for _, op := range u.operators {
		err = op.Prepare()
		if err != nil {
			return fmt.Errorf("prepare: %w", err)
		}
	}

	return nil
}

func (u *UpdateTransfer) BeforeRun() error { return nil }

func (u *UpdateTransfer) Exec(item map[string]any) ([]map[string]any, bool, error) {
	does := make(map[string]any)
	for _, op := range u.operators {
		do, err := op.Apply(item)
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
	for _, op := range u.operators {
		err := op.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *UpdateTransfer) destinations() ([]storage.Destination[storage.MapEntry], error) {
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

	return []storage.Destination[storage.MapEntry]{dest}, nil
}
