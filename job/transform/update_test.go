package transform

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/v3/job/extract/tag"
)

func Test_UpdateMode(t *testing.T) {
	tMostText := tag.NewMostText(_rule).ToMeans()
	tin1 := NewUpdate([]string{_keyName}, tMostText)
	err := tin1.Prepare()
	if err != nil {
		t.Fatal("tin1", err)
	}

	results := tin1.Apply(_item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	tMostKey := tag.NewMostKey(_rule).ToMeans()
	tin2 := NewUpdate([]string{_keyName}, tMostKey)
	err = tin2.Prepare()
	if err != nil {
		t.Fatal("tin1", err)
	}

	results = tin2.Apply(_item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)
}
