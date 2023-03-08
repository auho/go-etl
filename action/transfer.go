package action

import (
	"fmt"

	goSimpleDb "github.com/auho/go-simple-db/v2"
)

var _ Actor = (*Transfer)(nil)

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

func (t *Transfer) Do(item map[string]any) ([]map[string]any, bool) {
	newItem := make(map[string]any)
	for _, field := range t.fields {
		if ka, ok := t.alias[field]; ok {
			newItem[ka] = item[field]
		} else {
			newItem[field] = item[field]
		}
	}

	for k, v := range t.fixedData {
		if ka, ok := t.alias[k]; ok {
			newItem[ka] = v
		} else {
			newItem[k] = v
		}
	}

	return []map[string]any{newItem}, true
}

func (t *Transfer) PostBatchDo(items []map[string]any) {
	err := t.db.BulkInsertFromSliceMap(t.targetTable, items, batchSize)
	if err != nil {
		panic(err)
	}
}

func (t *Transfer) PostDo() {}
