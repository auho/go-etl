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
	t1 := getTable1()
	s1 := t1.Sql()
	fmt.Println(s1)

	t2 := NewSqlTable("efg", s1).Select([]string{"a11", "b11"}).Sql()

	fmt.Println(t2)
}

func TestTableJoin(t *testing.T) {
	t1 := getTable1()

	t2 := getTable2()

	s3 := NewTableJoin().Table(t1).LeftJoin(t2, []string{"a", "c"}, nil, nil).Limit(1, 11).Sql()
	fmt.Println(s3)
}

func TestInsert(t *testing.T) {
	t1 := getTable1()

	fmt.Println(t1.Insert("i1"))

	fmt.Println(t1.InsertWithFields("i2", []string{"a", "a11", "d11"}))

	fmt.Println(t1.InsertWithFields("i2", nil))

	t2 := getTableJoin()

	fmt.Println(t2.Insert("i1"))

	fmt.Println(t2.InsertWithFields("i2", []string{"a", "a11", "d11"}))

	fmt.Println(t2.InsertWithFields("i2", nil))
}

func TestDelete(t *testing.T) {
	t1 := getTable1()

	fmt.Println(t1.Delete())

	t2 := getTableJoin()

	fmt.Println(t2.Delete())
}

func getTable1() *Table {
	return NewTable("abc").Select([]string{"a", "b"}).
		SelectAlias(map[string]string{"a1": "a11", "b1": "b11"}).
		Aggregation(map[string]string{"COUNT(`a`)": "总数"}).
		Where("`a` = 1").
		GroupBy([]string{"c", "d"}).
		GroupByAlias(map[string]string{"c1": "c11", "d1": "d11"}).
		OrderBy(map[string]string{"a": command.SortDesc, "b": command.SortASC}).
		Limit(0, 11)
}

func getTable2() *Table {
	return NewTable("efg").Select([]string{"a", "b"}).
		SelectAlias(map[string]string{"a1": "a11", "b1": "b11"}).
		Aggregation(map[string]string{"COUNT(`a`)": "总数"}).
		Where("`a` = 1").
		GroupBy([]string{"c", "d"}).
		GroupByAlias(map[string]string{"c1": "c11", "d1": "d11"}).
		OrderBy(map[string]string{"a": command.SortDesc, "b": command.SortASC}).
		Limit(0, 11)
}

func getTableJoin() *TableJoin {
	t1 := getTable1()

	t2 := getTable2()

	return NewTableJoin().Table(t1).LeftJoin(t2, []string{"a", "c"}, nil, nil).Limit(1, 11)
}
