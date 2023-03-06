package mode

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/means/tag"
)

func Test_UpdateMode(t *testing.T) {
	tmtm := tag.NewMostText(ruleName, tag.WithDBRule(db))
	ti2 := NewUpdate([]string{keyName}, tmtm)
	results := ti2.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	tmkm := tag.NewMostKey(ruleName, tag.WithDBRule(db))
	ti3 := NewUpdate([]string{keyName}, tmkm)
	results = ti3.Do(item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)
}
