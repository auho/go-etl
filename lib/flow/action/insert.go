package action

import (
	"etl/lib/flow/mode"
	"etl/lib/storage"
)

type InsertAction struct {
	source      *storage.DbSource
	target      *storage.DbTarget
	mode        mode.InsertMode
	dataName    string
	idName      string
	affixFields []string
}

func NewInsertAction() *InsertAction {
	ia := &InsertAction{}

	return ia
}

func (ia *InsertAction) receive() {
}

func (ia *InsertAction) DoItem(item map[string]interface{}) {
	items := ia.mode.Do(item)
	if items == nil {
		return
	}

	if len(ia.affixFields) > 0 {
		for index := range items {
			for _, field := range ia.affixFields {
				items[index] = append(items[index], item[field])
			}
		}
	}
}
