package transform

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/v3/job/transform/collect"
	"github.com/auho/go-etl/v3/job/extract/segword"
	"github.com/auho/go-etl/v3/job/extract/tag"
)

func Test_InsertMode(t *testing.T) {
	// key
	_modeKeys := NewInsert(collect.NewKeys([]string{_keyName}), tag.NewKey(_rule), nil)
	err := _modeKeys.Prepare()
	if err != nil {
		t.Fatal("_modeKeys", err)
	}

	retKeys := _modeKeys.Apply(_item)
	if len(retKeys) <= 0 {
		t.Error("error")
	}

	fmt.Println(retKeys)

	// most text
	_modeMostText := NewInsert(collect.NewKeys([]string{_keyName}), tag.NewMostText(_rule), nil)
	err = _modeMostText.Prepare()
	if err != nil {
		t.Fatal("_modeMostText", err)
	}

	retMostText := _modeMostText.Apply(_item)
	if len(retMostText) <= 0 {
		t.Error("error")
	}
	fmt.Println(retMostText)

	// most key
	_modeMostKey := NewInsert(collect.NewKeys([]string{_keyName}), tag.NewMostKey(_rule), nil)
	err = _modeMostKey.Prepare()
	if err != nil {
		t.Fatal("_modeMostKey", err)
	}

	retMostKey := _modeMostKey.Apply(_item)
	if len(retMostKey) <= 0 {
		t.Error("error")
	}
	fmt.Println(retMostKey)

	// seg words
	_modeSegWords := NewInsert(collect.NewKeys([]string{_keyName}), segword.NewDefault(), nil)
	err = _modeSegWords.Prepare()
	if err != nil {
		t.Fatal("_modeSegWords", err)
	}

	retSegWords := _modeSegWords.Apply(_item)
	if len(retSegWords) <= 0 {
		t.Error("error")
	}
	fmt.Println(retSegWords)

	// insert stack
	_modeInsertStack := NewInsertStack(_modeKeys, _modeSegWords)
	err = _modeInsertStack.Prepare()
	if err != nil {
		t.Fatal("_modeInsertStack", err)
	}

	retStack := _modeInsertStack.Apply(_item)
	if len(retStack) <= 0 {
		t.Error("error")
	}
	fmt.Println(retStack)

	for _, v := range _modeInsertStack.Keys() {
		has := false

		for _, v1 := range _modeKeys.Keys() {
			if v == v1 {
				has = true
				break
			}
		}

		for _, v1 := range _modeSegWords.Keys() {
			if v == v1 {
				has = true
				break
			}
		}

		if has == false {
			t.Error(fmt.Sprintf("key[%s] is error", v))
		}
	}

	if len(retKeys)+len(retSegWords) != len(retStack) {
		t.Error("error")
	}

	// insert cross
	_modeInsertCross := NewInsertCross(_modeMostText, _modeSegWords)
	err = _modeInsertCross.Prepare()
	if err != nil {
		t.Fatal("_modeInsertCross", err)
	}

	retCross := _modeInsertCross.Apply(_item)
	if len(retCross) <= 0 {
		t.Error("error")
	}
	fmt.Println(retCross)

	for _, v := range _modeInsertStack.Keys() {
		has := false

		for _, v1 := range _modeMostText.Keys() {
			if v == v1 {
				has = true
				break
			}
		}

		for _, v1 := range _modeSegWords.Keys() {
			if v == v1 {
				has = true
				break
			}
		}

		if has == false {
			t.Error(fmt.Sprintf("key[%s] is error", v))
		}
	}

	if len(retMostText)*len(retSegWords) != len(retCross) {
		t.Error("error")
	}

	// insert spread
	_modeKeys = NewInsert(collect.NewKeys([]string{_keyName}), tag.NewKey(_rule).WithPluck([]string{_rule.NameAlias()}), nil)
	_modeInsertSpread := NewInsertSpread(_modeKeys, _modeSegWords)
	err = _modeInsertSpread.Prepare()
	if err != nil {
		t.Error("error")
	}

	retSpread := _modeInsertSpread.Apply(_item)
	if len(retSpread) <= 0 {
		t.Error("insert spread error")
	}
	fmt.Println(retSpread)

	for _, v := range _modeInsertSpread.Keys() {
		has := false

		for _, v1 := range _modeKeys.Keys() {
			if v == v1 {
				has = true
				break
			}
		}

		for _, v1 := range _modeSegWords.Keys() {
			if v == v1 {
				has = true
				break
			}
		}

		if has == false {
			t.Error(fmt.Sprintf("key[%s] is error", v))
		}
	}
}
