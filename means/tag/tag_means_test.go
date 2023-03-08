package tag

import (
	"testing"

	"github.com/auho/go-etl/v2/means/tag/rule"
)

func TestTagMeans(t *testing.T) {
	tm := NewTagMeans(ruleName, nil, WithDBRule(db, rule.WithDBRuleDataName("data")))

	keys := tm.GetKeys()
	if len(keys) < 4 {
		t.Error("error")
	}
}

func TestTagKeyMeans(t *testing.T) {
	tm := NewKey(ruleName, WithDBRule(db))
	resSlice := tm.Insert(contents)
	if len(resSlice) <= 0 {
		t.Error("error")
	}
}

func TestTagMostKeyMeans(t *testing.T) {
	tm := NewMostKey(ruleName, WithDBRule(db))
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
	tm := NewMostText(ruleName,
		WithAlias(map[string]string{
			"a":             "a_alias",
			"ab":            "ab_alias",
			"a_keyword":     "a_keyword_alias",
			"a_keyword_num": "a_keyword_num_alias",
		}),
		WithFixedTags(map[string]string{
			"c": "c_fixed",
			"d": "d_fixed",
		}),
		WithDBRule(
			db,
			rule.WithDBRuleTagsName([]string{"a"}),
			rule.WithDBRuleDataName("data"),
		),
	)

	for _, k := range tm.tagMatcher.tagsName {
		if k != "a_alias" && k != "ab_alias" {
			t.Error("alias is error")
		}
	}

	results := tm.Insert(contents)
	if len(results) <= 0 {
		t.Error("error")
	}

	keys := tm.GetKeys()
	for _, k := range keys {
		if k != "a_alias" && k != "ab_alias" && k != "a_keyword_alias" && k != "a_keyword_num_alias" && k != "c_alias" && k != "d_alias" {
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
