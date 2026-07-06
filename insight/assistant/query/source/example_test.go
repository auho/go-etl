package source

import (
	"github.com/auho/go-etl/v3/insight/assistant/sqlbuilder/dml"
)

func ExampleNewRows() {
	_, _ = NewRows(Base{}).
		Dataset()

	_ = NewRows(Base{
		Name:  "name",
		DB:    nil,
		Table: dml.NewTable("table_name").Select([]string{"field1", "field2"}),
	})
}

func ExampleNewPlaceholder() {
	_, _ = NewPlaceholder(Base{}).
		WithItems(nil).
		WithItemsCross(nil).
		Dataset()

	// one two

	// WithItems
	_ = NewPlaceholder(Base{
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
	_ = NewPlaceholder(Base{
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
	_, _ = NewPlaceholderStack(Base{}).
		WithCategories(nil).
		WithCategoriesCross(nil).
		WithStacks(nil).
		WithStacksCross(nil).
		Dataset()

	// one two three

	// WithCategories
	// WithStacks
	_ = NewPlaceholderStack(Base{
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
	_ = NewPlaceholderStack(Base{
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
