package table

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/v2/insight/model"
)

var _tableName = "table1"
var _tableIdName = "tid"
var _ruleName = "rule1"
var _data *model.Data
var _rule *model.Rule
var _rule1 *model.Rule
var _dataRule *model.DataRule
var _dataRule1 *model.DataRule
var _tag *model.Tag
var _tag1 *model.Tag
var _dcSegWords *model.DataContentSegWords
var _dcSpiltWords *model.DataContentSpiltWords

func TestData(t *testing.T) {
	dt := NewDataTable(_data)
	dt.AddPkBigInt("id")
	dt.AddInt("int1")
	dt.AddKeyInt("int2")
	dt.AddUniqueInt("int3")
	dt.AddString("s1")
	dt.AddStringWithLength("s2", 125)
	dt.AddKeyString("s3", 20, 0)
	dt.AddUniqueString("s4", 20)
	dt.AddTimestamp("ts1", false, false)
	dt.AddTimestamp("ts2", false, true)
	dt.AddTimestamp("ts3", true, false)
	dt.AddTimestamp("ts4", true, true)
	dt.AddText("t1")

	sql := dt.GetTable().SqlForCreate()
	fmt.Println(sql)
}

func TestRule(t *testing.T) {
	rt := NewRuleTable(_rule)

	sql := rt.GetTable().SqlForCreate()
	fmt.Println(sql)

	rt1 := NewRuleTable(_rule1)

	sql = rt1.GetTable().SqlForCreate()
	fmt.Println(sql)
}

func TestDataRule(t *testing.T) {
	dr := NewDataRuleTable(_dataRule)

	sql := dr.GetTable().SqlForCreate()
	fmt.Println(sql)

	dr1 := NewDataRuleTable(_dataRule1)

	sql = dr1.GetTable().SqlForCreate()
	fmt.Println(sql)
}

func TestTag(t *testing.T) {
	tt := NewTagTable(_tag)

	sql := tt.GetTable().SqlForCreate()
	fmt.Println(sql)

	tt1 := NewTagTable(_tag1)
	sql = tt1.GetTable().SqlForCreate()
	fmt.Println(sql)
}

func TestDataContent(t *testing.T) {
	dcSeg := NewDataContentSegWordsTable(_dcSegWords)

	sql := dcSeg.GetTable().SqlForCreate()
	fmt.Println(sql)

	dcSplit := NewDataContentSpiltWordsTable(_dcSpiltWords)

	sql = dcSplit.GetTable().SqlForCreate()
	fmt.Println(sql)
}

func init() {
	_data = model.NewData(_tableName, _tableIdName)
	_rule = model.NewRule(_ruleName, 20, 20, nil)
	_rule1 = model.NewRule(_ruleName, 20, 20, map[string]int{"r1": 10, "r2": 30})
	_dataRule = model.NewDataRule(_data, _rule)
	_dataRule1 = model.NewDataRule(_data, _rule1)
	_tag = model.NewTag(_data, _rule)
	_tag1 = model.NewTag(_data, _rule1)
	_dcSegWords = model.NewDataContentSegWords(_data, "abc", 30)
	_dcSpiltWords = model.NewDataContentSpiltWords(_data, "abc", 30)
}
