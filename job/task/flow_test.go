package task

import (
	"fmt"
	"testing"

	action2 "github.com/auho/go-etl/v2/job/action"
	"github.com/auho/go-etl/v2/job/means/tag"
	mode2 "github.com/auho/go-etl/v2/job/mode"
)

func Test_Update(t *testing.T) {
	m := mode2.NewUpdateMode([]string{_keyName}, tag.NewMostKey(_rule))
	ua := action2.NewUpdate(_source, []mode2.UpdateModer{m})

	RunTask(_source, []action2.Actor{ua})
	UpdateFlow(_source, []mode2.UpdateModer{m})

	var count int64
	err := _db.Table(_dataTable).Where(fmt.Sprintf("%s != ?", "a"), "").Count(&count).Error
	if err != nil {
		t.Error(err)
	}

	amount := getAmount(_dataTable, t)
	if float32(count/amount) > 0.5 {
		t.Error("update error")
	}
}

func Test_UpdateAndTransfer(t *testing.T) {
	m := mode2.NewUpdateMode([]string{_keyName}, tag.NewMostKey(_rule))
	UpdateAndTransferFlow(_source, _targetUpdateTransfer, []mode2.UpdateModer{m})

	dataCount := getAmount(_dataTable, t)
	transferCount := getAmount(_updateAndTransferTable, t)

	if dataCount != transferCount {
		t.Error("update and transfer error")
	}
}

func Test_Insert(t *testing.T) {
	m := mode2.NewInsert([]string{_keyName}, tag.NewKey(_rule))
	ia := action2.NewInsert(_targetTagA, m, []string{_source.GetIdName()})

	_ = _db.Drop(_targetTagA1.TableName())

	err := _db.Copy(_targetTagA.TableName(), _targetTagA1.TableName())
	if err != nil {
		t.Error(err)
	}

	_ = _db.Drop(_targetTagA2.TableName())
	err = _db.Copy(_targetTagA.TableName(), _targetTagA2.TableName())
	if err != nil {
		t.Error(err)
	}

	ia1 := action2.NewInsert(_targetTagA1, m, []string{_source.GetIdName()})
	ia2 := action2.NewInsert(_targetTagA2, m, []string{_source.GetIdName()})

	RunTask(_source, []action2.Actor{ia, ia1, ia2})
	InsertFlow(_source, _targetTagA, m, []string{_source.GetIdName()})
	InsertFlow(_source, _targetTagA1, m, []string{_source.GetIdName()})
	InsertFlow(_source, _targetTagA2, m, []string{_source.GetIdName()})

	dataCount := getAmount(_source.TableName(), t)
	tagCount := getAmount(_targetTagA.TableName(), t)

	if dataCount*4 != tagCount {
		t.Error("data *4 != tag")
	}

	count1 := getAmount(_targetTagA1.TableName(), t)
	count2 := getAmount(_targetTagA2.TableName(), t)

	if count1 != count2 {
		t.Error("count 1 != count 2")
	}

	if tagCount != count1 {
		t.Error("tag count != count 1")
	}

	_ = _db.Drop(_targetTagA1.TableName())
	_ = _db.Drop(_targetTagA2.TableName())
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

	m := mode2.NewTransferMode(nil, alias, map[string]any{"xyz": "xyz1"})
	TransferFlow(_source, _targetTransfer, m)
	dataCount := getAmount(_source.TableName(), t)
	tDataCount := getAmount(_targetTransfer.TableName(), t)
	xDataCount := getFieldAmount(_targetTransfer.TableName(), "xyz", "xyz1", t)
	if tDataCount != dataCount || xDataCount != dataCount {
		t.Error("tData != data")
	}
}

func Test_Clean(t *testing.T) {
	m := mode2.NewUpdateMode([]string{_keyName}, tag.NewMostKey(_rule))

	CleanTask(_source, _targetClean, []mode2.UpdateModer{m})
	dataCount := getAmount(_source.TableName(), t)
	cDataCount := getAmount(_targetClean.TableName(), t)
	if cDataCount == dataCount || dataCount/cDataCount < 3 {
		t.Error("cData != data")
	}
}

func getAmount(tableName string, t *testing.T) int64 {
	var count int64
	err := _db.Table(tableName).Count(&count).Error
	if err != nil {
		t.Error(err)
	}

	return count
}

func getFieldAmount(tableName string, field string, value any, t *testing.T) int64 {
	var count int64
	err := _db.Table(tableName).Where(fmt.Sprintf("%s = ?", field), value).Count(&count).Error
	if err != nil {
		t.Error(err)
	}

	return count
}
