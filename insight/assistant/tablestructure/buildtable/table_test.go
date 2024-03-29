package buildtable

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/model"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

var _tableName = "table1"
var _tableIdName = "tid"
var _ruleName = "rule1"
var _raw *model.Raw
var _data *model.Data
var _rule *model.Rule
var _rule1 *model.Rule
var _dataRule *model.DataRule
var _dataRule1 *model.DataRule
var _tagRule *model.TagDataRule
var _tagRule1 *model.TagDataRule
var _tagRules *model.TagDataRules
var _tagRules1 *model.TagDataRules
var _dcSegWords *model.DataContentSegWords
var _dcSpiltWords *model.DataContentSpiltWords

func TestRaw(t *testing.T) {
	rt := NewRawTable(_raw)
	rt.AddPkBigInt("id")
	rt.WithCommand(func(command *tablestructure.Command) {
		command.AddString("with")
	})

	sql := rt.Sql()
	fmt.Println(sql)
}

func TestData(t *testing.T) {
	dt := NewDataTable(_data)
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
	dt.WithCommand(func(command *tablestructure.Command) {
		command.AddString("with")
	})

	sql := dt.Sql()
	fmt.Println(sql)
}

func TestRule(t *testing.T) {
	rt := NewRuleTable(_rule).
		WithCommand(func(command *tablestructure.Command) {
			command.AddString("with")
		})

	sql := rt.Sql()
	fmt.Println(sql)

	rt1 := NewRuleTable(_rule1)

	sql = rt1.Sql()
	fmt.Println(sql)

	dr := NewRuleTable(_dataRule)

	sql = dr.Sql()
	fmt.Println(sql)

	dr1 := NewRuleTable(_dataRule1)

	sql = dr1.Sql()
	fmt.Println(sql)
}

func TestTag(t *testing.T) {
	tr := NewTagDataRuleTable(_tagRule).
		WithCommand(func(command *tablestructure.Command) {
			command.AddString("with")
		})

	sql := tr.Sql()
	fmt.Println(sql)

	tr1 := NewTagDataRuleTable(_tagRule1)
	sql = tr1.Sql()
	fmt.Println(sql)

	trs := NewTagDataRulesTable(_tagRules).
		WithCommand(func(command *tablestructure.Command) {
			command.AddString("with")
		})

	sql = trs.Sql()
	fmt.Println(sql)

	trs1 := NewTagDataRulesTable(_tagRules1)
	sql = trs1.Sql()
	fmt.Println(sql)
}

func TestDataContent(t *testing.T) {
	dcSeg := NewDataContentSegWordsTable(_dcSegWords).
		WithCommand(func(command *tablestructure.Command) {
			command.AddString("with")
		})

	sql := dcSeg.Sql()
	fmt.Println(sql)

	dcSplit := NewDataContentSpiltWordsTable(_dcSpiltWords).
		WithCommand(func(command *tablestructure.Command) {
			command.AddString("with")
		})

	sql = dcSplit.Sql()
	fmt.Println(sql)
}

func init() {
	_raw = model.NewRaw(_tableName, nil)
	_data = model.NewData(_tableName, _tableIdName, nil)
	_rule = model.NewRule(_ruleName, 20, 20, nil, nil)
	_rule1 = model.NewRule(_ruleName, 20, 20, map[string]int{"r1": 10, "r2": 30}, nil)
	_dataRule = model.NewDataRule(_data, _rule)
	_dataRule1 = model.NewDataRule(_data, _rule1)
	_tagRule = model.NewTagDataRule(_data, _rule, nil)
	_tagRule1 = model.NewTagDataRule(_data, _rule1, nil)
	_tagRules = model.NewTagDataRules("abc", _data, []assistant.Ruler{_rule, _rule1}, nil)
	_tagRules1 = model.NewTagDataRules("abc", _data, []assistant.Ruler{_rule1, _rule}, nil)
	_dcSegWords = model.NewDataContentSegWords(_data, "abc", nil)
	_dcSpiltWords = model.NewDataContentSpiltWords(_data, "abc", nil)
}
