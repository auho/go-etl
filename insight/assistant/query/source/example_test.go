package source

import (
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
)

func ExampleNewTable() {
	_ = NewTable(Source{
		Name:  "name",
		DB:    nil,
		Table: dml.NewTable("table_name").Select([]string{"field1", "field2"}),
	})
}

func ExampleNewPlaceholder() {
	// one two

	// WithItems
	_ = NewPlaceholder(Source{
		Name: "name",
		DB:   nil,
		Table: dml.NewTable("table_name").
			Select([]string{"field1", "field2"}).
			Where("`field1` = '##two##' AND `field2` = '##one##'"),
	}).WithItems([]map[string]any{
		{"one": "1", "two": "2"},
		{"one": "1", "two": "21"},
		{"one": "11", "two": "2"},
		{"one": "11", "two": "21"},
	})

	// WithItemsCross
	_ = NewPlaceholder(Source{
		Name: "name",
		DB:   nil,
		Table: dml.NewTable("table_name").
			Select([]string{"field1", "field2"}).
			Where("`field1` = '##two##' AND `field2` = '##one##' AND `field3` = ##three##"),
	}).WithItemsCross(map[string][]any{
		"one": {"1", "11"},
		"two": {"2", "21"},
	})
}

func ExampleNewPlaceholderStack() {
	// one two three

	// WithCategories
	// WithStacks
	_ = NewPlaceholderStack(Source{
		Name: "name",
		DB:   nil,
		Table: dml.NewTable("table_name").
			Select([]string{"field1", "field2"}).
			Where("`field1` = '##two##' AND `field2` = '##one##' AND `field3` = ##three##"),
	}).WithCategories([]map[string]any{
		{"three": 1},
		{"three": 2},
	}).WithStacks([]map[string]any{
		{"one": "1", "two": "2"},
		{"one": "1", "two": "21"},
		{"one": "11", "two": "2"},
		{"one": "11", "two": "21"},
	})

	// WithCategoriesCross
	// WithStacksCross
	_ = NewPlaceholderStack(Source{
		Name: "name",
		DB:   nil,
		Table: dml.NewTable("table_name").
			Select([]string{"field1", "field2"}).
			Where("`field1` = '##two##' AND `field2` = '##one##' AND `field3` = ##three## AND `field4`= ##four##"),
	}).WithCategoriesCross(map[string][]any{
		"three": {1, 2},
		"four":  {3, 4},
	}).WithStacksCross(map[string][]any{
		"one": {"1", "11"},
		"two": {"2", "21"},
	})
}
