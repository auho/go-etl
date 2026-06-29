package transform

import (
	"testing"

	"github.com/auho/go-etl/v3/job/extract/segword"
	"github.com/auho/go-etl/v3/job/transform/collector"
	"github.com/auho/go-etl/v3/job/transform/collector/keys"
	"github.com/auho/go-etl/v3/job/transform/collector/mode"
)

func TestInsertSpread(t *testing.T) {
	ins1 := newTestInsert(t)
	ins2 := NewInsert(collector.NewCollector(keys.New([]string{_keyName}), mode.NewAll(), segword.NewDefault()), nil)
	err := ins2.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	spread := NewInsertSpread(ins1, ins2)
	err = spread.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	ret, err := spread.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) != 1 {
		t.Fatalf("expected 1 result from spread, got %d", len(ret))
	}
}

func TestInsertSpread_SingleInsert(t *testing.T) {
	ins := newTestInsert(t)

	spread := NewInsertSpread(ins)
	err := spread.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	ret, err := spread.Apply(_item)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) != 1 {
		t.Fatalf("expected 1 result from spread, got %d", len(ret))
	}
}