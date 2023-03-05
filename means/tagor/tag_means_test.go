package tagor

import (
	"testing"
)

func TestTagMeans(t *testing.T) {
	tm := NewTagMeans(ruleName, db, nil, WithTagMatcherShortTableName("data_a"))

	if tm.tagMatcher.tableName != dataRuleTableName {
		t.Error("table name is error")
	}

	keys := tm.GetKeys()
	if len(keys) < 4 {
		t.Error("error")
	}
}

func TestTagKeyMeans(t *testing.T) {
	tm := NewKey(ruleName, db)
	resSlice := tm.Insert(contents)
	if len(resSlice) <= 0 {
		t.Error("error")
	}
}

func TestTagMostKeyMeans(t *testing.T) {
	tm := NewMostKey(ruleName, db)
	resSlice := tm.Insert(contents)
	if len(resSlice) <= 0 {
		t.Error("error")
	}

	keys := tm.GetKeys()
	for _, k := range keys {
		if k != "a" && k != "ab" && k != "a_keyword" && k != "a_keyword_num" {
			t.Error("keys error")
		}
	}

	resMap := tm.Update(contents)
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
	tm := NewMostText(ruleName, db,
		WithTagMatcherTags([]string{"a"}),
		WithTagMatcherDataName("data"),
		WithTagMatcherAlias(map[string]string{
			"a":             "aa",
			"ab":            "aabb",
			"a_keyword":     "aa_keyword",
			"a_keyword_num": "aa_keyword_num",
		}),
		WithTagMatcherFixedTags(map[string]string{
			"c": "c1",
			"d": "d1",
		}),
	)

	if tm.tagMatcher.tableName != dataRuleTableName {
		t.Error("table name error")
	}

	for _, k := range tm.tagMatcher.tagsName {
		if k != "aa" && k != "aabb" {
			t.Error("alias is error")
		}
	}

	results := tm.Insert(contents)
	if len(results) <= 0 {
		t.Error("error")
	}

	keys := tm.GetKeys()
	for _, k := range keys {
		if k != "aa" && k != "aabb" && k != "aa_keyword" && k != "aa_keyword_num" && k != "c" && k != "d" {
			t.Error("keys error")
		}
	}

	resMap := tm.Update(contents)
	if len(resMap) <= 0 {
		t.Error("error")
	}

	for _, k := range keys {
		if _, ok := resMap[k]; !ok {
			t.Error("update error")
		}
	}
}
