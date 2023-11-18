package dml

import (
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml/command"
)

// Table

func ExampleTable_Select() {
	_ = NewTable("one").Select([]string{"a", "b"})
}

func ExampleTable_SelectAlias() {
	_ = NewTable("one").SelectAlias(map[string]string{
		"a": "a_alias",
		"b": "b_alias",
	})
}

func ExampleTable_Aggregation() {
	_ = NewTable("one").Aggregation(map[string]string{
		"COUNT(`a`)": "总数",
	})
}

func ExampleTable_Where() {
	_ = NewTable("one").Where("`a` = 1")
}

func ExampleTable_GroupBy() {
	_ = NewTable("one").GroupBy([]string{"a", "b"})
}

func ExampleTable_GroupByAlias() {
	_ = NewTable("one").GroupByAlias(map[string]string{
		"a": "a_alias",
		"b": "b_alias",
	})
}

func ExampleTable_OrderBy() {
	_ = NewTable("one").OrderBy(
		"a", command.SortDesc,
		"b", command.SortAsc,
	)
}

func ExampleTable_OrderByAsc() {
	_ = NewTable("one").OrderByAsc("a").OrderByAsc("b")
}

func ExampleTable_OrderByDesc() {
	_ = NewTable("one").OrderByDesc("a").OrderByDesc("b")
}

func ExampleTable_Limit() {
	_ = NewTable("one").Limit(0, 11)
}

func ExampleTable_LeftJoin() {
	_ = NewTable("one").LeftJoin([]string{"a"}, NewTable("two"), []string{"two_a"})
}

func ExampleTable_SetField() {
	_ = NewTable("one").SetField(map[string]string{
		"a": "b",
		"c": "d",
	})
}

func ExampleTable_SetExpression() {
	_ = NewTable("one").SetExpression(map[string]string{
		"a": "`b` + 1 ",
		"c": "`d` * 2",
	})
}

func ExampleTable_SetValue() {
	_ = NewTable("one").SetValue(map[string]any{
		"a": "abc",
		"b": 1,
		"c": 1.1,
	})
}

// TableJoin

func ExampleNewTableJoin() {
	_t1 := NewTable("one")
	_t2 := NewTable("two")

	_ = NewTableJoin().Table(_t1).LeftJoin(_t2, []string{"a", "b"}, nil, nil)
}

func ExampleTableJoin_LeftJoin() {
	_t1 := NewTable("one")
	_t2 := NewTable("two")

	// 省略 table fields
	_ = NewTableJoin().Table(_t1).LeftJoin(_t2, []string{"a", "b"}, nil, nil)

	// 省略 fields
	// 两者结果相同
	_ = NewTableJoin().Table(_t1).LeftJoin(_t2, []string{"a", "b"}, _t1, nil)
	_ = NewTableJoin().Table(_t1).LeftJoin(_t2, []string{"a", "b"}, _t1, []string{"a", "b"})

	// 不同的 fields
	_ = NewTableJoin().Table(_t1).LeftJoin(_t2, []string{"a", "b"}, _t1, []string{"c", "d"})
}
