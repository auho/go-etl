package mode

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/means/tager"
)

func Test_UpdateMode(t *testing.T) {
	item := make(map[string]interface{})
	item[keyName] = content

	tmtm := tager.NewTagMostTextMeans(ruleName, db)
	ti2 := NewTagUpdate([]string{keyName}, tmtm)
	results := ti2.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	tmkm := tager.NewTagMostKeyMeans(ruleName, db)
	ti3 := NewTagUpdate([]string{keyName}, tmkm)
	results = ti3.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)
}
