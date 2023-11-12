package source

import (
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
)

func ExampleNewTable() {
	_ = NewTable(
		Source{
			Name:  "name",
			DB:    nil,
			Table: dml.NewTable("table_name").Select([]string{"field1", "field2"}),
		},
	)
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
	}).WithItems(
		[]map[string]string{
			{"one": "1", "two": "2"},
			{"one": "1", "two": "21"},
			{"one": "11", "two": "2"},
			{"one": "11", "two": "21"},
		},
	)

	// WithItemsCross
	_ = NewPlaceholder(Source{
		Name: "two",
		DB:   nil,
		Table: dml.NewTable("table_name").
			Select([]string{"field1", "field2"}).
			Where("`field1` = '##two##' AND `field2` = '##one##' AND `field3` = ##three##"),
	}).WithItemsCross(
		[]map[string][]string{
			{"one": []string{"1", "11"}},
			{"two": []string{"2", "21"}},
		},
	)
}

func ExampleNewPlaceholderStack() {
	// one two three

	// WithCategories
	// WithStacks
	_ = NewPlaceholderStack(Source{
		Name: "four",
		DB:   nil,
		Table: dml.NewTable("table_name").
			Select([]string{"field1", "field2"}).
			Where("`field1` = '##two##' AND `field2` = '##one##' AND `field3` = ##three##"),
	}).WithCategories([]map[string]string{
		{"three": "1"},
		{"three": "2"},
	}).WithStacks([]map[string]string{
		{"one": "1", "two": "2"},
		{"one": "1", "two": "21"},
		{"one": "11", "two": "2"},
		{"one": "11", "two": "21"},
	})

	// WithCategoriesCross
	// WithStacksCross
	_ = NewPlaceholderStack(Source{
		Name: "four",
		DB:   nil,
		Table: dml.NewTable("table_name").
			Select([]string{"field1", "field2"}).
			Where("`field1` = '##two##' AND `field2` = '##one##' AND `field3` = ##three##"),
	}).WithCategoriesCross([]map[string][]string{
		{"three": []string{"1", "2"}},
	}).WithStacksCross(
		[]map[string][]string{
			{"one": []string{"1", "11"}},
			{"two": []string{"2", "21"}},
		},
	)
}
