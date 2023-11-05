package mode

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/v2/job/means/segword"
	"github.com/auho/go-etl/v2/job/means/tag"
)

func Test_InsertMode(t *testing.T) {
	// key
	tKey := tag.NewKey(_rule)
	tin1 := NewInsert([]string{_keyName}, tKey)
	err := tin1.Prepare()
	if err != nil {
		t.Fatal("tin1", err)
	}

	results := tin1.Do(_item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	// most text
	tMostText := tag.NewMostText(_rule)
	tin2 := NewInsert([]string{_keyName}, tMostText)
	err = tin2.Prepare()
	if err != nil {
		t.Fatal("tin2", err)
	}

	results = tin2.Do(_item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	// most key
	tMostKey := tag.NewMostKey(_rule)
	tin3 := NewInsert([]string{_keyName}, tMostKey)
	err = tin3.Prepare()
	if err != nil {
		t.Fatal("tin3", err)
	}

	results = tin3.Do(_item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	// seg words
	tSegWords := segword.NewSegWordsMeans()
	tin4 := NewInsert([]string{_keyName}, tSegWords)
	err = tin4.Prepare()
	if err != nil {
		t.Fatal("tin4", err)
	}

	results = tin4.Do(_item)
	if len(results) <= 0 {
		t.Error("error")
	}
	fmt.Println(results)

	ti5 := NewInsertMulti([]string{_keyName}, tSegWords, tSegWords)
	err = ti5.Prepare()
	if err != nil {
		t.Fatal("ti5", err)
	}

	results2 := ti5.Do(_item)
	if len(results2) <= 0 {
		t.Error("error")
	}

	for _, v := range ti5.GetKeys() {
		has := false
		for _, v1 := range tSegWords.GetKeys() {
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
