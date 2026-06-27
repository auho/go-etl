package transform

import (
	"github.com/auho/go-etl/v3/job/extract/tag"
	"github.com/auho/go-etl/v3/job/transform/collect"
)

// newTestInsert creates a basic *Insert for use in composite insert tests.
func newTestInsert(t testHelper) *Insert {
	t.Helper()

	ins := NewInsert(collect.NewKeysAll([]string{_keyName}), tag.NewKey(_rule), nil)
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