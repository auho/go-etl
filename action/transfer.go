package action

import (
	"fmt"

	goSimpleDb "github.com/auho/go-simple-db/v2"
)

var _ Actionor = (*Transfer)(nil)

type Transfer struct {
	action
	db          *goSimpleDb.SimpleDB
	targetTable string
	fields      []string
	alias       map[string]string // alias map[table data name]output name
	fixedFields []string          // fixed data []key
	fixedData   map[string]any    // fixed data map[key]value
}

func NewTransfer(db *goSimpleDb.SimpleDB, targetTable string, alias map[string]string, fixedData map[string]any) *Transfer {
	t := &Transfer{}
	t.db = db
	t.targetTable = targetTable

	if len(alias) >= 0 {
		for k := range alias {
			t.fields = append(t.fields, k)
		}

		t.alias = alias
	} else {
		var err error
		t.fields, err = db.GetTableColumns(targetTable)
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

	t.Init()

	return t
}

func (t *Transfer) GetFields() []string {
	return t.fields
}

func (t *Transfer) Title() string {
	return fmt.Sprintf("Transfer[%s]", t.targetTable)
}

func (t *Transfer) Prepare() error {
	return nil
}

func (t *Transfer) Do(items []map[string]any) {
	newItems := make([]map[string]any, 0)

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

		newItems = append(newItems, _item)
	}

	err := t.db.BulkInsertFromSliceMap(t.targetTable, newItems, batchSize)
	if err != nil {
		panic(err)
	}
}

func (t *Transfer) AfterDo() {}
