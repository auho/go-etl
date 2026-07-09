package source

import (
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
