package filter

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

func ExampleNewFromCollector() {
	var _search extract.Extractor
	opt := NewFromCollector(collector.NewKeysAll([]string{"a"}, _search))

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
		Func(func(m map[string]any) (bool, error) {
			return m["a"] == 1, nil
		}),
		Func(func(m map[string]any) (bool, error) {
			return m["b"] == 2, nil
		}),
	)

	_ = opt
}

func ExampleNewOr() {
	opt := NewOr(
		Func(func(m map[string]any) (bool, error) {
			return m["a"] == 1, nil
		}),
		Func(func(m map[string]any) (bool, error) {
			return m["b"] == 2, nil
		}),
	)

	_ = opt
}
