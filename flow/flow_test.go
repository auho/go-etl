package flow

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/means/tag"
	"github.com/auho/go-etl/mode"
)

func Test_Update(t *testing.T) {
	m := mode.NewUpdate([]string{keyName}, tag.NewMostKey(ruleName, tag.WithDBRule(db)))
	ua := action.NewUpdate(db,
		dataTable,
		pkName,
		[]mode.UpdateModer{m},
	)

	RunFlow(db, dataTable, pkName, []action.Actioner{ua})
	UpdateFlow(db, dataTable, pkName, []mode.UpdateModer{m})

	var count int64
	err := db.Table(dataTable).Where(fmt.Sprintf("%s != ?", "a"), "").Count(&count).Error
	if err != nil {
		t.Error(err)
	}

	amount := getAmount(dataTable, t)
	if float32(count/amount) > 0.5 {
		t.Error("update error")
	}
}

func Test_UpdateAndTransfer(t *testing.T) {
	m := mode.NewUpdate([]string{keyName}, tag.NewMostKey(ruleName, tag.WithDBRule(db)))
	UpdateAndTransferFlow(db, dataTable, pkName, updateAndTransferTable, []mode.UpdateModer{m})

	dataCount := getAmount(dataTable, t)
	transferCount := getAmount(updateAndTransferTable, t)

	if dataCount != transferCount {
		t.Error("update and transfer error")
	}
}

func Test_Insert(t *testing.T) {
	m := mode.NewInsert([]string{keyName}, tag.NewKey(ruleName, tag.WithDBRule(db)))
	ia := action.NewInsert(db,
		tagATable,
		m,
		[]string{pkName},
	)

	_ = db.Drop(tagATable + "1")

	err := db.Copy(tagATable, tagATable+"1")
	if err != nil {
		t.Error(err)
	}

	_ = db.Drop(tagATable + "2")
	err = db.Copy(tagATable, tagATable+"2")
	if err != nil {
		t.Error(err)
	}

	ia1 := action.NewInsert(db, tagATable+"1", m, []string{pkName})
	ia2 := action.NewInsert(db, tagATable+"2", m, []string{pkName})

	RunFlow(db, dataTable, pkName, []action.Actioner{ia, ia1, ia2})
	InsertFlow(db, dataTable, pkName, tagATable, m, []string{pkName})
	InsertFlow(db, dataTable, pkName, tagATable+"1", m, []string{pkName})
	InsertFlow(db, dataTable, pkName, tagATable+"2", m, []string{pkName})

	dataCount := getAmount(dataTable, t)
	tagCount := getAmount(tagATable, t)

	if dataCount*4 != tagCount {
		t.Error("data *4 != tag")
	}

	count1 := getAmount(tagATable+"1", t)
	count2 := getAmount(tagATable+"2", t)

	if count1 != count2 {
		t.Error("count 1 != count 2")
	}

	if tagCount != count1 {
		t.Error("tag count != count 1")
	}

	_ = db.Drop(tagATable + "1")
	_ = db.Drop(tagATable + "2")

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

	TransferFlow(db, dataTable, pkName, transferTable, alias, map[string]interface{}{"xyz": "xyz1"})
	dataCount := getAmount(dataTable, t)
	tDataCount := getAmount(transferTable, t)
	xDataCount := getFieldAmount(transferTable, "xyz", "xyz1", t)
	if tDataCount != dataCount || xDataCount != dataCount {
		t.Error("tData != data")
	}
}

func Test_Clean(t *testing.T) {
	m := mode.NewUpdate([]string{keyName}, tag.NewMostKey(ruleName, tag.WithDBRule(db)))

	CleanFlow(db, dataTable, pkName, cleanTable, []mode.UpdateModer{m})
	dataCount := getAmount(dataTable, t)
	cDataCount := getAmount(cleanTable, t)
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
	err := db.Table(tableName).Where(fmt.Sprintf("%s = ?", field), value).Count(&count).Error
	if err != nil {
		t.Error(err)
	}

	return count
}
