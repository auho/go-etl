package mode

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/v2/means/segword"
	"github.com/auho/go-etl/v2/means/tag"
)

func Test_InsertMode(t *testing.T) {
	ttm := tag.NewKey(_rule)
	ti1 := NewInsert([]string{_keyName}, ttm)
	results := ti1.Do(_item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	tmtm := tag.NewMostText(_rule)
	ti2 := NewInsert([]string{_keyName}, tmtm)
	results = ti2.Do(_item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	tmkm := tag.NewMostKey(_rule)
	ti3 := NewInsert([]string{_keyName}, tmkm)
	results = ti3.Do(_item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	sw := segword.NewSegWordsMeans()
	ti4 := NewInsert([]string{_keyName}, sw)
	results = ti4.Do(_item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	ti5 := NewInsertMulti([]string{_keyName}, sw, sw)
	results2 := ti5.Do(_item)
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
