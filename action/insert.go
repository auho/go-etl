package action

import (
	"fmt"

	"github.com/auho/go-etl/mode"
	go_simple_db "github.com/auho/go-simple-db/v2"
)

var _ Actionor = (*Insert)(nil)

type Insert struct {
	action
	target      *go_simple_db.SimpleDB
	mode        mode.InsertModer
	tagTable    string
	affixFields []string
}

func NewInsert(db *go_simple_db.SimpleDB, tagTable string, moder mode.InsertModer, affixFields []string) *Insert {
	i := &Insert{}
	i.target = db
	i.mode = moder
	i.tagTable = tagTable
	i.affixFields = affixFields

	return i
}

func (i *Insert) GetFields() []string {
	return append(i.mode.GetFields(), i.affixFields...)
}

func (i *Insert) getKeys() []string {
	return append(i.mode.GetKeys(), i.affixFields...)
}

func (i *Insert) Title() string {
	return fmt.Sprintf("Insert[%s] {%s}", i.tagTable, i.mode.GetTitle())
}

func (i *Insert) Prepare() error {
	return nil
}

func (i *Insert) Do(items []map[string]any) {
	targetItems := make([]map[string]any, 0)

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
		targetItems = append(targetItems, _doItem...)
	}

	_ = i.target.BulkInsertFromSliceMap(i.tagTable, targetItems, 2000)
}

func (i *Insert) AfterDo() {}
