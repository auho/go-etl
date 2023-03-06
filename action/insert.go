package action

import (
	"fmt"

	"github.com/auho/go-etl/mode"
	goSimpleDb "github.com/auho/go-simple-db/v2"
)

var _ Actioner = (*Insert)(nil)

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

func (i *Insert) Do(items []map[string]any) {
	newItems := make([]map[string]any, 0)

	for _, item := range items {
		_doItem := i.mode.Do(item)
		if _doItem == nil {
			continue
		}

		if len(i.affixFields) > 0 {
			for index := range _doItem {
				for _, field := range i.affixFields {
					_doItem[index][field] = item[field]
				}
			}
		}

		i.AddAmount(1)
		newItems = append(newItems, _doItem...)
	}

	err := i.db.BulkInsertFromSliceMap(i.targetTable, newItems, batchSize)
	if err != nil {
		panic(err)
	}
}

func (i *Insert) AfterDo() {}
