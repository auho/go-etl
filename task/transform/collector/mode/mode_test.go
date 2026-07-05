package mode

import (
	"testing"

	"github.com/auho/go-etl/v3/task/extract"
)

// mockExtractor implements extract.Extractor for testing.
type mockExtractor struct {
	extractFn     func(contents []string) extract.Result
	extractCalled int
}

func (m *mockExtractor) Title() string                 { return "mock" }
func (m *mockExtractor) Prepare() error                { return nil }
func (m *mockExtractor) Keys() []string                { return nil }
func (m *mockExtractor) DefaultValues() map[string]any { return nil }
func (m *mockExtractor) Extract(contents []string) extract.Result {
	m.extractCalled++
	return m.extractFn(contents)
}
func (m *mockExtractor) Close() error { return nil }

func okExtractor(captured *[]string) *mockExtractor {
	return &mockExtractor{
		extractFn: func(contents []string) extract.Result {
			*captured = contents
			return extract.NewResult(true, nil)
		},
	}
}

// --- Semantic A: by key position ---

func TestAll(t *testing.T) {
	var got []string
	e := okExtractor(&got)
	m := NewAll()

	keys := []string{"a", "b"}
	kv := map[string]string{"a": "x", "b": "y"}

	res, err := m.Apply(keys, kv, e)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 2 || got[0] != "x" || got[1] != "y" {
		t.Fatalf("expected [x y], got %v", got)
	}
}

func TestFirst(t *testing.T) {
	var got []string
	e := okExtractor(&got)
	m := NewFirst()

	keys := []string{"a", "b"}
	kv := map[string]string{"a": "x", "b": "y"}

	res, err := m.Apply(keys, kv, e)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 1 || got[0] != "x" {
		t.Fatalf("expected [x], got %v", got)
	}
}

func TestLast(t *testing.T) {
	var got []string
	e := okExtractor(&got)
	m := NewLast()

	keys := []string{"a", "b"}
	kv := map[string]string{"a": "x", "b": "y"}

	res, err := m.Apply(keys, kv, e)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 1 || got[0] != "y" {
		t.Fatalf("expected [y], got %v", got)
	}
}

func TestFirstN(t *testing.T) {
	var got []string
	e := okExtractor(&got)
	m := NewFirstN(2)

	keys := []string{"a", "b", "c"}
	kv := map[string]string{"a": "x", "b": "y", "c": "z"}

	res, err := m.Apply(keys, kv, e)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 2 || got[0] != "x" || got[1] != "y" {
		t.Fatalf("expected [x y], got %v", got)
	}
}

func TestFirstN_Prepare(t *testing.T) {
	if err := NewFirstN(0).Prepare(); err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
	if err := NewFirstN(-1).Prepare(); err == nil {
		t.Fatal("expected error for n=-1, got nil")
	}
	if err := NewFirstN(1).Prepare(); err != nil {
		t.Fatal(err)
	}
}

func TestLastN(t *testing.T) {
	var got []string
	e := okExtractor(&got)
	m := NewLastN(2)

	keys := []string{"a", "b", "c"}
	kv := map[string]string{"a": "x", "b": "y", "c": "z"}

	res, err := m.Apply(keys, kv, e)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 2 || got[0] != "y" || got[1] != "z" {
		t.Fatalf("expected [y z], got %v", got)
	}
}

func TestLastN_Prepare(t *testing.T) {
	if err := NewLastN(0).Prepare(); err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
	if err := NewLastN(1).Prepare(); err != nil {
		t.Fatal(err)
	}
}

func TestMatchAny(t *testing.T) {
	t.Run("first key ok", func(t *testing.T) {
		e := &mockExtractor{extractFn: func(contents []string) extract.Result {
			return extract.NewResult(true, nil)
		}}
		m := NewMatchAny()

		keys := []string{"a", "b"}
		kv := map[string]string{"a": "x", "b": "y"}

		res, err := m.Apply(keys, kv, e)
		if err != nil {
			t.Fatal(err)
		}
		if !res.IsOK() {
			t.Fatal("expected ok")
		}
		if e.extractCalled != 1 {
			t.Fatalf("expected 1 call, got %d", e.extractCalled)
		}
	})

	t.Run("second key ok", func(t *testing.T) {
		calls := 0
		e := &mockExtractor{extractFn: func(contents []string) extract.Result {
			calls++
			if calls == 2 {
				return extract.NewResult(true, nil)
			}
			return extract.NewResult(false, nil)
		}}
		m := NewMatchAny()

		keys := []string{"a", "b"}
		kv := map[string]string{"a": "x", "b": "y"}

		res, err := m.Apply(keys, kv, e)
		if err != nil {
			t.Fatal(err)
		}
		if !res.IsOK() {
			t.Fatal("expected ok")
		}
		if e.extractCalled != 2 {
			t.Fatalf("expected 2 calls, got %d", e.extractCalled)
		}
	})

	t.Run("none ok", func(t *testing.T) {
		e := &mockExtractor{extractFn: func(contents []string) extract.Result {
			return extract.Result{}
		}}
		m := NewMatchAny()

		keys := []string{"a", "b"}
		kv := map[string]string{"a": "x", "b": "y"}

		res, err := m.Apply(keys, kv, e)
		if err != nil {
			t.Fatal(err)
		}
		if res.IsOK() {
			t.Fatal("expected not ok")
		}
		if e.extractCalled != 2 {
			t.Fatalf("expected 2 calls, got %d", e.extractCalled)
		}
	})
}

func TestMatchAnyN(t *testing.T) {
	t.Run("limited to n keys", func(t *testing.T) {
		e := &mockExtractor{extractFn: func(contents []string) extract.Result {
			return extract.NewResult(true, nil)
		}}
		m := NewMatchAnyN(2)

		keys := []string{"a", "b", "c"}
		kv := map[string]string{"a": "x", "b": "y", "c": "z"}

		res, err := m.Apply(keys, kv, e)
		if err != nil {
			t.Fatal(err)
		}
		if !res.IsOK() {
			t.Fatal("expected ok")
		}
		if e.extractCalled != 1 {
			t.Fatalf("expected 1 call (first key ok), got %d", e.extractCalled)
		}
	})

	t.Run("n >= len", func(t *testing.T) {
		e := &mockExtractor{extractFn: func(contents []string) extract.Result {
			return extract.Result{}
		}}
		m := NewMatchAnyN(10)

		keys := []string{"a", "b"}
		kv := map[string]string{"a": "x", "b": "y"}

		_, err := m.Apply(keys, kv, e)
		if err != nil {
			t.Fatal(err)
		}
		if e.extractCalled != 2 {
			t.Fatalf("expected 2 calls (all keys), got %d", e.extractCalled)
		}
	})
}

func TestMatchAnyN_Prepare(t *testing.T) {
	if err := NewMatchAnyN(0).Prepare(); err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
	if err := NewMatchAnyN(-1).Prepare(); err == nil {
		t.Fatal("expected error for n=-1, got nil")
	}
	if err := NewMatchAnyN(1).Prepare(); err != nil {
		t.Fatal(err)
	}
}

// --- Semantic B: by value (only non-empty valued keys) ---

func TestAllValued(t *testing.T) {
	var got []string
	e := okExtractor(&got)
	m := NewAllValued()

	// keys: a="", b="y", c="z" → only b and c are valued
	keys := []string{"a", "b", "c"}
	kv := map[string]string{"a": "", "b": "y", "c": "z"}

	res, err := m.Apply(keys, kv, e)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 2 || got[0] != "y" || got[1] != "z" {
		t.Fatalf("expected [y z], got %v", got)
	}
}

func TestFirstValued(t *testing.T) {
	var got []string
	e := okExtractor(&got)
	m := NewFirstValued()

	// keys: a="", b="y", c="z" → first valued is b
	keys := []string{"a", "b", "c"}
	kv := map[string]string{"a": "", "b": "y", "c": "z"}

	res, err := m.Apply(keys, kv, e)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 1 || got[0] != "y" {
		t.Fatalf("expected [y], got %v", got)
	}
}

func TestLastValued(t *testing.T) {
	var got []string
	e := okExtractor(&got)
	m := NewLastValued()

	// keys: a="x", b="y", c="" → last valued is b
	keys := []string{"a", "b", "c"}
	kv := map[string]string{"a": "x", "b": "y", "c": ""}

	res, err := m.Apply(keys, kv, e)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 1 || got[0] != "y" {
		t.Fatalf("expected [y], got %v", got)
	}
}

func TestFirstValuedN(t *testing.T) {
	var got []string
	e := okExtractor(&got)
	m := NewFirstValuedN(2)

	// keys: a="", b="y", c="z", d="w" → first 2 valued are b, c
	keys := []string{"a", "b", "c", "d"}
	kv := map[string]string{"a": "", "b": "y", "c": "z", "d": "w"}

	res, err := m.Apply(keys, kv, e)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 2 || got[0] != "y" || got[1] != "z" {
		t.Fatalf("expected [y z], got %v", got)
	}
}

func TestFirstValuedN_Prepare(t *testing.T) {
	if err := NewFirstValuedN(0).Prepare(); err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
	if err := NewFirstValuedN(1).Prepare(); err != nil {
		t.Fatal(err)
	}
}

func TestLastValuedN(t *testing.T) {
	var got []string
	e := okExtractor(&got)
	m := NewLastValuedN(2)

	// keys: a="x", b="y", c="", d="w" → last 2 valued are b, d
	keys := []string{"a", "b", "c", "d"}
	kv := map[string]string{"a": "x", "b": "y", "c": "", "d": "w"}

	res, err := m.Apply(keys, kv, e)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 2 || got[0] != "y" || got[1] != "w" {
		t.Fatalf("expected [y w], got %v", got)
	}
}

func TestLastValuedN_Prepare(t *testing.T) {
	if err := NewLastValuedN(0).Prepare(); err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
	if err := NewLastValuedN(1).Prepare(); err != nil {
		t.Fatal(err)
	}
}

func TestMatchAnyValued(t *testing.T) {
	t.Run("skips empty value keys", func(t *testing.T) {
		e := &mockExtractor{extractFn: func(contents []string) extract.Result {
			return extract.NewResult(true, nil)
		}}
		m := NewMatchAnyValued()

		// keys: a="", b="y" → only b is searched
		keys := []string{"a", "b"}
		kv := map[string]string{"a": "", "b": "y"}

		res, err := m.Apply(keys, kv, e)
		if err != nil {
			t.Fatal(err)
		}
		if !res.IsOK() {
			t.Fatal("expected ok")
		}
		if e.extractCalled != 1 {
			t.Fatalf("expected 1 call (skipped empty), got %d", e.extractCalled)
		}
	})

	t.Run("all empty values", func(t *testing.T) {
		e := &mockExtractor{extractFn: func(contents []string) extract.Result {
			return extract.Result{}
		}}
		m := NewMatchAnyValued()

		keys := []string{"a", "b"}
		kv := map[string]string{"a": "", "b": ""}

		res, err := m.Apply(keys, kv, e)
		if err != nil {
			t.Fatal(err)
		}
		if res.IsOK() {
			t.Fatal("expected not ok")
		}
		if e.extractCalled != 0 {
			t.Fatalf("expected 0 calls (all empty), got %d", e.extractCalled)
		}
	})
}

func TestMatchAnyValuedN(t *testing.T) {
	t.Run("limited to n valued keys", func(t *testing.T) {
		e := &mockExtractor{extractFn: func(contents []string) extract.Result {
			return extract.NewResult(true, nil)
		}}
		m := NewMatchAnyValuedN(2)

		// keys: a="", b="y", c="z", d="w" → first 2 valued are b, c; b matches
		keys := []string{"a", "b", "c", "d"}
		kv := map[string]string{"a": "", "b": "y", "c": "z", "d": "w"}

		res, err := m.Apply(keys, kv, e)
		if err != nil {
			t.Fatal(err)
		}
		if !res.IsOK() {
			t.Fatal("expected ok")
		}
		if e.extractCalled != 1 {
			t.Fatalf("expected 1 call (first valued key ok), got %d", e.extractCalled)
		}
	})

	t.Run("n >= valued count", func(t *testing.T) {
		e := &mockExtractor{extractFn: func(contents []string) extract.Result {
			return extract.Result{}
		}}
		m := NewMatchAnyValuedN(10)

		// keys: a="", b="y" → 1 valued key, n=10
		keys := []string{"a", "b"}
		kv := map[string]string{"a": "", "b": "y"}

		_, err := m.Apply(keys, kv, e)
		if err != nil {
			t.Fatal(err)
		}
		if e.extractCalled != 1 {
			t.Fatalf("expected 1 call (only 1 valued key), got %d", e.extractCalled)
		}
	})
}

func TestMatchAnyValuedN_Prepare(t *testing.T) {
	if err := NewMatchAnyValuedN(0).Prepare(); err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
	if err := NewMatchAnyValuedN(1).Prepare(); err != nil {
		t.Fatal(err)
	}
}

// --- helpers ---

func TestValuesByKeys(t *testing.T) {
	keys := []string{"a", "b", "c"}
	kv := map[string]string{"a": "x", "b": "y", "c": "z"}

	got := valuesByKeys(keys, kv)
	if len(got) != 3 || got[0] != "x" || got[1] != "y" || got[2] != "z" {
		t.Fatalf("expected [x y z], got %v", got)
	}
}

func TestValuedKeys(t *testing.T) {
	keys := []string{"a", "b", "c"}
	kv := map[string]string{"a": "", "b": "y", "c": "z"}

	got := valuedKeys(keys, kv)
	if len(got) != 2 || got[0] != "b" || got[1] != "c" {
		t.Fatalf("expected [b c], got %v", got)
	}
}

func TestTakeFirst(t *testing.T) {
	keys := []string{"a", "b", "c"}
	if got := takeFirst(keys, 2); len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("expected [a b], got %v", got)
	}
	if got := takeFirst(keys, 10); len(got) != 3 {
		t.Fatalf("expected all 3, got %v", got)
	}
}

func TestTakeLast(t *testing.T) {
	keys := []string{"a", "b", "c"}
	if got := takeLast(keys, 2); len(got) != 2 || got[0] != "b" || got[1] != "c" {
		t.Fatalf("expected [b c], got %v", got)
	}
	if got := takeLast(keys, 10); len(got) != 3 {
		t.Fatalf("expected all 3, got %v", got)
	}
}
