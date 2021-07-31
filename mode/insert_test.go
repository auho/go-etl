package mode

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/means"
	"github.com/auho/go-etl/means/segworder"
	"github.com/auho/go-etl/means/tager"
)

func Test_InsertMode(t *testing.T) {
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

	ti5 := NewInsertMultiMode([]string{keyName}, sw.GetKeys(), []means.InsertMeans{sw, sw})
	results2 := ti5.Do(item)
	if len(results2) <= 0 {
		t.Error("error")
	}

	for _, v := range ti5.GetKeys() {
		has := false
		for _, v1 := range sw.GetKeys() {
			if v == v1 {
				has = true
				break
			}
		}

		if has == false {
			t.Error(fmt.Sprintf("key[%s] is error", v))
		}
	}

	if len(results)*2 != len(results2) {
		t.Error("multi insert mode is error")
	}
}
