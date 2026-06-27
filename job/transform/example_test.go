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

func ExampleNewInsertStack() {
	i1 := NewInsert(collect.NewKeysAll([]string{_keyName}), match.NewKey(_rule), nil)
	i2 := NewInsert(collect.NewKeysAll([]string{_keyName}), match.NewKey(_rule), nil)

	_ = NewInsertStack(i1, i2)
}

func ExampleNewInsertSpread() {
	i1 := NewInsert(collect.NewKeysAll([]string{_keyName}), match.NewKey(_rule), nil)
	i2 := NewInsert(collect.NewKeysAll([]string{_keyName}), match.NewKey(_rule), nil)

	_ = NewInsertSpread(i1, i2)
}

func ExampleNewInsertComposeSpread() {
	i1 := NewInsert(collect.NewKeysAll([]string{_keyName}), match.NewKey(_rule), nil)
	i2 := NewInsert(collect.NewKeysAll([]string{_keyName}), match.NewKey(_rule), nil)

	_ = NewInsertComposeSpread(i1, i2)
}

func ExampleNewTransfer() {
	_ = NewTransfer(
		[]string{"a", "b"},
		map[string]string{"a": "x"},
		map[string]any{"c": "d"},
	)
}