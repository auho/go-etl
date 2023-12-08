package tag

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestMatcher(t *testing.T) {
	// keyword: a
	// labels: b c
	items := []map[string]string{
		{"a": "b", "b": "b", "c": "c"},
		{"a": "123", "b": "b", "c": "c"},
		{"a": "中文", "b": "b2", "c": "c2"},
		{"a": "中_文", "b": "b3", "c": "c3"},
		{"a": "_中1_a文", "b": "b4", "c": "c4"},
		{"a": "，。【】", "b": "b5", "c": "c5"},
		{"a": ".+*?()|[]{}^$`))", "b": "b6", "c": "c6"},
	}

	m := NewMatcher(WithMatcherKeyFormatFunc(func(s string) string {
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

	var results []Result
	var tagResults []LabelResult

	fmt.Println("\n Match")
	results = m.Match(_contents)
	for _, result := range results {
		fmt.Println(result)
	}

	if len(results) != 31 || results[0].Key != "b" {
		t.Fatal("Match")
	}

	fmt.Println("\n MatchText")
	results = m.MatchText(_contents)
	for _, result := range results {
		fmt.Println(result)
	}

	if len(results) != 17 || results[1].Key != "123" || results[1].Texts["123"] != 7 || results[2].Num != 4 {
		t.Fatal("MatchText")
	}

	fmt.Println("\n MatchKey")
	results = m.MatchKey(_contents)
	for _, result := range results {
		fmt.Println(result)
	}

	if len(results) != 4 || results[1].Texts["123"] != 7 || results[2].Key != "中文" || results[3].Num != 17 {
		t.Fatal("MatchKey")
	}

	fmt.Println("\n MatchFirstText")
	results = m.MatchFirstText(_contents)
	for _, result := range results {
		fmt.Println(result)
	}

	if len(results) != 1 || results[0].Key != "b" || results[0].Num != 1 {
		t.Fatal("MatchFirstText")
	}

	fmt.Println("\n MatchLastText")
	results = m.MatchLastText(_contents)
	for _, result := range results {
		fmt.Println(result)
	}

	if len(results) != 1 || results[0].Key != "中_文" || results[0].Num != 1 || results[0].Texts["中123文"] != 1 {
		t.Fatal("MatchLastText")
	}

	fmt.Println("\n MatchMostKey")
	results = m.MatchMostKey(_contents)
	for _, result := range results {
		fmt.Println(result)
	}

	if len(results) != 1 || results[0].Key != "中_文" || results[0].Num != 17 || results[0].Texts["中00文"] != 2 {
		t.Fatal("MatchMostKey")
	}

	fmt.Println("\n MatchMostText")
	results = m.MatchMostText(_contents)
	for _, result := range results {
		fmt.Println(result)
	}

	if len(results) != 1 || results[0].Key != "123" || results[0].Num != 7 || results[0].Texts["123"] != 7 {
		t.Fatal("MatchMostText")
	}

	fmt.Println("\n MatchLabel")
	tagResults = m.MatchLabel(_contents)
	for _, tagResult := range tagResults {
		fmt.Println(tagResult)
	}

	if len(tagResults) != 3 || tagResults[1].MatchAmount != 10 || tagResults[1].Identity != "-b-c" {
		t.Fatal("MatchLabel")
	}

	fmt.Println("\n MatchLabelMostText")
	tagResults = m.MatchLabelMostText(_contents)
	for _, tagResult := range tagResults {
		fmt.Println(tagResult)
	}

	if len(tagResults) != 1 || tagResults[0].MatchAmount != 17 || tagResults[0].Identity != "-b3-c3" {
		t.Fatal("MatchLabelMostText")
	}
}
