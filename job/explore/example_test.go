package explore

import (
	"github.com/auho/go-etl/v2/job/explore/collect"
	"github.com/auho/go-etl/v2/job/explore/condition"
	"github.com/auho/go-etl/v2/job/means/match"
)

func ExampleGenExplore() {
	insert := GenExplore().
		SetCollect(collect.NewKeys([]string{_keyName})).
		SetSearch(match.NewKey(_rule)).
		SetCondition(condition.NewContainAll("a", []string{"a1", "a2"})).
		ToInsert()

	update := GenExplore().
		SetCollect(collect.NewKeys([]string{_keyName})).
		SetSearch(match.NewKey(_rule)).
		SetCondition(condition.NewContainAll("a", []string{"a1", "a2"})).
		ToUpdate()

	_ = insert
	_ = update
}

func ExampleNewInsert() {
	_ = NewInsert(collect.NewKeys([]string{_keyName}), match.NewKey(_rule), nil)
}

func ExampleNewUpdate() {
	_ = NewInsert(collect.NewKeys([]string{_keyName}), match.NewKey(_rule), nil)
}

func ExampleNewInsertCross() {
	i1 := NewInsert(collect.NewKeys([]string{_keyName}), match.NewKey(_rule), nil)
	i2 := NewInsert(collect.NewKeys([]string{_keyName}), match.NewKey(_rule), nil)

	_ = NewInsertCross(i1, i2)
}
