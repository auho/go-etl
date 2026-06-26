package filter

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collect"
)

func ExampleNewMatcher() {
	var _search extract.Extractor
	opt := NewMatcher(collect.NewKeys([]string{"a"}), _search)

	_ = opt
}

func ExampleNewContainsAll() {
	opt := NewContainsAll("a", []string{"a1", "a2"})

	_ = opt
}

func ExampleNewContainsAny() {
	opt := NewContainsAny("a", []string{"a1", "a2"})

	_ = opt
}

func ExampleNewAND() {
	opt := NewAND(
		func(m map[string]any) bool {
			return m["a"] == 1
		}, func(m map[string]any) bool {
			return m["b"] == 2
		},
	).ToPredicate()

	_ = opt
}

func ExampleNewOR() {
	opt := NewOR(
		func(m map[string]any) bool {
			return m["a"] == 1
		}, func(m map[string]any) bool {
			return m["b"] == 2
		},
	).ToPredicate()

	_ = opt
}
