package filter

import (
	"testing"

	"github.com/auho/go-etl/v3/job/extract"
)

// mockCollector implements collect.Collector for testing.
// Its Extract returns a Result whose IsOK reflects the ok field.
type mockCollector struct {
	ok bool
}

func (m *mockCollector) Title() string                        { return "mock" }
func (m *mockCollector) SourceKeys() []string                  { return nil }
func (m *mockCollector) Extract(item map[string]any, e extract.Extractor) (extract.Result, error) {
	r := extract.Result{}
	if m.ok {
		r.SetOK()
	}
	return r, nil
}

// mockExtractor implements extract.Extractor for testing.
// It is only used to satisfy the NewFilter signature; its methods are not
// exercised because mockCollector.Search ignores the searcher argument.
type mockExtractor struct{}

func (m *mockExtractor) Title() string                           { return "mock" }
func (m *mockExtractor) Prepare() error                          { return nil }
func (m *mockExtractor) NewExport() extract.FieldSpec            { return nil }
func (m *mockExtractor) Search(contents []string) extract.Result { return extract.Result{} }
func (m *mockExtractor) Close() error                            { return nil }

func TestNewMatcher(t *testing.T) {
	// collector reports OK
	pred := NewFilter(&mockCollector{ok: true}, &mockExtractor{})
	ok, err := pred(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true, got false")
	}

	// collector reports not OK
	pred2 := NewFilter(&mockCollector{ok: false}, &mockExtractor{})
	ok, err = pred2(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false, got true")
	}
}

func TestMatcher_OK(t *testing.T) {
	m := &Filter{collector: &mockCollector{ok: true}, extractor: &mockExtractor{}}
	ok, err := m.OK(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true, got false")
	}

	m2 := &Filter{collector: &mockCollector{ok: false}, extractor: &mockExtractor{}}
	ok, err = m2.OK(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false, got true")
	}
}

func TestMatcher_ToPredicate(t *testing.T) {
	m := &Filter{collector: &mockCollector{ok: true}, extractor: &mockExtractor{}}
	pred := m.ToPredicate()
	ok, err := pred(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true, got false")
	}

	m2 := &Filter{collector: &mockCollector{ok: false}, extractor: &mockExtractor{}}
	pred2 := m2.ToPredicate()
	ok, err = pred2(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false, got true")
	}
}

func TestExpression_ToPredicate(t *testing.T) {
	// AND.ToPredicate: true only when all operands match
	andPred := NewAND(opInt(1), opString("1")).ToPredicate()
	ok, err := andPred(map[string]any{"int": 1, "string": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("AND: expected true for all match")
	}
	ok, err = andPred(map[string]any{"int": 2, "string": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("AND: expected false when first operand fails")
	}
	ok, err = andPred(map[string]any{"int": 1, "string": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("AND: expected false when second operand fails")
	}
	ok, err = andPred(map[string]any{"int": 2, "string": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("AND: expected false when all operands fail")
	}

	// OR.ToPredicate: true when any operand matches
	orPred := NewOR(opInt(1), opString("2")).ToPredicate()
	ok, err = orPred(map[string]any{"int": 1, "string": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("OR: expected true when first operand matches")
	}
	ok, err = orPred(map[string]any{"int": 2, "string": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("OR: expected true when second operand matches")
	}
	ok, err = orPred(map[string]any{"int": 1, "string": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("OR: expected true when both operands match")
	}
	ok, err = orPred(map[string]any{"int": 2, "string": "3"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("OR: expected false when no operand matches")
	}
}
