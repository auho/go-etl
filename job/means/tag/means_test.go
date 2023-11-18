package tag

import (
	"fmt"
	"strings"
	"testing"
)

func TestMeans(t *testing.T) {
	tm := NewMeans(_rule, nil)

	keys := tm.GetKeys()
	if len(keys) < 3 {
		t.Error("error")
	}
}

func TestKey(t *testing.T) {
	tm := NewKey(_rule)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	resSlice := tm.Insert(_contents)
	if len(resSlice) <= 0 {
		t.Error("error")
	}

	fmt.Println(resSlice)
}

func TestLabel(t *testing.T) {
	tm := NewLabel(_rule)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	resSlice := tm.Insert(_contents)
	if len(resSlice) <= 0 {
		t.Error("error")
	}

	fmt.Println(resSlice)
}

func TestMostKey(t *testing.T) {
	tm := NewMostKey(_rule)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	resSlice := tm.Insert(_contents)
	if len(resSlice) <= 0 {
		t.Error("error")
	}

	fmt.Println(resSlice)

	keys := tm.GetKeys()
	for _, k := range keys {
		if k != "a" && k != "ab" && k != "a_keyword" && k != "a_keyword_num" {
			t.Error("keys error")
		}
	}

	resMap := tm.Update(_contents)
	if len(resMap) <= 0 {
		t.Error("error")
	}

	fmt.Println(resMap)

	for _, k := range keys {
		if _, ok := resMap[k]; !ok {
			t.Error("update error")
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
			t.Error("alias is error")
		}
	}

	resSlice := tm.Insert(_contents)
	if len(resSlice) <= 0 {
		t.Error("error")
	}

	fmt.Println(resSlice)

	keys := tm.GetKeys()
	for _, k := range keys {
		if !strings.HasSuffix(k, "_alias") {
			t.Error("keys error")
		}
	}

	resMap := tm.Update(_contents)
	if len(resMap) <= 0 {
		t.Error("error")
	}

	fmt.Println(resMap)

	for _, k := range keys {
		if _, ok := resMap[k]; !ok {
			t.Error("update error")
		}
	}

}
