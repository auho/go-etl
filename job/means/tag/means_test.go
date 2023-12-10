package tag

import (
	"fmt"
	"strings"
	"testing"
)

func TestMeans(t *testing.T) {
	tm := NewMeans(_rule, nil)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	keys := tm.GetKeys()
	if len(keys) < 3 {
		t.Fatal("error")
	}
}

func TestWholeLabels(t *testing.T) {
	tm := NewWholeLabels(_rule)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	resSlice := tm.Insert(_contents)
	_printResult(resSlice)
	if len(resSlice) != 1 {
		t.Fatal()
	}

	if resSlice[0]["a"] != "123|a|ab|中1文|中文" || resSlice[0]["a_label_num"] != 5 ||
		resSlice[0]["a_keyword_num"] != 6 || resSlice[0]["a_keyword_amount"] != 33 {
		t.Fatal(0)
	}

	if resSlice[0]["a_label_num"] != len(strings.Split(resSlice[0]["a"].(string), "|")) {
		t.Fatal("label num")
	}

	amount := 0
	keyword := resSlice[0]["a_keyword"].(string)
	for _, _s := range strings.Split(keyword, "|") {
		amount += len(strings.Split(_s, ","))
	}

	if resSlice[0]["a_keyword_num"] != amount {
		t.Fatal("keyword num")
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

	if resSlice[0]["a"] != "a" || resSlice[0]["a_keyword"] != "b" || resSlice[0]["a_keyword_amount"] != 3 {
		t.Fatal(0)
	}

	if resSlice[2]["a"] != "123" || resSlice[2]["a_keyword"] != "123" || resSlice[2]["a_keyword_amount"] != 7 {
		t.Fatal(2)
	}

	if resSlice[3]["a"] != "a" || resSlice[3]["a_keyword"] != "a" || resSlice[3]["a_keyword_amount"] != 1 {
		t.Fatal(3)
	}

	if resSlice[5]["a"] != "中1文" || resSlice[5]["a_keyword"] != "中_文" || resSlice[5]["a_keyword_amount"] != 17 {
		t.Fatal(5)
	}

	_amount := 0
	for _, _m := range resSlice {
		_amount += _m["a_keyword_amount"].(int)
	}

	if _amount != 33 {
		t.Fatal("amount")
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
		t.Fatal("error")
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
		t.Fatal("error")
	}

	if resSlice[0][_rule.KeywordName()] != "中_文" {
		t.Fatal("error")
	}

	keys := tm.GetKeys()

	resMap := tm.Update(_contents)
	fmt.Println(resMap)
	if len(resMap) <= 0 {
		t.Fatal("error")
	}

	for _, k := range keys {
		if _, ok := resMap[k]; !ok {
			t.Fatal("error")
		}
	}
}

func TestMostText(t *testing.T) {
	tm := NewMostText(_ruleAliasFixed)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	for _, k := range tm.rule.LabelsAlias() {
		if !strings.HasSuffix(k, "_alias") {
			t.Fatal("error")
		}
	}

	resSlice := tm.Insert(_contents)
	_printResult(resSlice)
	if len(resSlice) != 1 {
		t.Error("error")
	}

	if resSlice[0][_ruleAliasFixed.KeywordNameAlias()] != "123" {
		t.Fatal("error")
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
		t.Error("error")
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
		t.Fatal("error")
	}

	if resSlice[0][_rule.KeywordName()] != "b" {
		t.Fatal("error")
	}
}

func _printResult[T any](sm []T) {
	for _, m := range sm {
		fmt.Println(fmt.Sprintf("%+v", m))
	}
}
