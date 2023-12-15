package match

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestMeans(t *testing.T) {
	tm := NewMatch(_rule, nil)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	keys := tm.GetKeys()
	if len(keys) < 3 {
		t.Fatal()
	}
}

func TestWholeLabels(t *testing.T) {
	tm := NewWholeLabels(_rule).WithFuzzy(FuzzyConfig{
		Window: 3,
		Sep:    "_",
	}).WithDebug()
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	resSlice := tm.Insert(_contents)
	_printResult(resSlice)
	if len(resSlice) != 1 {
		t.Fatal()
	}

	if resSlice[0][_rule.NameAlias()] != "123|a|ab|中1文|中文" || resSlice[0][_rule.LabelNumNameAlias()] != 5 ||
		resSlice[0][_rule.KeywordNumNameAlias()] != 6 || resSlice[0][_rule.KeywordAmountNameAlias()] != 33 {
		t.Fatal()
	}

	if resSlice[0][_rule.LabelNumNameAlias()] != len(strings.Split(resSlice[0]["a"].(string), "|")) {
		t.Fatal("label num")
	}

	amount := 0
	keyword := resSlice[0][_rule.KeywordNameAlias()].(string)
	for _, _s := range strings.Split(keyword, "|") {
		amount += len(strings.Split(_s, ","))
	}

	if resSlice[0][_rule.KeywordNumNameAlias()] != amount {
		t.Fatal("keyword num")
	}
}

func TestLabel(t *testing.T) {
	tm := NewLabel(_rule)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	resSlice := tm.Insert(_contents)
	_printResult(resSlice)
	if len(resSlice) != 5 {
		t.Fatal()
	}

	if resSlice[0]["a"] != "a" || resSlice[0][_rule.KeywordAmountNameAlias()] != 4 {
		t.Fatal()
	}

	amount := 0
	for _, _kw := range strings.Split(resSlice[0][_rule.KeywordNameAlias()].(string), ",") {
		for _, _s := range strings.Split(_kw, " ") {
			_n, _ := strconv.Atoi(_s)
			amount += _n
		}
	}
	if amount != resSlice[0][_rule.KeywordAmountNameAlias()] {
		t.Fatal()
	}

	if resSlice[3]["a"] != "中文" || resSlice[3][_rule.KeywordNameAlias()] == "中文" || resSlice[3][_rule.KeywordAmountNameAlias()] != 4 {
		t.Fatal()
	}

	if resSlice[4]["a"] != "中1文" || resSlice[4][_rule.KeywordAmountNameAlias()] != 17 {
		t.Fatal()
	}

	amount = 0
	for _, _m := range resSlice {
		amount += _m[_rule.KeywordAmountNameAlias()].(int)
	}
	if amount != 33 {
		t.Fatal("amount")
	}
}

func TestKey(t *testing.T) {
	tm := NewKey(_rule)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	resSlice := tm.Insert(_contents)
	_printResult(resSlice)
	if len(resSlice) != 6 {
		t.Fatal()
	}

	if resSlice[0]["a"] != "a" || resSlice[0][_rule.KeywordNameAlias()] != "b" || resSlice[0][_rule.KeywordAmountNameAlias()] != 3 {
		t.Fatal()
	}

	if resSlice[2]["a"] != "123" || resSlice[2][_rule.KeywordNameAlias()] != "123" || resSlice[2][_rule.KeywordAmountNameAlias()] != 7 {
		t.Fatal()
	}

	if resSlice[3]["a"] != "a" || resSlice[3][_rule.KeywordNameAlias()] != "a" || resSlice[3][_rule.KeywordAmountNameAlias()] != 1 {
		t.Fatal()
	}

	if resSlice[5]["a"] != "中1文" || resSlice[5][_rule.KeywordNameAlias()] != "中_文" || resSlice[5][_rule.KeywordAmountNameAlias()] != 17 {
		t.Fatal()
	}

	_amount := 0
	for _, _m := range resSlice {
		_amount += _m[_rule.KeywordAmountNameAlias()].(int)
	}

	if _amount != 33 {
		t.Fatal("amount")
	}
}

func TestMostKey(t *testing.T) {
	tm := NewMostKey(_rule)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	resSlice := tm.Insert(_contents)
	_printResult(resSlice)
	if len(resSlice) != 1 {
		t.Fatal()
	}

	if resSlice[0][_rule.NameAlias()] != "中1文" || resSlice[0][_rule.KeywordName()] != "中_文" || resSlice[0][_rule.KeywordAmountNameAlias()] != 17 {
		t.Fatal()
	}

	keys := tm.GetKeys()

	resMap := tm.Update(_contents)
	fmt.Println(resMap)
	if len(resMap) <= 0 {
		t.Fatal()
	}

	for _, k := range keys {
		if _, ok := resMap[k]; !ok {
			t.Fatal()
		}
	}
}

func TestMostText(t *testing.T) {
	tm := NewMostText(_ruleAliasFixed)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	for _, k := range _ruleAliasFixed.LabelsAlias() {
		if !strings.HasSuffix(k, "_alias") {
			t.Fatal()
		}
	}

	resSlice := tm.Insert(_contents)
	_printResult(resSlice)
	if len(resSlice) != 1 {
		t.Fatal()
	}

	if resSlice[0][_ruleAliasFixed.NameAlias()] != "123" || resSlice[0][_ruleAliasFixed.KeywordNameAlias()] != "123" ||
		resSlice[0][_ruleAliasFixed.KeywordAmountNameAlias()] != 7 || resSlice[0]["c_alias"] != "c_fixed" {
		t.Fatal()
	}

	keys := tm.GetKeys()
	for _, k := range keys {
		if !strings.HasSuffix(k, "_alias") {
			t.Error("keys error")
		}
	}

	resMap := tm.Update(_contents)
	fmt.Println(resMap)
	if len(resMap) <= 0 {
		t.Fatal()
	}

	for _, k := range keys {
		if _, ok := resMap[k]; !ok {
			t.Error("update error")
		}
	}
}

func TestFirst(t *testing.T) {
	tm := NewFirst(_rule)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	resSlice := tm.Insert(_contents)
	_printResult(resSlice)
	if len(resSlice) != 1 {
		t.Fatal()
	}

	if resSlice[0][_rule.NameAlias()] != "a" || resSlice[0][_rule.KeywordName()] != "b" || resSlice[0][_rule.KeywordAmountNameAlias()] != 1 {
		t.Fatal()
	}
}

func _printResult[T any](sm []T) {
	for _, m := range sm {
		fmt.Println(fmt.Sprintf("%+v", m))
	}
}
