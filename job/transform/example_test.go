package transform

import (
	"github.com/auho/go-etl/v3/job/extract/match"
	"github.com/auho/go-etl/v3/job/transform/collect"
	"github.com/auho/go-etl/v3/job/transform/filter"
)

func ExampleNewPipeline() {
	insert := NewPipeline().
		SetCollector(collect.NewKeysAll([]string{_keyName})).
		SetExtractor(match.NewKey(_rule)).
		SetPredicate(filter.NewContainsAll("a", []string{"a1", "a2"})).
		ToInsert()

	update := NewPipeline().
		SetCollector(collect.NewKeysAll([]string{_keyName})).
		SetExtractor(match.NewKey(_rule)).
		SetPredicate(filter.NewContainsAll("a", []string{"a1", "a2"})).
		ToUpdate()

	_ = insert
	_ = update
}

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
