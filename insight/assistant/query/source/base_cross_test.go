package source

import (
	"fmt"
	"testing"
)

func Test_baseCross_expandItems(t *testing.T) {
	bc := &baseCross{}
	items := map[string][]any{
		"one": {"1", "2", "3"},
		"tow": {"1", "2", "3"},
	}

	_items := bc.expandItemsCross(items)
	for _, item := range _items {
		fmt.Println(item)
	}
}
