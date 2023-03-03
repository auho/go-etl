package flow

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/means/tagor"
	"github.com/auho/go-etl/mode"
)

func Test_Update(t *testing.T) {
	m := mode.NewUpdate([]string{keyName}, tagor.NewMostKey(ruleName, db))
	ua := action.NewUpdate(db,
		dataTableName,
		pkName,
		[]mode.UpdateModer{m},
	)

	RunFlow(db, dataTableName, pkName, []action.Actionor{ua})
	UpdateFlow(db, dataTableName, pkName, []mode.UpdateModer{m})

	var count int64
	err := db.Table(dataTableName).Where(fmt.Sprintf("%s != ?", "a"), "").Count(&count).Error
	if err != nil {
		t.Error(err)
	}

	amount := getAmount(dataTableName, t)
	if float32(count/amount) > 0.5 {
		t.Error("update error")
	}
}

func Test_Insert(t *testing.T) {
	m := mode.NewInsert([]string{keyName}, tagor.NewKey(ruleName, db))
	ia := action.NewInsert(db,
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

	ia1 := action.NewInsert(db, tagTableName+"1", m, []string{pkName})
	ia2 := action.NewInsert(db, tagTableName+"2", m, []string{pkName})

	RunFlow(db, dataTableName, pkName, []action.Actionor{ia, ia1, ia2})
	InsertFlow(db, dataTableName, pkName, tagTableName, m, []string{pkName})
	InsertFlow(db, dataTableName, pkName, tagTableName+"1", m, []string{pkName})
	InsertFlow(db, dataTableName, pkName, tagTableName+"2", m, []string{pkName})

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

	TransferFlow(db, dataTableName, pkName, tDataTableName, alias, map[string]interface{}{"xyz": "xyz1"})
	dataCount := getAmount(dataTableName, t)
	tDataCount := getAmount(tDataTableName, t)
	xDataCount := getFieldAmount(tDataTableName, "xyz", "xyz1", t)
	if tDataCount != dataCount || xDataCount != dataCount {
		t.Error("tData != data")
	}
}

func Test_Clean(t *testing.T) {
	m := mode.NewUpdate([]string{keyName}, tagor.NewMostKey(ruleName, db))

	CleanFlow(db, dataTableName, pkName, cDataTableName, []mode.UpdateModer{m})
	dataCount := getAmount(dataTableName, t)
	cDataCount := getAmount(cDataTableName, t)
	if cDataCount == dataCount || dataCount/cDataCount < 3 {
		t.Error("cData != data")
	}
}

func getAmount(tableName string, t *testing.T) int64 {
	var count int64
	err := db.Table(tableName).Count(&count).Error
	if err != nil {
		t.Error(err)
	}

	return count
}

func getFieldAmount(tableName string, field string, value interface{}, t *testing.T) int64 {
	var count int64
	err := db.Table(tableName).Where(fmt.Sprintf("%s = ?", value)).Count(&count).Error
	if err != nil {
		t.Error(err)
	}

	return count
}
