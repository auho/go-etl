package source

import (
	"fmt"
	"testing"
)

func TestNewPlaceholder(t *testing.T) {
	p := NewPlaceholder(Source{}).
		WithItems([]map[string]any{
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

	p.WithItemsCross(map[string][]any{
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
