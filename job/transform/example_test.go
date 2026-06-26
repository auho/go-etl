package transform

import (
	"github.com/auho/go-etl/v3/job/transform/collect"
	"github.com/auho/go-etl/v3/job/transform/filter"
	"github.com/auho/go-etl/v3/job/extract/match"
)

func ExampleNewPipeline() {
	insert := NewPipeline().
		SetCollect(collect.NewKeys([]string{_keyName})).
		SetSearch(match.NewKey(_rule)).
		SetCondition(filter.NewContainsAll("a", []string{"a1", "a2"})).
		ToInsert()

	update := NewPipeline().
		SetCollect(collect.NewKeys([]string{_keyName})).
		SetSearch(match.NewKey(_rule)).
		SetCondition(filter.NewContainsAll("a", []string{"a1", "a2"})).
		ToUpdate()

	_ = insert
	_ = update
}

func ExampleNewInsert() {
	_ = NewInsert(collect.NewKeys([]string{_keyName}), match.NewKey(_rule), nil)
}

func ExampleNewUpdate() {
	_ = NewUpdate(collect.NewKeys([]string{_keyName}), match.NewKey(_rule), nil)
}

func ExampleNewInsertCross() {
	i1 := NewInsert(collect.NewKeys([]string{_keyName}), match.NewKey(_rule), nil)
	i2 := NewInsert(collect.NewKeys([]string{_keyName}), match.NewKey(_rule), nil)

	_ = NewInsertCross(i1, i2)
}
