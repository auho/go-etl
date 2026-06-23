package dml

import (
	"fmt"
	"os"
	"testing"

	"github.com/auho/go-etl/v2/insight/assistant/sqlbuilder/dml/command"
)

func TestMain(t *testing.M) {
	code := t.Run()

	os.Exit(code)
}

func TestTable(t *testing.T) {
	t1 := getTable1()
	s1 := t1.SQL()
	fmt.Println(s1)

	t2 := NewSqlTable("efg", s1).Select([]string{"b11", "a11"})
	s2 := t2.SQL()
	fmt.Println(s2)

	fmt.Println(t1.GetSelectFields())
	fmt.Println(t2.GetSelectFields())
}

func TestTableJoin(t *testing.T) {
	t1 := getTable1()

	t2 := getTable2()

	t3 := NewTableJoin().Table(t1).LeftJoin(t2, []string{"a", "c"}, nil, nil).Limit(1, 11)
	s3 := t3.SQL()
	fmt.Println(s3)
	fmt.Println(t3.GetSelectFields())

	t4 := t1.CreateJoin().LeftJoin(t2, []string{"a", "c"}, nil, nil).Limit(2, 22)
	s4 := t4.SQL()
	fmt.Println(s4)
	fmt.Println(t4.GetSelectFields())
}

func TestInsert(t *testing.T) {
	t1 := getTable1()

	fmt.Println(t1.InsertSQL("i1"))

	fmt.Println(t1.InsertWithFieldsSQL("i2", []string{"a", "a11", "d11"}))

	fmt.Println(t1.InsertWithFieldsSQL("i2", nil))

	t2 := getTableJoin()

	fmt.Println(t2.InsertSQL("i1"))

	fmt.Println(t2.InsertWithFieldsSQL("i2", []string{"a", "a11", "d11"}))

	fmt.Println(t2.InsertWithFieldsSQL("i2", nil))
}

func TestUpdate(t *testing.T) {
	t1 := getTable1().
		SetField(map[string]string{"a": "b", "c": "d"}).
		SetExpression(map[string]string{"a": "`b` + 1 ", "c": "`d` * 2"}).
		SetValue(map[string]any{"e": "abc", "f": 1, "g": 1.11})
	fmt.Println(t1.UpdateSQL())

	t2 := getTable2()
	t3 := NewTableJoin().Table(t1).LeftJoin(t2, []string{"a", "c"}, nil, nil).Limit(1, 11).
		SetField(t1, []string{"a", "b"}, t2, []string{"c", "d"}).
		SetExpression(t1, []string{"a", "b"}, t2, []string{"`c` * 3", "`d` + 4 "}).
		SetValue(t1, []string{"a", "b"}, t2, []any{"abc", 1})
	fmt.Println(t3.UpdateSQL())
}

func TestDelete(t *testing.T) {
	t1 := getTable1()

	fmt.Println(t1.DeleteSQL())

	t2 := getTableJoin()

	fmt.Println(t2.DeleteSQL())
}

func getTable1() *Table {
	return NewTable("abc").Select([]string{"a", "b"}).
		SelectAlias(map[string]string{"a1": "a11", "b1": "b11"}).
		Aggregation(map[string]string{"COUNT(`a`)": "总数"}).
		Where("`a` = 1").
		GroupBy([]string{"c", "d"}).
		GroupByAlias("c1", "c11", "d1", "d11").
		OrderBy("a", command.SortDesc, "b", command.SortAsc).
		OrderByAsc("c").
		OrderByDesc("d").
		Limit(0, 11)
}

func getTable2() *Table {
	return NewTable("efg").Select([]string{"a", "b"}).
		SelectAlias(map[string]string{"a1": "a11", "b1": "b11", "c": "cc"}).
		Aggregation(map[string]string{"COUNT(`a`)": "总数"}).
		Where("`a` = 1").
		GroupBy([]string{"c", "d"}).
		GroupByAlias("c1", "c11", "d1", "d11").
		OrderBy("b", command.SortAsc, "a", command.SortDesc).
		OrderByAsc("c").
		OrderByDesc("d").
		Limit(0, 11)
}

func getTableJoin() *TableJoin {
	t1 := getTable1()

	t2 := getTable2()

	return NewTableJoin().Table(t1).LeftJoin(t2, []string{"a", "c"}, nil, nil).Limit(1, 11)
}
