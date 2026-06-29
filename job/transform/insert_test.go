package transform

import (
	"testing"

	"github.com/auho/go-etl/v3/job/extract/segword"
	"github.com/auho/go-etl/v3/job/extract/tag"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

func TestInsert_Key(t *testing.T) {
	ins := NewInsert(collector.NewKeysAll([]string{_keyName}, tag.NewKey(_rule)), nil)
	err := ins.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	ret, err := ins.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) <= 0 {
		t.Error("expected results from Key insert")
	}
}

func TestInsert_MostText(t *testing.T) {
	ins := NewInsert(collector.NewKeysAll([]string{_keyName}, tag.NewMostText(_rule)), nil)
	err := ins.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	ret, err := ins.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) <= 0 {
		t.Error("expected results from MostText insert")
	}
}

func TestInsert_MostKey(t *testing.T) {
	ins := NewInsert(collector.NewKeysAll([]string{_keyName}, tag.NewMostKey(_rule)), nil)
	err := ins.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	ret, err := ins.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) <= 0 {
		t.Error("expected results from MostKey insert")
	}
}

func TestInsert_SegWords(t *testing.T) {
	ins := NewInsert(collector.NewKeysAll([]string{_keyName}, segword.NewDefault()), nil)
	err := ins.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	ret, err := ins.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) <= 0 {
		t.Error("expected results from SegWords insert")
	}
}

func TestInsert_NilPredicate(t *testing.T) {
	ins := NewInsert(collector.NewKeysAll([]string{_keyName}, tag.NewKey(_rule)), nil)
	err := ins.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	ret, err := ins.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) <= 0 {
		t.Error("expected results with nil predicate")
	}
}