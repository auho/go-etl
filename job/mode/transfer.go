package mode

import (
	"sort"
)

var _ TransferModer = (*TransferMode)(nil)

type TransferMode struct {
	Mode
	alias     map[string]string // alias map[table data name]output name
	aliasKeys []string          // alias data []key
	fixed     map[string]any    // fixed data map[key]value
	fixedKeys []string          // fixed data []key
}

func NewTransfer(keys []string, alias map[string]string, fixed map[string]any) *TransferMode {
	tm := &TransferMode{}
	tm.keys = keys
	tm.alias = alias
	tm.fixed = fixed

	tm.initAlias(alias)
	tm.initFixed(fixed)

	return tm
}

func (tm *TransferMode) initAlias(alias map[string]string) {
	for k := range alias {
		tm.aliasKeys = append(tm.aliasKeys, k)
	}

	sort.Slice(tm.aliasKeys, func(i, j int) bool {
		return tm.aliasKeys[i] < tm.aliasKeys[j]
	})

	tm.alias = alias
}

func (tm *TransferMode) initFixed(fixed map[string]any) {
	for k := range fixed {
		tm.fixedKeys = append(tm.fixedKeys, k)
	}

	sort.Slice(tm.fixedKeys, func(i, j int) bool {
		return tm.fixedKeys[i] < tm.fixedKeys[j]
	})

	tm.fixed = fixed
}

func (tm *TransferMode) GetTitle() string {
	return "TransferMode " + tm.Mode.getTitle()
}

func (tm *TransferMode) GetFields() []string {
	return tm.keys
}

func (tm *TransferMode) Prepare() error {
	return nil
}

func (tm *TransferMode) Do(item map[string]any) map[string]any {
	newItem := make(map[string]any)
	for _, field := range tm.keys {
		if ka, ok := tm.alias[field]; ok {
			newItem[ka] = item[field]
		} else {
			newItem[field] = item[field]
		}
	}

	for k, v := range tm.fixed {
		if ka, ok := tm.alias[k]; ok {
			newItem[ka] = v
		} else {
			newItem[k] = v
		}
	}

	return newItem
}

func (tm *TransferMode) Close() error {
	return nil
}
