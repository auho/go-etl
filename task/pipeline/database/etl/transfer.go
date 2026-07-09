package etl

import (
	"fmt"

	"github.com/auho/go-etl/v3/task/transform"
	"github.com/auho/go-toolkit-flow/v3/processor/producer"
	"github.com/auho/go-toolkit-flow/v3/storage"
	"github.com/auho/go-toolkit-flow/v3/storage/database/destination"
)

var _ producer.Item[storage.MapEntry, storage.MapEntry] = (*Transfer)(nil)

// Transfer is a producer that transforms each source row via operator and writes it
// to the target table (target is truncated first).
type Transfer struct {
	producer.Processor
	TaskBase

	operator transform.TransferOperator
	target   Table
}

// NewTransfer creates a Transfer producer writing to target via operator.
func NewTransfer(target Table, operator transform.TransferOperator) *Transfer {
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

func (t *Transfer) Prepare() error { return nil }

func (t *Transfer) BeforeRun() error { return nil }

func (t *Transfer) Exec(item storage.MapEntry) ([]storage.MapEntry, bool, error) {
	newItem, err := t.operator.Apply(item)
	if err != nil {
		return nil, false, fmt.Errorf("operator.Apply: %w", err)
	}

	return []map[string]any{newItem}, true, nil
}

func (t *Transfer) AppendState()    {}
func (t *Transfer) AfterRun() error { return nil }
func (t *Transfer) Close() error {
	return t.operator.Close()
}

// BuildDestinations constructs and returns the insert destination for the target table.
func (t *Transfer) BuildDestinations() ([]storage.Destination[storage.MapEntry], error) {
	dest, err := destination.NewBulkInsertMapWithGorm(
		destination.BulkConfig{
			IsTruncate:  true,
			Concurrency: t.Concurrency(),
			BatchSize:   int64(batchSize),
		},
		destination.WriteConfig{TableName: t.target.TableName()},
		t.target.DB().GormDB(),
	)
	if err != nil {
		return nil, fmt.Errorf("NewBulkInsertMapWithGorm: %w", err)
	}

	return []storage.Destination[storage.MapEntry]{dest}, nil
}
