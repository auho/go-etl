package source

import (
	"fmt"
	"testing"
)

func TestBaseCrossExpandItems(t *testing.T) {
	bc := &baseCross{}
	items := map[string][]any{
		"one": {"1", "2", "3"},
		"tow": {"1", "2", "3"},
	}

	_items := bc.expandItemsCross(items)
	if len(_items) != 9 {
		t.Fatalf("expect[9] != actual[%d]", len(_items))
	}
}

func TestExpandItemsCrossEmpty(t *testing.T) {
	bc := &baseCross{}
	_items := bc.expandItemsCross(nil)
	if len(_items) != 0 {
		t.Fatalf("expect[0] != actual[%d]", len(_items))
	}
}

func TestExpandItemsCrossSingleKey(t *testing.T) {
	bc := &baseCross{}
	_items := bc.expandItemsCross(map[string][]any{
		"one": {"a", "b"},
	})
	if len(_items) != 2 {
		t.Fatalf("expect[2] != actual[%d]", len(_items))
	}
}

func TestExpandItemsCrossEmptyValues(t *testing.T) {
	bc := &baseCross{}
	_items := bc.expandItemsCross(map[string][]any{
		"one": {},
		"two": {"a", "b"},
	})
	if len(_items) != 0 {
		t.Fatalf("expect[0] != actual[%d]", len(_items))
	}
}

func TestExpandItemsCrossContent(t *testing.T) {
	bc := &baseCross{}
	items := map[string][]any{
		"one": {"a", "b"},
		"two": {"x", "y"},
	}

	result := bc.expandItemsCross(items)
	if len(result) != 4 {
		t.Fatalf("expect[4] != actual[%d]", len(result))
	}

	// verify all 4 combinations exist
	combos := make(map[string]bool)
	for _, item := range result {
		key := fmt.Sprintf("%v_%v", item["one"], item["two"])
		combos[key] = true
	}

	for _, expected := range []string{"a_x", "a_y", "b_x", "b_y"} {
		if !combos[expected] {
			t.Fatalf("missing combination: %s", expected)
		}
	}
}
