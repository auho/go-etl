package transform

import (
	"testing"

	"github.com/auho/go-etl/v3/job/extract/segword"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

func TestInsertCross(t *testing.T) {
	ins1 := newTestInsert(t)
	ret1, err := ins1.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}

	ins2 := NewInsert(collector.NewKeysAll([]string{_keyName}, segword.NewSegWordsAll()), nil)
	err = ins2.Prepare()
	if err != nil {
		t.Fatal(err)
	}
	ret2, err := ins2.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}

	cross := NewInsertCross(ins1, ins2)
	err = cross.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	retCross, err := cross.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}

	expected := len(ret1) * len(ret2)
	if len(retCross) != expected {
		t.Fatalf("expected %d results from cross, got %d", expected, len(retCross))
	}
}

func TestInsertCross_SingleInsert(t *testing.T) {
	ins := newTestInsert(t)

	cross := NewInsertCross(ins)
	err := cross.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	ret, err := cross.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) <= 0 {
		t.Error("expected results from single-insert cross")
	}
}
