package transform

import (
	"testing"

	"github.com/auho/go-etl/v3/job/extract/segword"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

func TestInsertStack(t *testing.T) {
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

	stack := NewInsertStack(ins1, ins2)
	err = stack.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	retStack, err := stack.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}

	if len(retStack) != len(ret1)+len(ret2) {
		t.Fatalf("expected %d results from stack, got %d", len(ret1)+len(ret2), len(retStack))
	}
}

func TestInsertStack_SingleInsert(t *testing.T) {
	ins := newTestInsert(t)

	stack := NewInsertStack(ins)
	err := stack.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	ret, err := stack.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) <= 0 {
		t.Error("expected results from single-insert stack")
	}
}

func TestInsertStack_Keys(t *testing.T) {
	ins := newTestInsert(t)
	ins2 := NewInsert(collector.NewKeysAll([]string{_keyName}, segword.NewSegWordsAll()), nil)
	err := ins2.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	stack := NewInsertStack(ins, ins2)
	err = stack.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	_Keys := stack.Keys()
	if len(_Keys) <= 0 {
		t.Error("expected non-empty keys")
	}
}
