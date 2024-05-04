package condition

import (
	"github.com/auho/go-etl/v2/job/explore/collect"
	"github.com/auho/go-etl/v2/job/explore/search"
)

func ExampleNewCondition() {
	var _search search.Searcher
	opt := NewCondition(collect.NewKeys([]string{"a"}), _search)

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
	).ToOperation()

	_ = opt
}

func ExampleNewOR() {
	opt := NewOR(
		func(m map[string]any) bool {
			return m["a"] == 1
		}, func(m map[string]any) bool {
			return m["b"] == 2
		},
	).ToOperation()

	_ = opt
}
