package transform

import (
	"sort"
)

var _ TransferOperator = (*Transfer)(nil)

type Transfer struct {
	base
	alias     map[string]string // alias map[table data name]output name
	aliasKeys []string          // alias data []key
	fixed     map[string]any    // fixed data map[key]value
	fixedKeys []string          // fixed data []key
}

func NewTransfer(keys []string, alias map[string]string, fixed map[string]any) *Transfer {
	tm := &Transfer{}
	tm.keys = keys
	tm.alias = alias
	tm.fixed = fixed

	tm.initAlias(alias)
	tm.initFixed(fixed)

	return tm
}

func (tm *Transfer) initAlias(alias map[string]string) {
	for k := range alias {
		tm.aliasKeys = append(tm.aliasKeys, k)
	}

	sort.Slice(tm.aliasKeys, func(i, j int) bool {
		return tm.aliasKeys[i] < tm.aliasKeys[j]
	})

	tm.alias = alias
}

func (tm *Transfer) initFixed(fixed map[string]any) {
	for k := range fixed {
		tm.fixedKeys = append(tm.fixedKeys, k)
	}

	sort.Slice(tm.fixedKeys, func(i, j int) bool {
		return tm.fixedKeys[i] < tm.fixedKeys[j]
	})

	tm.fixed = fixed
}

func (tm *Transfer) Title() string {
	return tm.genTitle("Transfer", "")
}

func (tm *Transfer) Fields() []string {
	return tm.keys
}

func (tm *Transfer) Prepare() error {
	return nil
}

func (tm *Transfer) Apply(item map[string]any) (map[string]any, error) {
	newItem := make(map[string]any)
	for _, field := range tm.keys {
		if ka, ok := tm.alias[field]; ok {
			newItem[ka] = item[field]
		} else {
			newItem[field] = item[field]
		}
	}

	// Apply fixed values. Alias results take precedence: if a fixed key
	// collides with an alias output key already set above, the alias
	// value wins and the fixed value is skipped.
	for k, v := range tm.fixed {
		if ka, ok := tm.alias[k]; ok {
			newItem[ka] = v
		} else if _, exists := newItem[k]; !exists {
			newItem[k] = v
		}
	}

	return newItem, nil
}

func (tm *Transfer) Close() error {
	return nil
}
