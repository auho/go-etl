package flow

import (
	goEtl "github.com/auho/go-etl"
	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/mode"
)

func InsertFlow(config goEtl.DbConfig, dataName string, idName string, tagTableName string, moder mode.InsertModer, affixFields []string) {
	a := action.NewInsert(config, tagTableName, moder, affixFields)
	RunFlow(config, dataName, idName, []action.Actionor{a})
}
