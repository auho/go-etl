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
	fmt.Println(resSlice)
	if len(resSlice) != 1 {
		t.Fatal("error")
	}
}

func TestKey(t *testing.T) {
	tm := NewKey(_rule)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	resSlice := tm.Insert(_contents)
	fmt.Println(resSlice)
	if len(resSlice) != 5 {
		t.Fatal("error")
	}

}

func TestLabel(t *testing.T) {
	tm := NewLabel(_rule)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	resSlice := tm.Insert(_contents)
	fmt.Println(resSlice)

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
	fmt.Println(resSlice)
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
	fmt.Println(resSlice)
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
	fmt.Println(resSlice)
	if len(resSlice) != 1 {
		t.Fatal("error")
	}

	if resSlice[0][_rule.KeywordName()] != "b" {
		t.Fatal("error")
	}
}
