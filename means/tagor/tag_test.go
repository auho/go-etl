package tagor

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestTagMatcher_GetResultInsertKeys(t *testing.T) {
	tm := newTagMatcher(ruleName, db)
	keys := tm.getResultInsertKeys()
	if len(keys) < 4 {
		t.Error("error")
	}
}

func TestMatcher(t *testing.T) {
	items := make([]map[string]string, 0)
	items = append(items, map[string]string{"a": "b", "b": "b", "c": "c"})
	items = append(items, map[string]string{"a": "123", "b": "b", "c": "c"})
	items = append(items, map[string]string{"a": "中文", "b": "b2", "c": "c2"})
	items = append(items, map[string]string{"a": "中_文", "b": "b3", "c": "c3"})
	items = append(items, map[string]string{"a": "_中1_a文", "b": "b4", "c": "c4"})
	items = append(items, map[string]string{"a": "，。【】", "b": "b5", "c": "c5"})
	items = append(items, map[string]string{"a": ".+*?()|[]{}^$`))", "b": "b6", "c": "c6"})

	m := NewMatcher(WithTagMatcherKeyFunc(func(s string) string {
		res, err := regexp.MatchString(`^[\w+._\s()]+$`, s)
		if err != nil {
			return s
		}

		if res {
			return fmt.Sprintf(`\b%s\b`, s)
		} else {
			return strings.ReplaceAll(s, "_", `.{1,3}`)
		}
	}))

	m.prepare("a", items)

	var results []*Result
	var tagResults []*LabelResult

	fmt.Println("\n Match")
	results = m.Match(contents)
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("\n MatchText")
	results = m.MatchText(contents)
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("\n MatchKey")
	results = m.MatchKey(contents)
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("\n MatchFirstText")
	results = m.MatchFirstText(contents)
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("\n MatchLastText")
	results = m.MatchLastText(contents)
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("\n MatchMostKey")
	results = m.MatchMostKey(contents)
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("\n MatchMostText")
	results = m.MatchMostText(contents)
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("\n MatchLabel")
	tagResults = m.MatchLabel(contents)
	for _, tagResult := range tagResults {
		fmt.Println(tagResult)
	}

	fmt.Println("\n MatchLabelMostText")
	tagResults = m.MatchLabelMostText(contents)
	for _, tagResult := range tagResults {
		fmt.Println(tagResult)
	}
}
