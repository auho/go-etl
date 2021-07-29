package flow

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/means/tager"
	"github.com/auho/go-etl/mode"
)

func Test_FlowUpdate(t *testing.T) {
	m := mode.NewTagUpdate([]string{keyName}, tager.NewTagMostKeyMeans(ruleName, db))
	ua := action.NewUpdateAction(dbConfig,
		dataTableName,
		pkName,
		[]mode.UpdateModer{m},
	)

	RunFlow(dbConfig, dataTableName, pkName, []action.Action{ua})
	UpdateFlow(dbConfig, dataTableName, pkName, []mode.UpdateModer{m})

	query := fmt.Sprintf("SELECT COUNT(*) AS _count FROM `%s` WHERE `%s` != ?", dataTableName, "a")
	res, err := db.QueryFieldInterface("_count", query, "")
	if err != nil {
		t.Error(err)
	}

	count := res.(int64)
	if count <= 0 {
		t.Error("update error")
	}

	amount := getAmount(dataTableName, t)
	if float32(count/amount) > 0.5 {
		t.Error("update error")
	}
}

func Test_FlowInsert(t *testing.T) {
	m := mode.NewInsertMode([]string{keyName}, tager.NewTagKeyMeans(ruleName, db))
	ia := action.NewInsertAction(dbConfig,
		tagTableName,
		m,
		[]string{pkName},
	)

	_ = db.Drop(tagTableName + "1")

	err := db.Copy(tagTableName, tagTableName+"1")
	if err != nil {
		t.Error(err)
	}

	_ = db.Drop(tagTableName + "2")
	err = db.Copy(tagTableName, tagTableName+"2")
	if err != nil {
		t.Error(err)
	}

	ia1 := action.NewInsertAction(dbConfig, tagTableName+"1", m, []string{pkName})
	ia2 := action.NewInsertAction(dbConfig, tagTableName+"2", m, []string{pkName})

	RunFlow(dbConfig, dataTableName, pkName, []action.Action{ia, ia1, ia2})
	InsertFlow(dbConfig, dataTableName, pkName, tagTableName, m, []string{pkName})
	InsertFlow(dbConfig, dataTableName, pkName, tagTableName+"1", m, []string{pkName})
	InsertFlow(dbConfig, dataTableName, pkName, tagTableName+"2", m, []string{pkName})

	dataCount := getAmount(dataTableName, t)
	tagCount := getAmount(tagTableName, t)

	if dataCount*4 != tagCount {
		t.Error("data *2 != tag")
	}

	count1 := getAmount(tagTableName+"1", t)
	count2 := getAmount(tagTableName+"2", t)

	if count1 != count2 {
		t.Error("count 1 != count 2")
	}

	if tagCount != count1 {
		t.Error("tag count != count 1")
	}

	_ = db.Drop(tagTableName + "1")
	_ = db.Drop(tagTableName + "2")
}

func getAmount(tableName string, t *testing.T) int64 {
	query := fmt.Sprintf("SELECT COUNT(*) AS _count FROM `%s`", tableName)
	res, err := db.QueryFieldInterface("_count", query)
	if err != nil {
		t.Error(err)
	}

	count, err := strconv.Atoi(string(res.([]uint8)))
	if err != nil {
		t.Error(err)
	}

	return int64(count)
}
