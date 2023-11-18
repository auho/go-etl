package source

import (
	"fmt"
	"testing"
)

func TestNewPlaceholderStack(t *testing.T) {
	ps := NewPlaceholderStack(Source{}).
		WithCategories([]map[string]any{
			{"three": "1"},
			{"three": "2"},
		}).WithStacks([]map[string]any{
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

	ps.WithCategoriesCross(map[string][]any{
		"three": {1, 2, 3}, "four": {4, 5, 6},
	}).WithStacksCross(map[string][]any{
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
