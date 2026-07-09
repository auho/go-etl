package source

import (
	"fmt"
	"testing"
)

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
