package tager

import (
	"testing"
)

func TestTagMeans(t *testing.T) {
	tm := &TagMeans{}
	tm.prepare(ruleName, db)

	keys := tm.GetKeys()
	if len(keys) < 4 {
		t.Error("error")
	}
}

func TestTagKeyMeans(t *testing.T) {
	tm := NewTagKeyMeans(ruleName, db)
	resSlice := tm.Insert(contents)
	if len(resSlice) <= 0 {
		t.Error("error")
	}
}

func TestTagMostKeyMeans(t *testing.T) {
	tm := NewTagMostKeyMeans(ruleName, db)
	resSlice := tm.Insert(contents)
	if len(resSlice) <= 0 {
		t.Error("error")
	}

	resMap := tm.Update(contents)
	if len(resMap) <= 0 {
		t.Error("error")
	}
}

func TestTagMostTextMeans(t *testing.T) {
	tm := NewTagMostTextMeans(ruleName, db)
	results := tm.Insert(contents)
	if len(results) <= 0 {
		t.Error("error")
	}

	resMap := tm.Update(contents)
	if len(resMap) <= 0 {
		t.Error("error")
	}
}
