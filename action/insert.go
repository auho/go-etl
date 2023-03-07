package action

import (
	"fmt"

	"github.com/auho/go-etl/mode"
	goSimpleDb "github.com/auho/go-simple-db/v2"
)

var _ Actor = (*Insert)(nil)

type Insert struct {
	action
	db          *goSimpleDb.SimpleDB
	mode        mode.InsertModer
	targetTable string
	affixFields []string
}

// NewInsertAndTransfer
// insert and transfer
func NewInsertAndTransfer(db *goSimpleDb.SimpleDB, targetTable string, moder mode.InsertModer) *Insert {
	columns, err := db.GetTableColumns(targetTable)
	if err != nil {
		panic(err)
	}

	return NewInsert(db, targetTable, moder, columns)
}

// NewInsert
// insert
func NewInsert(db *goSimpleDb.SimpleDB, targetTable string, moder mode.InsertModer, affixFields []string) *Insert {
	i := &Insert{}
	i.db = db
	i.mode = moder
	i.targetTable = targetTable
	i.affixFields = affixFields

	i.Init()

	return i
}

func (i *Insert) GetFields() []string {
	return append(i.mode.GetFields(), i.affixFields...)
}

func (i *Insert) getKeys() []string {
	return append(i.mode.GetKeys(), i.affixFields...)
}

func (i *Insert) Title() string {
	return fmt.Sprintf("Insert[%s] {%s}", i.targetTable, i.mode.GetTitle())
}

func (i *Insert) Prepare() error {
	return nil
}

func (i *Insert) Do(item map[string]any) ([]map[string]any, bool) {
	newItems := i.mode.Do(item)
	if newItems == nil {
		return nil, false
	}

	if len(i.affixFields) > 0 {
		for index := range newItems {
			for _, field := range i.affixFields {
				newItems[index][field] = item[field]
			}
		}
	}

	return newItems, true
}

func (i *Insert) PostBatchDo(items []map[string]any) {
	err := i.db.BulkInsertFromSliceMap(i.targetTable, items, batchSize)
	if err != nil {
		panic(err)
	}
}

func (i *Insert) PostDo() {}
