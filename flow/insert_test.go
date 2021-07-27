package flow

import (
	"fmt"
	"strconv"
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

	query := fmt.Sprintf("SELECT COUNT(*) AS _count FROM `%s`", dataTableName)
	res, err := db.QueryFieldInterface("_count", query)
	if err != nil {
		t.Error(err)
	}

	dataCount, err := strconv.Atoi(string(res.([]uint8)))
	if err != nil {
		t.Error(err)
	}

	query = fmt.Sprintf("SELECT COUNT(*) AS _count FROM `%s`", tagTableName)
	res, err = db.QueryFieldInterface("_count", query)
	if err != nil {
		t.Error(err)
	}

	tagCount, err := strconv.Atoi(string(res.([]uint8)))

	if dataCount*2 != tagCount {
		t.Error("data *2 != tag")
	}
}
