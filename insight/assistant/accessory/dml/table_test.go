package dml

import (
	"fmt"
	"os"
	"testing"

	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml/command"
)

func TestMain(t *testing.M) {
	code := t.Run()

	os.Exit(code)
}

func TestTable(t *testing.T) {
	t1 := getTable1()
	s1 := t1.Sql()
	fmt.Println(s1)

	t2 := NewSqlTable("efg", s1).Select([]string{"b11", "a11"})
	s2 := t2.Sql()
	fmt.Println(s2)

	fmt.Println(t1.GetSelectFields())
	fmt.Println(t2.GetSelectFields())
}

func TestTableJoin(t *testing.T) {
	t1 := getTable1()

	t2 := getTable2()

	t3 := NewTableJoin().Table(t1).LeftJoin(t2, []string{"a", "c"}, nil, nil).Limit(1, 11)
	s3 := t3.Sql()
	fmt.Println(s3)
	fmt.Println(t3.GetSelectFields())

	t4 := t1.CreateJoin().LeftJoin(t2, []string{"a", "c"}, nil, nil).Limit(2, 22)
	s4 := t4.Sql()
	fmt.Println(s4)
	fmt.Println(t4.GetSelectFields())
}

func TestInsert(t *testing.T) {
	t1 := getTable1()

	fmt.Println(t1.InsertSql("i1"))

	fmt.Println(t1.InsertWithFieldsSql("i2", []string{"a", "a11", "d11"}))

	fmt.Println(t1.InsertWithFieldsSql("i2", nil))

	t2 := getTableJoin()

	fmt.Println(t2.InsertSql("i1"))

	fmt.Println(t2.InsertWithFieldsSql("i2", []string{"a", "a11", "d11"}))

	fmt.Println(t2.InsertWithFieldsSql("i2", nil))
}

func TestUpdate(t *testing.T) {
	t1 := getTable1().
		SetField(map[string]string{"a": "b", "c": "d"}).
		SetExpression(map[string]string{"a": "`b` + 1 ", "c": "`d` * 2"}).
		SetValue(map[string]any{"e": "abc", "f": 1, "g": 1.11})
	fmt.Println(t1.UpdateSql())

	t2 := getTable2()
	t3 := NewTableJoin().Table(t1).LeftJoin(t2, []string{"a", "c"}, nil, nil).Limit(1, 11).
		SetField(t1, []string{"a", "b"}, t2, []string{"c", "d"}).
		SetExpression(t1, []string{"a", "b"}, t2, []string{"`c` * 3", "`d` + 4 "}).
		SetValue(t1, []string{"a", "b"}, t2, []any{"abc", 1})
	fmt.Println(t3.UpdateSql())
}

func TestDelete(t *testing.T) {
	t1 := getTable1()

	fmt.Println(t1.DeleteSql())

	t2 := getTableJoin()

	fmt.Println(t2.DeleteSql())
}

func getTable1() *Table {
	return NewTable("abc").Select([]string{"a", "b"}).
		SelectAlias(map[string]string{"a1": "a11", "b1": "b11"}).
		Aggregation(map[string]string{"COUNT(`a`)": "总数"}).
		Where("`a` = 1").
		GroupBy([]string{"c", "d"}).
		GroupByAlias(map[string]string{"c1": "c11", "d1": "d11"}).
		OrderBy("a", command.SortDesc, "b", command.SortASC).
		Limit(0, 11)
}

func getTable2() *Table {
	return NewTable("efg").Select([]string{"a", "b"}).
		SelectAlias(map[string]string{"a1": "a11", "b1": "b11"}).
		Aggregation(map[string]string{"COUNT(`a`)": "总数"}).
		Where("`a` = 1").
		GroupBy([]string{"c", "d"}).
		GroupByAlias(map[string]string{"c1": "c11", "d1": "d11"}).
		OrderBy("b", command.SortASC, "a", command.SortDesc).
		Limit(0, 11)
}

func getTableJoin() *TableJoin {
	t1 := getTable1()

	t2 := getTable2()

	return NewTableJoin().Table(t1).LeftJoin(t2, []string{"a", "c"}, nil, nil).Limit(1, 11)
}
