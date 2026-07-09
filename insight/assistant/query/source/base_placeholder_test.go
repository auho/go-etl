package source

import (
	"testing"
)

func TestBuildKeys(t *testing.T) {
	bph := &basePlaceholder{}
	sql := "SELECT * FROM t WHERE a = '##one##' AND b = ##two##"
	keys := bph.buildKeys(sql)
	if len(keys) != 2 {
		t.Fatalf("expect[2] != actual[%d]", len(keys))
	}
	if keys[0] != "one" || keys[1] != "two" {
		t.Fatalf("expect[one two] != actual[%v]", keys)
	}
}

func TestBuildKeysNoPlaceholder(t *testing.T) {
	bph := &basePlaceholder{}
	keys := bph.buildKeys("SELECT * FROM t")
	if len(keys) != 0 {
		t.Fatalf("expect[0] != actual[%d]", len(keys))
	}
}

func TestBuildPlaceholderItemSql(t *testing.T) {
	bph := &basePlaceholder{}
	sql := "SELECT * FROM t WHERE a = '##one##' AND b = ##two##"
	item := map[string]any{"one": "val1", "two": 2}
	result := bph.buildPlaceholderItemSql(sql, item)
	expect := "SELECT * FROM t WHERE a = 'val1' AND b = 2"
	if result != expect {
		t.Fatalf("expect[%s] != actual[%s]", expect, result)
	}
}

func TestBuildPlaceholderItemSqlNoPlaceholder(t *testing.T) {
	bph := &basePlaceholder{}
	sql := "SELECT * FROM t"
	item := map[string]any{"one": "val1"}
	result := bph.buildPlaceholderItemSql(sql, item)
	if result != sql {
		t.Fatalf("expect[%s] != actual[%s]", sql, result)
	}
}

func TestBuildPlaceholderItemsSqlSet(t *testing.T) {
	bph := &basePlaceholder{}
	s := Base{}
	sql := "SELECT * FROM t WHERE a = '##one##'"
	keys := []string{"one"}
	items := []map[string]any{
		{"one": "a"},
		{"one": "a"}, // duplicate
		{"one": "b"},
	}
	itemIds, itemsSql := bph.buildPlaceholderItemsSqlSet(s, sql, keys, items)
	if len(itemIds) != 2 {
		t.Fatalf("expect[2] != actual[%d]", len(itemIds))
	}
	if len(itemsSql) != 2 {
		t.Fatalf("expect[2] != actual[%d]", len(itemsSql))
	}
}

func TestBuildPlaceholderItemsSqlList(t *testing.T) {
	bph := &basePlaceholder{}
	sql := "SELECT * FROM t WHERE a = '##one##'"
	items := []map[string]any{
		{"one": "a"},
		{"one": "b"},
	}
	list := bph.buildPlaceholderItemsSqlList(sql, items)
	if len(list) != 2 {
		t.Fatalf("expect[2] != actual[%d]", len(list))
	}
}

func TestBuildPlaceholderItemsSqlSetDedupIDs(t *testing.T) {
	bph := &basePlaceholder{}
	s := Base{}
	sql := "SELECT * FROM t WHERE a = '##one##'"
	keys := []string{"one"}
	items := []map[string]any{
		{"one": "a"},
		{"one": "a"}, // duplicate
		{"one": "b"},
		{"one": "b"}, // duplicate
		{"one": "c"},
	}
	itemIds, _ := bph.buildPlaceholderItemsSqlSet(s, sql, keys, items)
	if len(itemIds) != 3 {
		t.Fatalf("expect[3] != actual[%d]", len(itemIds))
	}

	// verify IDs
	idMap := make(map[string]bool)
	for _, id := range itemIds {
		idMap[id] = true
	}
	for _, expected := range []string{"a", "b", "c"} {
		if !idMap[expected] {
			t.Fatalf("missing itemId: %s", expected)
		}
	}
}
