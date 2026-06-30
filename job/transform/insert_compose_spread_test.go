package transform

import (
	"testing"

	"github.com/auho/go-etl/v3/job/extract/segword"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

func TestInsertComposeSpread(t *testing.T) {
	ins1 := newTestInsert(t)
	ins2 := NewInsert(collector.NewKeysAll([]string{_keyName}, segword.NewSegWordsAll()), nil)
	err := ins2.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	cs := NewInsertComposeSpread(ins1, ins2)
	err = cs.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	_ = cs.GetFields()
	_ = cs.Keys()
	_ = cs.DefaultValues()

	ret, err := cs.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) != 1 {
		t.Fatalf("expected 1 result from compose spread, got %d", len(ret))
	}
}

func TestInsertComposeSpread_SingleOperator(t *testing.T) {
	ins := newTestInsert(t)

	cs := NewInsertComposeSpread(ins)
	err := cs.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	ret, err := cs.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) <= 0 {
		t.Error("expected results from single-operator compose spread")
	}
}

func TestInsertComposeSpread_State(t *testing.T) {
	ins := newTestInsert(t)

	cs := NewInsertComposeSpread(ins)
	err := cs.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	state := cs.State()
	if len(state) <= 0 {
		t.Error("expected non-empty state")
	}
}

func TestInsertComposeSpread_Close(t *testing.T) {
	ins := newTestInsert(t)

	cs := NewInsertComposeSpread(ins)
	err := cs.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	err = cs.Close()
	if err != nil {
		t.Fatal(err)
	}
}
