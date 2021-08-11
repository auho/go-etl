package flow

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/means/tagor"
	"github.com/auho/go-etl/mode"
)

func Test_Update(t *testing.T) {
	m := mode.NewUpdate([]string{keyName}, tagor.NewMostKey(ruleName, db))
	ua := action.NewUpdate(dbConfig,
		dataTableName,
		pkName,
		[]mode.UpdateModer{m},
	)

	RunFlow(dbConfig, dataTableName, pkName, []action.Actionor{ua})
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

func Test_Insert(t *testing.T) {
	m := mode.NewInsert([]string{keyName}, tagor.NewKey(ruleName, db))
	ia := action.NewInsert(dbConfig,
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

	ia1 := action.NewInsert(dbConfig, tagTableName+"1", m, []string{pkName})
	ia2 := action.NewInsert(dbConfig, tagTableName+"2", m, []string{pkName})

	RunFlow(dbConfig, dataTableName, pkName, []action.Actionor{ia, ia1, ia2})
	InsertFlow(dbConfig, dataTableName, pkName, tagTableName, m, []string{pkName})
	InsertFlow(dbConfig, dataTableName, pkName, tagTableName+"1", m, []string{pkName})
	InsertFlow(dbConfig, dataTableName, pkName, tagTableName+"2", m, []string{pkName})

	dataCount := getAmount(dataTableName, t)
	tagCount := getAmount(tagTableName, t)

	if dataCount*2 != tagCount {
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

func Test_Transfer(t *testing.T) {
	alias := map[string]string{
		"did":           "did",
		"name":          "name",
		"a":             "a1",
		"ab":            "ab1",
		"a_keyword":     "a_keyword",
		"a_keyword_num": "a_keyword_num",
	}

	TransferFlow(db, dbConfig, dataTableName, pkName, tDataTableName, alias, map[string]interface{}{"xyz": "xyz1"})
	dataCount := getAmount(dataTableName, t)
	tDataCount := getAmount(tDataTableName, t)
	xDataCount := getFieldAmount(tDataTableName, "xyz", "xyz1", t)
	if tDataCount != dataCount || xDataCount != dataCount {
		t.Error("tData != data")
	}
}

func Test_Clean(t *testing.T) {
	m := mode.NewUpdate([]string{keyName}, tagor.NewMostKey(ruleName, db))

	CleanFlow(db, dbConfig, dataTableName, pkName, cDataTableName, []mode.UpdateModer{m})
	dataCount := getAmount(dataTableName, t)
	cDataCount := getAmount(cDataTableName, t)
	if cDataCount == dataCount || dataCount/cDataCount < 3 {
		t.Error("cData != data")
	}
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

func getFieldAmount(tableName string, field string, value interface{}, t *testing.T) int64 {
	query := fmt.Sprintf("SELECT COUNT(*) AS _count FROM `%s` WHERE `%s` = ?", tableName, field)
	res, err := db.QueryFieldInterface("_count", query, value)
	if err != nil {
		t.Error(err)
	}

	return res.(int64)
}
