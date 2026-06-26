package transform

import (
	"testing"

	"github.com/auho/go-etl/v3/job/transform/collect"
	"github.com/auho/go-etl/v3/job/extract/tag"
)

func Test_UpdateMode(t *testing.T) {
	_mu := NewUpdate(collect.NewKeys([]string{_keyName}), tag.NewMostText(_rule), nil)
	_mu.Prepare()
	defer _mu.Close()

	results := _mu.Apply(_item)
	if len(results) <= 0 {
		t.Error("update most text error")
	}

	_mu2 := NewUpdate(collect.NewKeys([]string{_keyName}), tag.NewMostKey(_rule), nil)
	_mu2.Prepare()
	defer _mu2.Close()

	results2 := _mu2.Apply(_item)
	if len(results2) <= 0 {
		t.Error("update most key error")
	}
}
