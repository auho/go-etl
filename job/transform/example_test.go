package transform

import (
	"github.com/auho/go-etl/v3/job/extract/match"
	"github.com/auho/go-etl/v3/job/transform/collect"
)

func ExampleNewInsert() {
	_ = NewInsert(collect.NewKeysAll([]string{_keyName}), match.NewKey(_rule), nil)
}

func ExampleNewUpdate() {
	_ = NewUpdate(collect.NewKeysAll([]string{_keyName}), match.NewKey(_rule), nil)
}

func ExampleNewInsertCross() {
	i1 := NewInsert(collect.NewKeysAll([]string{_keyName}), match.NewKey(_rule), nil)
	i2 := NewInsert(collect.NewKeysAll([]string{_keyName}), match.NewKey(_rule), nil)

	_ = NewInsertCross(i1, i2)
}
