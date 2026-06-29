package transform

import (
	"github.com/auho/go-etl/v3/job/extract/tag"
	"github.com/auho/go-etl/v3/job/transform/collector"
	"github.com/auho/go-etl/v3/job/transform/collector/keys"
	"github.com/auho/go-etl/v3/job/transform/collector/mode"
)

// newTestInsert creates a basic *Insert for use in composite insert tests.
func newTestInsert(t testHelper) *Insert {
	t.Helper()

	ins := NewInsert(collector.NewCollector(keys.New([]string{_keyName}), mode.NewAll(), tag.NewKey(_rule)), nil)
	err := ins.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	return ins
}

// testHelper is satisfied by *testing.T and *testing.B.
type testHelper interface {
	Helper()
	Fatal(args ...any)
}