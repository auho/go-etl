package action

import (
	"fmt"
	"sort"
)

var _ Actor = (*Transfer)(nil)

type Transfer struct {
	action
	target    Target
	keys      []string
	alias     map[string]string // alias map[table data name]output name
	aliasKeys []string          // alias data []key
	fixed     map[string]any    // fixed data map[key]value
	fixedKeys []string          // fixed data []key
}

func NewTransfer(target Target, keys []string, alias map[string]string, fixed map[string]any) *Transfer {
	t := &Transfer{}
	t.target = target

	var err error
	if len(keys) <= 0 {
		t.keys, err = target.GetDB().GetTableColumns(target.TableName())
		if err != nil {
			panic(err)
		}
	} else {
		t.keys = keys
	}

	t.initAlias(alias)
	t.initFixed(fixed)
	t.Init()

	return t
}

func (t *Transfer) initAlias(alias map[string]string) {
	for k := range alias {
		t.aliasKeys = append(t.aliasKeys, k)
	}

	sort.Slice(t.aliasKeys, func(i, j int) bool {
		return t.aliasKeys[i] < t.aliasKeys[j]
	})

	t.alias = alias
}

func (t *Transfer) initFixed(fixed map[string]any) {
	for k := range fixed {
		t.fixedKeys = append(t.fixedKeys, k)
	}

	sort.Slice(t.fixedKeys, func(i, j int) bool {
		return t.fixedKeys[i] < t.fixedKeys[j]
	})

	t.fixed = fixed
}

func (t *Transfer) GetFields() []string {
	return t.keys
}

func (t *Transfer) Title() string {
	return fmt.Sprintf("Transfer[%s]", t.target.TableName())
}

func (t *Transfer) Prepare() error {
	return nil
}

func (t *Transfer) Do(item map[string]any) ([]map[string]any, bool) {
	newItem := make(map[string]any)
	for _, field := range t.keys {
		if ka, ok := t.alias[field]; ok {
			newItem[ka] = item[field]
		} else {
			newItem[field] = item[field]
		}
	}

	for k, v := range t.fixed {
		if ka, ok := t.alias[k]; ok {
			newItem[ka] = v
		} else {
			newItem[k] = v
		}
	}

	return []map[string]any{newItem}, true
}

func (t *Transfer) PostBatchDo(items []map[string]any) {
	err := t.target.GetDB().BulkInsertFromSliceMap(t.target.TableName(), items, batchSize)
	if err != nil {
		panic(err)
	}
}

func (t *Transfer) PostDo() {}
