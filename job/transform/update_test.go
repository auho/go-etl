package transform

import (
	"testing"

	"github.com/auho/go-etl/v3/job/extract/tag"
	"github.com/auho/go-etl/v3/job/transform/collector"
	"github.com/auho/go-etl/v3/job/transform/collector/keys"
	"github.com/auho/go-etl/v3/job/transform/collector/mode"
)

func TestUpdate(t *testing.T) {
	_mu := NewUpdate(collector.NewCollector(keys.New([]string{_keyName}), mode.NewAll(), tag.NewMostText(_rule)), nil)
	err := _mu.Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer _mu.Close()

	results, err := _mu.Apply(_item)
	if err != nil {
		t.Fatal("update most text", err)
	}
	if len(results) <= 0 {
		t.Error("update most text error")
	}

	_mu2 := NewUpdate(collector.NewCollector(keys.New([]string{_keyName}), mode.NewAll(), tag.NewMostKey(_rule)), nil)
	err = _mu2.Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer _mu2.Close()

	results2, err := _mu2.Apply(_item)
	if err != nil {
		t.Fatal("update most key", err)
	}
	if len(results2) <= 0 {
		t.Error("update most key error")
	}
}
