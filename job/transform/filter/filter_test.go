package filter

import (
	"testing"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector"
	"github.com/auho/go-etl/v3/job/transform/collector/keys"
	"github.com/auho/go-etl/v3/job/transform/collector/mode"
)

// mockExtractor returns a predictable Result based on ok.
type mockExtractor struct {
	ok bool
}

func (m *mockExtractor) Title() string                                                    { return "mock" }
func (m *mockExtractor) Prepare() error                                                   { return nil }
func (m *mockExtractor) NewExport() extract.FieldSpec                                     { return nil }
func (m *mockExtractor) Search(contents []string) extract.Result {
	r := extract.Result{}
	if m.ok {
		r.SetOK()
	}
	return r
}
func (m *mockExtractor) Close() error { return nil }

func TestNewFilter(t *testing.T) {
	// collector reports OK
	pred := NewFilterPredicate(collector.NewCollector(keys.New([]string{"a"}), mode.NewAll(), &mockExtractor{ok: true}))
	ok, err := pred(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true, got false")
	}

	// collector reports not OK
	pred2 := NewFilterPredicate(collector.NewCollector(keys.New([]string{"a"}), mode.NewAll(), &mockExtractor{ok: false}))
	ok, err = pred2(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false, got true")
	}
}

func TestFilter_OK(t *testing.T) {
	m := &Filter{collector: collector.NewCollector(keys.New([]string{"a"}), mode.NewAll(), &mockExtractor{ok: true})}
	ok, err := m.OK(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true, got false")
	}

	m2 := &Filter{collector: collector.NewCollector(keys.New([]string{"a"}), mode.NewAll(), &mockExtractor{ok: false})}
	ok, err = m2.OK(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false, got true")
	}
}

func TestFilter_ToPredicate(t *testing.T) {
	m := &Filter{collector: collector.NewCollector(keys.New([]string{"a"}), mode.NewAll(), &mockExtractor{ok: true})}
	pred := m.ToPredicate()
	ok, err := pred(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true, got false")
	}

	m2 := &Filter{collector: collector.NewCollector(keys.New([]string{"a"}), mode.NewAll(), &mockExtractor{ok: false})}
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
	// And.ToPredicate: true only when all operands match
	andPred := NewAnd(opInt(1), opString("1")).ToPredicate()
	ok, err := andPred(map[string]any{"int": 1, "string": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("And: expected true for all match")
	}
	ok, err = andPred(map[string]any{"int": 2, "string": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("And: expected false when first operand fails")
	}
	ok, err = andPred(map[string]any{"int": 1, "string": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("And: expected false when second operand fails")
	}
	ok, err = andPred(map[string]any{"int": 2, "string": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("And: expected false when all operands fail")
	}

	// Or.ToPredicate: true when any operand matches
	orPred := NewOr(opInt(1), opString("2")).ToPredicate()
	ok, err = orPred(map[string]any{"int": 1, "string": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("Or: expected true when first operand matches")
	}
	ok, err = orPred(map[string]any{"int": 2, "string": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("Or: expected true when second operand matches")
	}
	ok, err = orPred(map[string]any{"int": 1, "string": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("Or: expected true when both operands match")
	}
	ok, err = orPred(map[string]any{"int": 2, "string": "3"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("Or: expected false when no operand matches")
	}
}
