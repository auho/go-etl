package mode

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/v2/means/tag"
)

func Test_UpdateMode(t *testing.T) {
	tmtm := tag.NewMostText(_rule)
	ti2 := NewUpdate([]string{_keyName}, tmtm)
	results := ti2.Do(_item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	tmkm := tag.NewMostKey(_rule)
	ti3 := NewUpdate([]string{_keyName}, tmkm)
	results = ti3.Do(_item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)
}
