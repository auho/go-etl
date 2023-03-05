package mode

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/means/segword"
	"github.com/auho/go-etl/means/tag"
)

func Test_InsertMode(t *testing.T) {
	ttm := tag.NewKey(ruleName, tag.WithDBRule(db))
	ti1 := NewInsert([]string{keyName}, ttm)
	results := ti1.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	tmtm := tag.NewMostText(ruleName, tag.WithDBRule(db))
	ti2 := NewInsert([]string{keyName}, tmtm)
	results = ti2.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	tmkm := tag.NewMostKey(ruleName, tag.WithDBRule(db))
	ti3 := NewInsert([]string{keyName}, tmkm)
	results = ti3.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	sw := segword.NewSegWordsMeans()
	ti4 := NewInsert([]string{keyName}, sw)
	results = ti4.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	ti5 := NewInsertMulti([]string{keyName}, sw.GetKeys(), sw, sw)
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
		t.Error("multi means mode is error")
	}
}
