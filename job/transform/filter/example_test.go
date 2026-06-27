package filter

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collect"
)

func ExampleNewFilterPredicate() {
	var _search extract.Extractor
	opt := NewFilterPredicate(collect.NewKeysAll([]string{"a"}), _search)

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

func ExampleNewAnd() {
	opt := NewAnd(
		func(m map[string]any) (bool, error) {
			return m["a"] == 1, nil
		}, func(m map[string]any) (bool, error) {
			return m["b"] == 2, nil
		},
	).ToPredicate()

	_ = opt
}

func ExampleNewOr() {
	opt := NewOr(
		func(m map[string]any) (bool, error) {
			return m["a"] == 1, nil
		}, func(m map[string]any) (bool, error) {
			return m["b"] == 2, nil
		},
	).ToPredicate()

	_ = opt
}
