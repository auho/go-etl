package task

import (
	"fmt"

	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/transform"
)

var _ itemProducer = (*Transfer)(nil)

type Transfer struct {
	producerTask

	mode mode.TransferOperator
}

func NewTransfer(target job.Target, moder mode.TransferOperator) *Transfer {
	t := &Transfer{}
	t.target = target
	t.mode = moder

	return t
}

func (t *Transfer) GetFields() []string {
	return t.mode.GetFields()
}

func (t *Transfer) Summary() string {
	return fmt.Sprintf("Transfer[%s]", t.target.TableName())
}

func (t *Transfer) Prepare() error {
	return nil
}

func (t *Transfer) Exec(item map[string]any) ([]map[string]any, bool, error) {
	return []map[string]any{t.mode.Do(item)}, true, nil
}

func (t *Transfer) AppendState()     {}
func (t *Transfer) BeforeRun() error { return nil }
func (t *Transfer) AfterRun() error  { return nil }
func (t *Transfer) Close() error {
	return t.mode.Close()
}
