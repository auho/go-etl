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

func ExampleNewContainAll() {
	opt := NewContainAll("a", []string{"a1", "a2"})

	_ = opt
}

func ExampleNewContainAny() {
	opt := NewContainAny("a", []string{"a1", "a2"})

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
