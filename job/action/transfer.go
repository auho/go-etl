package action

import (
	"fmt"

	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/mode"
)

var _ Actor = (*Transfer)(nil)

type Transfer struct {
	action

	mode mode.TransferModer
}

func NewTransfer(target job.Target, moder mode.TransferModer) *Transfer {
	t := &Transfer{}
	t.target = target
	t.mode = moder

	t.Init()

	return t
}

func (t *Transfer) GetFields() []string {
	return t.mode.GetFields()
}

func (t *Transfer) Title() string {
	return fmt.Sprintf("Transfer[%s]", t.target.TableName())
}

func (t *Transfer) Prepare() error {
	return nil
}

func (t *Transfer) Do(item map[string]any) ([]map[string]any, bool) {
	return []map[string]any{t.mode.Do(item)}, true
}

func (t *Transfer) PostBatchDo(items []map[string]any) {
	err := t.target.GetDB().BulkInsertFromSliceMap(t.target.TableName(), items, batchSize)
	if err != nil {
		panic(err)
	}
}

func (t *Transfer) PostDo() {}
