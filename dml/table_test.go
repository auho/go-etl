package dml

import (
	"fmt"
	"os"
	"testing"

	"github.com/auho/go-etl/dml/command"
)

func TestMain(t *testing.M) {
	RegisterDriver(DriverMysql)
	code := t.Run()

	os.Exit(code)
}

func TestTable(t *testing.T) {
	t1 := NewTable("abc")
	s1 := t1.Select([]string{"a", "b"}).
		SelectAlias(map[string]string{"a1": "a11", "b1": "b11"}).
		Aggregation(map[string]string{"COUNT(`a`)": "总数"}).
		Where("`a` = 1").
		GroupBy([]string{"c", "d"}).
		GroupByAlias(map[string]string{"c1": "c11", "d1": "d11"}).
		OrderBy(map[string]string{"a": command.SortDesc, "b": command.SortASC}).
		Limit(0, 11).
		Sql()

	fmt.Println(s1)

	t2 := NewSqlTable("efg", s1).Select([]string{"a11", "b11"}).Sql()

	fmt.Println(t2)
}

func TestTableJoin(t *testing.T) {
	t1 := NewTable("abc").Select([]string{"a", "b"}).
		SelectAlias(map[string]string{"a1": "a11", "b1": "b11"}).
		Aggregation(map[string]string{"COUNT(`a`)": "总数"}).
		Where("`a` = 1").
		GroupBy([]string{"c", "d"}).
		GroupByAlias(map[string]string{"c1": "c11", "d1": "d11"}).
		OrderBy(map[string]string{"a": command.SortDesc, "b": command.SortASC}).
		Limit(0, 11)

	t2 := NewTable("efg").Select([]string{"a", "b"}).
		SelectAlias(map[string]string{"a1": "a11", "b1": "b11"}).
		Aggregation(map[string]string{"COUNT(`a`)": "总数"}).
		Where("`a` = 1").
		GroupBy([]string{"c", "d"}).
		GroupByAlias(map[string]string{"c1": "c11", "d1": "d11"}).
		OrderBy(map[string]string{"a": command.SortDesc, "b": command.SortASC}).
		Limit(0, 11)

	s3 := NewTableJoin().Table(t1).LeftJoin(t2, []string{"a", "c"}, nil, nil).Limit(1, 11).Sql()
	fmt.Println(s3)
}

func TestInsert(t *testing.T) {
	t1 := NewTable("abc").Select([]string{"a", "b"}).
		SelectAlias(map[string]string{"a1": "a11", "b1": "b11"}).
		Aggregation(map[string]string{"COUNT(`a`)": "总数"}).
		Where("`a` = 1").
		GroupBy([]string{"c", "d"}).
		GroupByAlias(map[string]string{"c1": "c11", "d1": "d11"}).
		OrderBy(map[string]string{"a": command.SortDesc, "b": command.SortASC}).
		Limit(0, 11)

	i1 := NewInsert("insert_table", t1, nil)
	s := i1.Sql()

	fmt.Println(s)
}
