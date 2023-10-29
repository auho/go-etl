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
var _data_rule *model.DataRule
var _data_rule1 *model.DataRule
var _tag *model.Tag
var _tag1 *model.Tag

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
	dr := NewDataRuleTable(_data_rule)

	sql := dr.GetTable().SqlForCreate()
	fmt.Println(sql)

	dr1 := NewDataRuleTable(_data_rule1)

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

func init() {
	_data = model.NewData(_tableName, _tableIdName)
	_rule = model.NewRule(_ruleName, 20, 20, nil)
	_rule1 = model.NewRule(_ruleName, 20, 20, map[string]int{"r1": 10, "r2": 30})
	_data_rule = model.NewDataRule(_data, _rule)
	_data_rule1 = model.NewDataRule(_data, _rule1)
	_tag = model.NewTag(_data, _rule)
	_tag1 = model.NewTag(_data, _rule1)
}
