package task

import (
	"fmt"

	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-etl/v3/job/transform"
	"github.com/auho/go-toolkit-flow/v3/storage"
	"github.com/auho/go-toolkit-flow/v3/storage/database/destination"
)

var _ itemProducer = (*Transfer)(nil)

type Transfer struct {
	producerTask

	mode   transform.TransferOperator
	target job.Table
}

func NewTransfer(target job.Table, mode transform.TransferOperator) *Transfer {
	t := &Transfer{}
	t.target = target
	t.mode = mode

	return t
}

func (t *Transfer) Fields() ([]string, error) {
	return t.mode.Fields(), nil
}

func (t *Transfer) Summary() string {
	return fmt.Sprintf("Transfer[%s]", t.target.TableName())
}

func (t *Transfer) Prepare() error {
	return nil
}

func (t *Transfer) Exec(item map[string]any) ([]map[string]any, bool, error) {
	newItem, err := t.mode.Apply(item)
	if err != nil {
		return nil, false, fmt.Errorf("mode.Apply: %w", err)
	}

	return []map[string]any{newItem}, true, nil
}

func (t *Transfer) AppendState()     {}
func (t *Transfer) BeforeRun() error { return nil }
func (t *Transfer) AfterRun() error  { return nil }
func (t *Transfer) Close() error {
	return t.mode.Close()
}

func (t *Transfer) destinations() ([]storage.Destination[storage.MapEntry], error) {
	var ds []storage.Destination[storage.MapEntry]
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

	ds = append(ds, dest)
	return ds, nil
}
