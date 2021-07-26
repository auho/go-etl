package means

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestMatcher(t *testing.T) {
	items := make([]map[string]string, 0)
	items = append(items, map[string]string{"a": "b", "b": "b", "c": "c"})
	items = append(items, map[string]string{"a": "123", "b": "b", "c": "c"})
	items = append(items, map[string]string{"a": "中文", "b": "b2", "c": "c2"})
	items = append(items, map[string]string{"a": "中_文", "b": "b3", "c": "c3"})
	items = append(items, map[string]string{"a": "_中1_a文", "b": "b4", "c": "c4"})
	items = append(items, map[string]string{"a": "，。【】", "b": "b5", "c": "c5"})
	items = append(items, map[string]string{"a": ".+*?()|[]{}^$`))", "b": "b6", "c": "c6"})

	contents := []string{
		`b一ab一bc一abc一123b一b123一123一0123一1234一01234一`,
		`中文一b中文123一123中文b一中bb文一中123文一中00文一中aa文一中00文一中aa文一中中文文一中二二文一
123一一`,
	}
	tm := NewMatcher(WithTagMatcherKeyFun(func(s string) string {
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

	tm.init("a", items)

	var results []*Result
	var result *Result
	var tagResults []*TagResult
	var tagResult *TagResult

	fmt.Println("\n Match")
	results = tm.Match(contents)
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("\n MatchText")
	results = tm.MatchText(contents)
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("\n MatchKey")
	results = tm.MatchKey(contents)
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("\n MatchFirstText")
	result = tm.MatchFirstText(contents)
	fmt.Println(result)

	fmt.Println("\n MatchLastText")
	result = tm.MatchLastText(contents)
	fmt.Println(result)

	fmt.Println("\n MatchMostKey")
	result = tm.MatchMostKey(contents)
	fmt.Println(result)

	fmt.Println("\n MatchMostText")
	result = tm.MatchMostText(contents)
	fmt.Println(result)

	fmt.Println("\n MatchTag")
	tagResults = tm.MatchTag(contents)
	for _, tagResult := range tagResults {
		fmt.Println(tagResult)
	}

	fmt.Println("\n MatchTagMostText")
	tagResult = tm.MatchTagMostText(contents)
	fmt.Println(tagResult)

}
