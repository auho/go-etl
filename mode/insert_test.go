package mode

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/means"
	"github.com/auho/go-etl/means/segworder"
	"github.com/auho/go-etl/means/tager"
)

func Test_InsertMode(t *testing.T) {
	item := make(map[string]interface{})
	item[keyName] = content

	ttm := tager.NewTagKeyMeans(ruleName, db)
	ti1 := NewInsertMode([]string{keyName}, ttm)
	results := ti1.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	tmtm := tager.NewTagMostTextMeans(ruleName, db)
	ti2 := NewInsertMode([]string{keyName}, tmtm)
	results = ti2.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	tmkm := tager.NewTagMostKeyMeans(ruleName, db)
	ti3 := NewInsertMode([]string{keyName}, tmkm)
	results = ti3.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	sw := segworder.NewSegWordsMeans()
	ti4 := NewInsertMode([]string{keyName}, sw)
	results = ti4.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	ti5 := NewMultiInsertMode([]string{keyName}, sw.GetKeys(), []means.InsertMeans{sw, sw})
	results = ti5.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)
}
