package tag

import (
	"strings"
	"testing"
)

func TestTagMeans(t *testing.T) {
	tm := NewTagMeans(_rule, nil)

	keys := tm.GetKeys()
	if len(keys) < 4 {
		t.Error("error")
	}
}

func TestTagKeyMeans(t *testing.T) {
	tm := NewKey(_rule)
	resSlice := tm.Insert(_contents)
	if len(resSlice) <= 0 {
		t.Error("error")
	}
}

func TestTagMostKeyMeans(t *testing.T) {
	tm := NewMostKey(_rule)
	resSlice := tm.Insert(_contents)
	if len(resSlice) <= 0 {
		t.Error("error")
	}

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

	for _, k := range keys {
		if _, ok := resMap[k]; !ok {
			t.Error("update error")
		}
	}
}

func TestTagMostTextMeans(t *testing.T) {
	tm := NewMostText(_ruleAliasFixed)

	for _, k := range tm.rule.LabelsAlias() {
		if !strings.HasSuffix(k, "_alias") {
			t.Error("alias is error")
		}
	}

	results := tm.Insert(_contents)
	if len(results) <= 0 {
		t.Error("error")
	}

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

	for _, k := range keys {
		if _, ok := resMap[k]; !ok {
			t.Error("update error")
		}
	}
}
