package create

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/entity"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
)

var _tableName = "table1"
var _tableIdName = "tid"
var _ruleName = "rule1"
var _raw *entity.Raw
var _data *entity.Data
var _rule *entity.Rule
var _rule1 *entity.Rule
var _dataRule *entity.DataRule
var _dataRule1 *entity.DataRule
var _tagRule *entity.TagDataRule
var _tagRule1 *entity.TagDataRule
var _tagRules *entity.TagDataRules
var _tagRules1 *entity.TagDataRules
var _dcSegWords *entity.DataContentSegWords
var _dcSplitWords *entity.DataContentSplitWords

func TestRaw(t *testing.T) {
	rt := NewRawTable(_raw)
	rt.AddPKBigInt("id")
	rt.WithCommand(func(command *schema.Command) {
		command.AddString("with")
	})

	sql := rt.SQL()
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
	dt.WithCommand(func(command *schema.Command) {
		command.AddString("with")
	})

	sql := dt.SQL()
	fmt.Println(sql)
}

func TestRule(t *testing.T) {
	rt := NewRuleTable(_rule).
		WithCommand(func(command *schema.Command) {
			command.AddString("with")
		})

	sql := rt.SQL()
	fmt.Println(sql)

	rt1 := NewRuleTable(_rule1)

	sql = rt1.SQL()
	fmt.Println(sql)

	dr := NewRuleTable(_dataRule)

	sql = dr.SQL()
	fmt.Println(sql)

	dr1 := NewRuleTable(_dataRule1)

	sql = dr1.SQL()
	fmt.Println(sql)
}

func TestTag(t *testing.T) {
	tr := NewTagDataRuleTable(_tagRule).
		WithCommand(func(command *schema.Command) {
			command.AddString("with")
		})

	sql := tr.SQL()
	fmt.Println(sql)

	tr1 := NewTagDataRuleTable(_tagRule1)
	sql = tr1.SQL()
	fmt.Println(sql)

	trs := NewTagDataRulesTable(_tagRules).
		WithCommand(func(command *schema.Command) {
			command.AddString("with")
		})

	sql = trs.SQL()
	fmt.Println(sql)

	trs1 := NewTagDataRulesTable(_tagRules1)
	sql = trs1.SQL()
	fmt.Println(sql)
}

func TestDataContent(t *testing.T) {
	dcSeg := NewDataContentSegWordsTable(_dcSegWords).
		WithCommand(func(command *schema.Command) {
			command.AddString("with")
		})

	sql := dcSeg.SQL()
	fmt.Println(sql)

	dcSplit := NewDataContentSplitWordsTable(_dcSplitWords).
		WithCommand(func(command *schema.Command) {
			command.AddString("with")
		})

	sql = dcSplit.SQL()
	fmt.Println(sql)
}

func init() {
	_raw = entity.NewRaw(_tableName, nil)
	_data = entity.NewData(_tableName, _tableIdName, nil)
	_rule = entity.NewRule(_ruleName, 20, 20, nil, nil)
	_rule1 = entity.NewRule(_ruleName, 20, 20, map[string]int{"r1": 10, "r2": 30}, nil)
	_dataRule = entity.NewDataRule(_data, _rule)
	_dataRule1 = entity.NewDataRule(_data, _rule1)
	_tagRule = entity.NewTagDataRule(_data, _rule, nil)
	_tagRule1 = entity.NewTagDataRule(_data, _rule1, nil)
	_tagRules = entity.NewTagDataRules("abc", _data, []assistant.Rule{_rule, _rule1}, nil)
	_tagRules1 = entity.NewTagDataRules("abc", _data, []assistant.Rule{_rule1, _rule}, nil)
	_dcSegWords = entity.NewDataContentSegWords(_data, "abc", nil)
	_dcSplitWords = entity.NewDataContentSplitWords(_data, "abc", nil)
}
