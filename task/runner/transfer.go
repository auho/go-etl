package runner

import (
	"fmt"

	"github.com/auho/go-etl/v3/task"
	"github.com/auho/go-etl/v3/task/transform"
	"github.com/auho/go-toolkit-flow/v3/storage"
	"github.com/auho/go-toolkit-flow/v3/storage/database/destination"
)

var _ itemProducer = (*Transfer)(nil)

// Transfer is a producer that transforms each source row via operator and writes it
// to the target table (target is truncated first).
type Transfer struct {
	producerTask

	operator transform.TransferOperator
	target   task.Table
}

// NewTransfer creates a Transfer producer writing to target via operator.
func NewTransfer(target task.Table, operator transform.TransferOperator) *Transfer {
	t := &Transfer{}
	t.target = target
	t.operator = operator

	return t
}

func (t *Transfer) Fields() ([]string, error) {
	return t.operator.Fields(), nil
}

func (t *Transfer) Summary() string {
	return fmt.Sprintf("Transfer[%s]", t.target.TableName())
}

func (t *Transfer) Prepare() error {
	return nil
}

func (t *Transfer) Exec(item map[string]any) ([]map[string]any, bool, error) {
	newItem, err := t.operator.Apply(item)
	if err != nil {
		return nil, false, fmt.Errorf("operator.Apply: %w", err)
	}

	return []map[string]any{newItem}, true, nil
}

func (t *Transfer) AppendState()     {}
func (t *Transfer) BeforeRun() error { return nil }
func (t *Transfer) AfterRun() error  { return nil }
func (t *Transfer) Close() error {
	return t.operator.Close()
}

func (t *Transfer) destinations() ([]storage.Destination[storage.MapEntry], error) {
	dest, err := destination.NewBulkInsertMapWithGorm(
		destination.BulkConfig{
			IsTruncate:  true,
			Concurrency: t.Concurrency(),
			PageSize:    batchSize,
		},
		destination.WriteConfig{
			TableName: t.target.TableName(),
		},
		t.target.GetDB().GormDB(),
	)
	if err != nil {
		return nil, fmt.Errorf("NewBulkInsertMapWithGorm: %w", err)
	}

	return []storage.Destination[storage.MapEntry]{dest}, nil
}