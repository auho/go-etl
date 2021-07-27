package flow

import (
	"testing"

	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/means/tager"
	"github.com/auho/go-etl/mode"
)

func TestInsertFlow(t *testing.T) {
	ia := action.NewInsertAction(dbConfig,
		tagTableName,
		mode.NewTagInsert([]string{keyName}, tager.NewTagKeyMeans(ruleName, db)),
		[]string{pkName})

	RunInsertFlow(dbConfig, dataTableName, pkName, []*action.InsertAction{ia})
}
