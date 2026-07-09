package source

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/v3/insight/assistant/sqlbuilder/dml"
)

// --- unit tests ---

func TestNewPlaceholderStack(t *testing.T) {
	ps := NewPlaceholderStack(Base{}).
		AppendCategories([]map[string]any{
			{"three": "1"},
			{"three": "2"},
		}).AppendStacks([]map[string]any{
		{"one": "1", "two": "2"},
		{"one": "1", "two": "21"},
		{"one": "11", "two": "2"},
		{"one": "11", "two": "21"},
	})

	fmt.Println(ps.categories)
	fmt.Println(ps.stacks)
	expect := 2
	actual := len(ps.categories)
	if expect != actual {
		t.Fatalf("catgory expect[%d] != actual[%d]", expect, actual)
	}

	expect = 4
	actual = len(ps.stacks)
	if expect != actual {
		t.Fatalf("stack expect[%d] != actual[%d]", expect, actual)
	}

	ps.SetCategoriesCross(map[string][]any{
		"three": {1, 2, 3}, "four": {4, 5, 6},
	}).SetStacksCross(map[string][]any{
		"one": {"1", "11"},
		"two": {"2", "21"},
	},
	)

	fmt.Println(ps.categories)
	fmt.Println(ps.stacks)
	expect = 9
	actual = len(ps.categories)
	if expect != actual {
		t.Fatalf("catgory1 expect[%d] != actual[%d]", expect, actual)
	}

	expect = 4
	actual = len(ps.stacks)
	if expect != actual {
		t.Fatalf("stack1 expect[%d] != actual[%d]", expect, actual)
	}
}

func TestPlaceholderStackAppendAccumulate(t *testing.T) {
	ps := NewPlaceholderStack(Base{}).
		AppendCategories([]map[string]any{
			{"three": "1"},
		})

	ps.AppendCategories([]map[string]any{
		{"three": "2"},
	})

	if len(ps.categories) != 2 {
		t.Fatalf("expect[2] != actual[%d]", len(ps.categories))
	}
}

func TestPlaceholderStackSetReplace(t *testing.T) {
	ps := NewPlaceholderStack(Base{}).
		AppendStacks([]map[string]any{
			{"one": "1"},
			{"one": "2"},
		})

	ps.SetStacks([]map[string]any{
		{"one": "3"},
	})

	if len(ps.stacks) != 1 {
		t.Fatalf("expect[1] != actual[%d]", len(ps.stacks))
	}
}

func TestPlaceholderStackDatasetEmptyCategories(t *testing.T) {
	pss := NewPlaceholderStack(Base{})
	_, err := pss.Dataset()
	if err == nil {
		t.Fatal("expect error for empty categories")
	}
}

func TestCategoryToID(t *testing.T) {
	pss := &PlaceholderStackSource{Base: Base{Name: "test"}}
	keys := []string{"one", "two"}
	category := map[string]any{"one": "a", "two": "b"}
	id := pss.categoryToID(category, keys)
	// keys order: one=a, two=b
	if id != "a_b" {
		t.Fatalf("expect[a_b] != actual[%s]", id)
	}
}

func TestCategoryToIDCrossValues(t *testing.T) {
	pss := &PlaceholderStackSource{Base: Base{Name: "test"}}
	keys := []string{"one", "two"}
	category := map[string]any{"one": "b", "two": "a"}
	id := pss.categoryToID(category, keys)
	// keys order: one=b, two=a (not sorted)
	if id != "b_a" {
		t.Fatalf("expect[b_a] != actual[%s]", id)
	}
}

func TestCategoryToIDPartialKeys(t *testing.T) {
	pss := &PlaceholderStackSource{Base: Base{Name: "test"}}
	keys := []string{"one", "two", "three"}
	category := map[string]any{"three": 1}
	id := pss.categoryToID(category, keys)
	// only "three" is in category
	if id != "1" {
		t.Fatalf("expect[1] != actual[%s]", id)
	}
}

// --- integration tests ---

func TestPlaceholderStackSourceDataset(t *testing.T) {
	skipIfNoDB(t)

	s := NewPlaceholderStack(Base{
		Name: "stack",
		Table: dml.NewTable(_testTable).
			Select([]string{"name"}).
			Where("category = '##category##' AND value = ##value##"),
		DB: _simpleDB,
	})

	s.AppendCategories([]map[string]any{
		{"category": "cat1"},
		{"category": "cat2"},
	})

	s.AppendStacks([]map[string]any{
		{"value": 10},
		{"value": 20},
		{"value": 30},
		{"value": 40},
	})

	ds, err := s.Dataset()
	if err != nil {
		t.Fatal(err)
	}

	// 2 categories
	if len(ds.Sets) != 2 {
		t.Fatalf("expect[2] != actual[%d]", len(ds.Sets))
	}
}

func TestPlaceholderStackSourceCategoryDedup(t *testing.T) {
	skipIfNoDB(t)

	s := NewPlaceholderStack(Base{
		Name: "stack_dedup",
		Table: dml.NewTable(_testTable).
			Select([]string{"name"}).
			Where("category = '##category##' AND value = ##value##"),
		DB: _simpleDB,
	})

	s.AppendCategories([]map[string]any{
		{"category": "cat1"},
		{"category": "cat1"}, // duplicate
		{"category": "cat2"},
	})

	s.AppendStacks([]map[string]any{
		{"value": 10},
	})

	ds, err := s.Dataset()
	if err != nil {
		t.Fatal(err)
	}

	// duplicate category removed, only 2 sets
	if len(ds.Sets) != 2 {
		t.Fatalf("expect[2] != actual[%d]", len(ds.Sets))
	}
}
