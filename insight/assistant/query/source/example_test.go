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
		AppendItems(nil).
		AppendItemsCross(nil).
		Dataset()

	// one two

	// AppendItems
	_ = NewPlaceholder(Base{
		Name: "name",
		DB:   nil,
		Table: dml.NewTable("table_name").
			Select([]string{"field1", "field2"}).
			Where("`field1` = '##two##' AND `field2` = '##one##'"),
	}).AppendItems([]map[string]any{
		{"one": "1", "two": "2"},
		{"one": "1", "two": "21"},
		{"one": "11", "two": "2"},
		{"one": "11", "two": "21"},
	})

	// AppendItemsCross
	_ = NewPlaceholder(Base{
		Name: "name",
		DB:   nil,
		Table: dml.NewTable("table_name").
			Select([]string{"field1", "field2"}).
			Where("`field1` = '##two##' AND `field2` = '##one##' AND `field3` = ##three##"),
	}).AppendItemsCross(map[string][]any{
		"one": {"1", "11"},
		"two": {"2", "21"},
	})
}

func ExampleNewPlaceholderStack() {
	_, _ = NewPlaceholderStack(Base{}).
		AppendCategories(nil).
		AppendCategoriesCross(nil).
		AppendStacks(nil).
		AppendStacksCross(nil).
		Dataset()

	// one two three

	// AppendCategories
	// AppendStacks
	_ = NewPlaceholderStack(Base{
		Name: "name",
		DB:   nil,
		Table: dml.NewTable("table_name").
			Select([]string{"field1", "field2"}).
			Where("`field1` = '##two##' AND `field2` = '##one##' AND `field3` = ##three##"),
	}).AppendCategories([]map[string]any{
		{"three": 1},
		{"three": 2},
	}).AppendStacks([]map[string]any{
		{"one": "1", "two": "2"},
		{"one": "1", "two": "21"},
		{"one": "11", "two": "2"},
		{"one": "11", "two": "21"},
	})

	// AppendCategoriesCross
	// AppendStacksCross
	_ = NewPlaceholderStack(Base{
		Name: "name",
		DB:   nil,
		Table: dml.NewTable("table_name").
			Select([]string{"field1", "field2"}).
			Where("`field1` = '##two##' AND `field2` = '##one##' AND `field3` = ##three## AND `field4`= ##four##"),
	}).AppendCategoriesCross(map[string][]any{
		"three": {1, 2},
		"four":  {3, 4},
	}).AppendStacksCross(map[string][]any{
		"one": {"1", "11"},
		"two": {"2", "21"},
	})
}
