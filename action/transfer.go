package action

import (
	"fmt"

	go_simple_db "github.com/auho/go-simple-db/v2"
)

var _ Actionor = (*Transfer)(nil)

type Transfer struct {
	action
	target         *go_simple_db.SimpleDB
	targetDataName string
	alias          map[string]string
	fields         []string
	fixedFields    []string
	fixedData      map[string]any
}

func NewTransfer(db *go_simple_db.SimpleDB, targetDataName string, alias map[string]string, fixedData map[string]any) *Transfer {
	t := &Transfer{}
	t.target = db
	t.targetDataName = targetDataName

	if len(alias) >= 0 {
		for k := range alias {
			t.fields = append(t.fields, k)
		}

		t.alias = alias
	} else {
		var err error
		t.fields, err = db.GetTableColumns(targetDataName)
		if err != nil {
			panic(err)
		}
	}

	if len(fixedData) > 0 {
		for k := range fixedData {
			t.fixedFields = append(t.fixedFields, k)
		}

		t.fixedData = fixedData
	}

	return t
}

func (t *Transfer) GetFields() []string {
	return t.fields
}

func (t *Transfer) Title() string {
	return fmt.Sprintf("Transfer[%s]", t.targetDataName)
}

func (t *Transfer) Prepare() error {
	return nil
}

func (t *Transfer) Do(items []map[string]any) {
	targetItems := make([]map[string]any, 0)

	for _, item := range items {
		_item := make(map[string]any)
		for _, field := range t.fields {
			if ka, ok := t.alias[field]; ok {
				_item[ka] = item[field]
			} else {
				_item[field] = item[field]
			}
		}

		for k, v := range t.fixedData {
			if ka, ok := t.alias[k]; ok {
				_item[ka] = v
			} else {
				_item[k] = v
			}
		}

		targetItems = append(targetItems, _item)
	}

	_ = t.target.BulkInsertFromSliceMap(t.targetDataName, targetItems, 2000)
}

func (t *Transfer) AfterDo() {}
