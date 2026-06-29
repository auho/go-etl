package filter

import (
	"errors"
	"testing"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

// mockExtractor returns a predictable Result based on ok.
type mockExtractor struct {
	ok bool
}

func (m *mockExtractor) Title() string                { return "mock" }
func (m *mockExtractor) Prepare() error               { return nil }
func (m *mockExtractor) NewExport() extract.FieldSpec { return nil }
func (m *mockExtractor) Search(contents []string) extract.Result {
	r := extract.Result{}
	if m.ok {
		r.SetOK()
	}
	return r
}
func (m *mockExtractor) Close() error { return nil }

// errorMockExtractor fails on Prepare.
type errorMockExtractor struct{}

func (m *errorMockExtractor) Title() string                { return "mock" }
func (m *errorMockExtractor) Prepare() error               { return errors.New("prepare failed") }
func (m *errorMockExtractor) NewExport() extract.FieldSpec { return nil }
func (m *errorMockExtractor) Search(contents []string) extract.Result {
	return extract.Result{}
}
func (m *errorMockExtractor) Close() error { return nil }

func TestNewFromCollector(t *testing.T) {
	// collector reports OK
	pred := NewFromCollector(collector.NewKeysAll([]string{"a"}, &mockExtractor{ok: true}))
	ok, err := pred.Match(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true, got false")
	}

	// collector reports not OK
	pred2 := NewFromCollector(collector.NewKeysAll([]string{"a"}, &mockExtractor{ok: false}))
	ok, err = pred2.Match(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false, got true")
	}
}

func TestFromCollector_Match(t *testing.T) {
	m := &FromCollector{collector: collector.NewKeysAll([]string{"a"}, &mockExtractor{ok: true})}
	ok, err := m.Match(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true, got false")
	}

	m2 := &FromCollector{collector: collector.NewKeysAll([]string{"a"}, &mockExtractor{ok: false})}
	ok, err = m2.Match(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false, got true")
	}
}

func TestFromCollector_Prepare(t *testing.T) {
	// Prepare delegates to collector.Prepare; with mockExtractor it succeeds.
	m := &FromCollector{collector: collector.NewKeysAll([]string{"a"}, &mockExtractor{ok: true})}
	if err := m.Prepare(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestFromCollector_PrepareError(t *testing.T) {
	// Prepare propagates the collector's Prepare error.
	m := &FromCollector{collector: collector.NewKeysAll([]string{"a"}, &errorMockExtractor{})}
	if err := m.Prepare(); err == nil {
		t.Fatal("expected error, got nil")
	}
}
