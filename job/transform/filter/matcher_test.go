package filter

import (
	"testing"

	"github.com/auho/go-etl/v3/job/extract"
)

// mockCollector implements collect.Collector for testing.
// Its Search returns a Result whose IsOK reflects the ok field.
type mockCollector struct {
	ok bool
}

func (m *mockCollector) Title() string     { return "mock" }
func (m *mockCollector) Keys() []string     { return nil }
func (m *mockCollector) Search(item map[string]any, search extract.Extractor) extract.Result {
	r := extract.Result{}
	if m.ok {
		r.SetOK()
	}
	return r
}

// mockExtractor implements extract.Extractor for testing.
// It is only used to satisfy the NewMatcher signature; its methods are not
// exercised because mockCollector.Search ignores the searcher argument.
type mockExtractor struct{}

func (m *mockExtractor) Title() string                              { return "mock" }
func (m *mockExtractor) Prepare() error                              { return nil }
func (m *mockExtractor) NewExport() extract.FieldSpec                { return nil }
func (m *mockExtractor) Search(contents []string) extract.Result      { return extract.Result{} }
func (m *mockExtractor) Close() error                                { return nil }

func TestNewMatcher(t *testing.T) {
	// collector reports OK
	pred := NewMatcher(&mockCollector{ok: true}, &mockExtractor{})
	if !pred(map[string]any{"a": 1}) {
		t.Fatal("expected true, got false")
	}

	// collector reports not OK
	pred2 := NewMatcher(&mockCollector{ok: false}, &mockExtractor{})
	if pred2(map[string]any{"a": 1}) {
		t.Fatal("expected false, got true")
	}
}

func TestMatcher_OK(t *testing.T) {
	m := &Matcher{collect: &mockCollector{ok: true}, search: &mockExtractor{}}
	if !m.OK(map[string]any{"a": 1}) {
		t.Fatal("expected true, got false")
	}

	m2 := &Matcher{collect: &mockCollector{ok: false}, search: &mockExtractor{}}
	if m2.OK(map[string]any{"a": 1}) {
		t.Fatal("expected false, got true")
	}
}

func TestMatcher_ToPredicate(t *testing.T) {
	m := &Matcher{collect: &mockCollector{ok: true}, search: &mockExtractor{}}
	pred := m.ToPredicate()
	if !pred(map[string]any{"a": 1}) {
		t.Fatal("expected true, got false")
	}

	m2 := &Matcher{collect: &mockCollector{ok: false}, search: &mockExtractor{}}
	pred2 := m2.ToPredicate()
	if pred2(map[string]any{"a": 1}) {
		t.Fatal("expected false, got true")
	}
}

func TestExpression_ToPredicate(t *testing.T) {
	// AND.ToPredicate: true only when all operands match
	andPred := NewAND(opInt(1), opString("1")).ToPredicate()
	if !andPred(map[string]any{"int": 1, "string": "1"}) {
		t.Fatal("AND: expected true for all match")
	}
	if andPred(map[string]any{"int": 2, "string": "1"}) {
		t.Fatal("AND: expected false when first operand fails")
	}
	if andPred(map[string]any{"int": 1, "string": "2"}) {
		t.Fatal("AND: expected false when second operand fails")
	}
	if andPred(map[string]any{"int": 2, "string": "2"}) {
		t.Fatal("AND: expected false when all operands fail")
	}

	// OR.ToPredicate: true when any operand matches
	orPred := NewOR(opInt(1), opString("2")).ToPredicate()
	if !orPred(map[string]any{"int": 1, "string": "1"}) {
		t.Fatal("OR: expected true when first operand matches")
	}
	if !orPred(map[string]any{"int": 2, "string": "2"}) {
		t.Fatal("OR: expected true when second operand matches")
	}
	if !orPred(map[string]any{"int": 1, "string": "2"}) {
		t.Fatal("OR: expected true when both operands match")
	}
	if orPred(map[string]any{"int": 2, "string": "3"}) {
		t.Fatal("OR: expected false when no operand matches")
	}
}
