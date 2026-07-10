package source

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/v3/insight/assistant/sqlbuilder/dml"
)

// --- unit tests ---

func TestNewPlaceholder(t *testing.T) {
	p := NewPlaceholder(Base{}).
		AppendItems([]map[string]any{
			{"one": "1", "two": "2"},
			{"one": "1", "two": "21"},
			{"one": "11", "two": "2"},
			{"one": "11", "two": "21"},
		},
		)

	fmt.Println(p.items)
	expect := 4
	actual := len(p.items)
	if expect != actual {
		t.Fatalf("expect[%d] != actual[%d]", expect, actual)
	}

	p.SetItemsCross(map[string][]any{
		"one":   {"1", "11"},
		"two":   {"2", "21"},
		"three": {"3", "31"},
	})

	fmt.Println(p.items)
	expect = 8
	actual = len(p.items)
	if expect != actual {
		t.Fatalf("expect[%d] != actual[%d]", expect, actual)
	}
}

func TestPlaceholderAppendAccumulate(t *testing.T) {
	p := NewPlaceholder(Base{}).
		AppendItems([]map[string]any{
			{"one": "1"},
		})

	p.AppendItems([]map[string]any{
		{"one": "2"},
	})

	if len(p.items) != 2 {
		t.Fatalf("expect[2] != actual[%d]", len(p.items))
	}
}

func TestPlaceholderSetReplace(t *testing.T) {
	p := NewPlaceholder(Base{}).
		AppendItems([]map[string]any{
			{"one": "1"},
			{"one": "2"},
		})

	p.SetItems([]map[string]any{
		{"one": "3"},
	})

	if len(p.items) != 1 {
		t.Fatalf("expect[1] != actual[%d]", len(p.items))
	}
}

func TestPlaceholderDatasetEmptyItems(t *testing.T) {
	p := NewPlaceholder(Base{})
	_, err := p.Dataset()
	if err == nil {
		t.Fatal("expect error for empty items")
	}
}

// --- integration tests ---

func TestPlaceholderDataset(t *testing.T) {
	skipIfNoDB(t)

	s := NewPlaceholder(Base{
		Name: "placeholder",
		Table: dml.NewTable(_testTable).
			Select([]string{"name", "value"}).
			Where("category = '##category##'"),
		DB: _simpleDB,
	})

	s.AppendItems([]map[string]any{
		{"category": "cat1"},
		{"category": "cat2"},
	})

	ds, err := s.Dataset()
	if err != nil {
		t.Fatal(err)
	}

	if len(ds.Keys) != 1 || ds.Keys[0] != "category" {
		t.Fatalf("expect[category] != actual[%v]", ds.Keys)
	}
	if len(ds.Sets) != 2 {
		t.Fatalf("expect[2] != actual[%d]", len(ds.Sets))
	}
	// cat1 has 2 rows, cat2 has 2 rows
	for _, set := range ds.Sets {
		if set.Amount != 2 {
			t.Fatalf("expect[2] != actual[%d] for set[%s]", set.Amount, set.Name)
		}
	}
}

func TestPlaceholderMultiPlaceholder(t *testing.T) {
	skipIfNoDB(t)

	s := NewPlaceholder(Base{
		Name: "multi",
		Table: dml.NewTable(_testTable).
			Select([]string{"name"}).
			Where("category = '##category##' AND value = ##value##"),
		DB: _simpleDB,
	})

	s.AppendItems([]map[string]any{
		{"category": "cat1", "value": 10},
		{"category": "cat1", "value": 20},
		{"category": "cat2", "value": 30},
		{"category": "cat2", "value": 40},
		// duplicate of first item
		{"category": "cat1", "value": 10},
	})

	ds, err := s.Dataset()
	if err != nil {
		t.Fatal(err)
	}

	// 4 unique items (1 duplicate removed)
	if len(ds.Sets) != 4 {
		t.Fatalf("expect[4] != actual[%d]", len(ds.Sets))
	}
	// each query returns exactly 1 row
	for _, set := range ds.Sets {
		if set.Amount != 1 {
			t.Fatalf("expect[1] != actual[%d] for set[%s]", set.Amount, set.Name)
		}
	}
}

func TestPlaceholderHasNamePrefix(t *testing.T) {
	skipIfNoDB(t)

	s := NewPlaceholder(Base{
		Name:          "prefix",
		HasNamePrefix: true,
		Table: dml.NewTable(_testTable).
			Select([]string{"name"}).
			Where("category = '##category##'"),
		DB: _simpleDB,
	})

	s.AppendItems([]map[string]any{
		{"category": "cat1"},
	})

	ds, err := s.Dataset()
	if err != nil {
		t.Fatal(err)
	}

	if len(ds.Sets) != 1 {
		t.Fatalf("expect[1] != actual[%d]", len(ds.Sets))
	}
	// ID should have name prefix: "prefix_cat1"
	if ds.Sets[0].Name != "prefix_cat1" {
		t.Fatalf("expect[prefix_cat1] != actual[%s]", ds.Sets[0].Name)
	}
}

func TestPlaceholderEmptyResult(t *testing.T) {
	skipIfNoDB(t)

	s := NewPlaceholder(Base{
		Name: "empty",
		Table: dml.NewTable(_testTable).
			Select([]string{"name"}).
			Where("category = '##category##' AND value = ##value##"),
		DB: _simpleDB,
	})

	s.AppendItems([]map[string]any{
		{"category": "cat1", "value": 999}, // no matching row
	})

	ds, err := s.Dataset()
	if err != nil {
		t.Fatal(err)
	}

	if len(ds.Sets) != 1 {
		t.Fatalf("expect[1] != actual[%d]", len(ds.Sets))
	}
	if ds.Sets[0].Amount != 0 {
		t.Fatalf("expect[0] != actual[%d]", ds.Sets[0].Amount)
	}
}
