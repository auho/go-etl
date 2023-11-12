package source

import (
	"fmt"
	"testing"
)

func Test_baseCross_expandItems(t *testing.T) {
	bc := &baseCross{}
	items := []map[string][]string{
		{"one": []string{"1", "2", "3"}},
		{"tow": []string{"1", "2", "3"}},
	}

	_items := bc.expandItemsCross(items)
	for _, item := range _items {
		fmt.Println(item)
	}
}
