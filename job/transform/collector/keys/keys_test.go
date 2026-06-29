package keys

import (
	"testing"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector/mode"
)

// mockExtractor implements extract.Extractor for testing.
type mockExtractor struct {
	searchFn func(contents []string) extract.Result
}

func (m *mockExtractor) Title() string                { return "mock" }
func (m *mockExtractor) Prepare() error               { return nil }
func (m *mockExtractor) NewExport() extract.FieldSpec { return nil }
func (m *mockExtractor) Search(contents []string) extract.Result {
	return m.searchFn(contents)
}
func (m *mockExtractor) Close() error { return nil }

func TestKeys_New(t *testing.T) {
	k := New([]string{"a", "b"})
	if k.Title() != "keys{a,b}" {
		t.Fatalf("expected %q, got %q", "keys{a,b}", k.Title())
	}
	got := k.Keys()
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("expected [a b], got %v", got)
	}
}

func TestKeys_Prepare(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		k := New([]string{"a"})
		if err := k.Prepare(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("empty", func(t *testing.T) {
		k := New([]string{})
		if err := k.Prepare(); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestKeys_Contents(t *testing.T) {
	k := New([]string{"a", "b"})

	t.Run("all present", func(t *testing.T) {
		keys, keysValue, err := k.Contents(map[string]any{"a": "x", "b": "y"})
		if err != nil {
			t.Fatal(err)
		}
		if len(keys) != 2 || keys[0] != "a" || keys[1] != "b" {
			t.Fatalf("expected keys [a b], got %v", keys)
		}
		if keysValue["a"] != "x" || keysValue["b"] != "y" {
			t.Fatalf("expected keysValue {a:x b:y}, got %v", keysValue)
		}
	})

	t.Run("key missing", func(t *testing.T) {
		_, _, err := k.Contents(map[string]any{"a": "x"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestMode_All(t *testing.T) {
	var got []string
	e := &mockExtractor{searchFn: func(contents []string) extract.Result {
		got = contents
		r := extract.Result{}
		r.SetOK()
		return r
	}}

	keys := []string{"a", "b"}
	keysValue := map[string]string{"a": "x", "b": "y"}
	m := mode.NewAll()

	res, err := m.Apply(keys, keysValue, e)
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

func TestMode_First(t *testing.T) {
	var got []string
	e := &mockExtractor{searchFn: func(contents []string) extract.Result {
		got = contents
		r := extract.Result{}
		r.SetOK()
		return r
	}}

	keys := []string{"a", "b"}
	keysValue := map[string]string{"a": "x", "b": "y"}
	m := mode.NewFirst()

	res, err := m.Apply(keys, keysValue, e)
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

func TestMode_Last(t *testing.T) {
	var got []string
	e := &mockExtractor{searchFn: func(contents []string) extract.Result {
		got = contents
		r := extract.Result{}
		r.SetOK()
		return r
	}}

	keys := []string{"a", "b"}
	keysValue := map[string]string{"a": "x", "b": "y"}
	m := mode.NewLast()

	res, err := m.Apply(keys, keysValue, e)
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

func TestMode_FirstN(t *testing.T) {
	var got []string
	e := &mockExtractor{searchFn: func(contents []string) extract.Result {
		got = contents
		r := extract.Result{}
		r.SetOK()
		return r
	}}

	keys := []string{"a", "b", "c"}
	keysValue := map[string]string{"a": "x", "b": "y", "c": "z"}
	m := mode.NewFirstN(2)

	res, err := m.Apply(keys, keysValue, e)
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

func TestMode_FirstN_Prepare(t *testing.T) {
	if err := mode.NewFirstN(0).Prepare(); err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
	if err := mode.NewFirstN(-1).Prepare(); err == nil {
		t.Fatal("expected error for n=-1, got nil")
	}
	if err := mode.NewFirstN(1).Prepare(); err != nil {
		t.Fatal(err)
	}
}

func TestMode_LastN(t *testing.T) {
	var got []string
	e := &mockExtractor{searchFn: func(contents []string) extract.Result {
		got = contents
		r := extract.Result{}
		r.SetOK()
		return r
	}}

	keys := []string{"a", "b", "c"}
	keysValue := map[string]string{"a": "x", "b": "y", "c": "z"}
	m := mode.NewLastN(2)

	res, err := m.Apply(keys, keysValue, e)
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

func TestMode_MatchAny(t *testing.T) {
	t.Run("first key ok", func(t *testing.T) {
		calls := 0
		e := &mockExtractor{searchFn: func(contents []string) extract.Result {
			calls++
			r := extract.Result{}
			r.SetOK()
			return r
		}}

		keys := []string{"a", "b"}
		keysValue := map[string]string{"a": "x", "b": "y"}
		m := mode.NewMatchAny()

		res, err := m.Apply(keys, keysValue, e)
		if err != nil {
			t.Fatal(err)
		}
		if !res.IsOK() {
			t.Fatal("expected ok")
		}
		if calls != 1 {
			t.Fatalf("expected 1 call, got %d", calls)
		}
	})

	t.Run("second key ok", func(t *testing.T) {
		calls := 0
		e := &mockExtractor{searchFn: func(contents []string) extract.Result {
			calls++
			r := extract.Result{}
			if calls == 2 {
				r.SetOK()
			}
			return r
		}}

		keys := []string{"a", "b"}
		keysValue := map[string]string{"a": "x", "b": "y"}
		m := mode.NewMatchAny()

		res, err := m.Apply(keys, keysValue, e)
		if err != nil {
			t.Fatal(err)
		}
		if !res.IsOK() {
			t.Fatal("expected ok")
		}
		if calls != 2 {
			t.Fatalf("expected 2 calls, got %d", calls)
		}
	})

	t.Run("none ok", func(t *testing.T) {
		calls := 0
		e := &mockExtractor{searchFn: func(contents []string) extract.Result {
			calls++
			return extract.Result{}
		}}

		keys := []string{"a", "b"}
		keysValue := map[string]string{"a": "x", "b": "y"}
		m := mode.NewMatchAny()

		res, err := m.Apply(keys, keysValue, e)
		if err != nil {
			t.Fatal(err)
		}
		if res.IsOK() {
			t.Fatal("expected not ok")
		}
		if calls != 2 {
			t.Fatalf("expected 2 calls, got %d", calls)
		}
	})
}

func TestMode_MatchAnyN(t *testing.T) {
	t.Run("limited to n keys", func(t *testing.T) {
		calls := 0
		e := &mockExtractor{searchFn: func(contents []string) extract.Result {
			calls++
			r := extract.Result{}
			r.SetOK()
			return r
		}}

		keys := []string{"a", "b", "c"}
		keysValue := map[string]string{"a": "x", "b": "y", "c": "z"}
		m := mode.NewMatchAnyN(2)

		res, err := m.Apply(keys, keysValue, e)
		if err != nil {
			t.Fatal(err)
		}
		if !res.IsOK() {
			t.Fatal("expected ok")
		}
		if calls != 1 {
			t.Fatalf("expected 1 call (first key ok), got %d", calls)
		}
	})

	t.Run("n >= len", func(t *testing.T) {
		calls := 0
		e := &mockExtractor{searchFn: func(contents []string) extract.Result {
			calls++
			return extract.Result{}
		}}

		keys := []string{"a", "b"}
		keysValue := map[string]string{"a": "x", "b": "y"}
		m := mode.NewMatchAnyN(10)

		_, err := m.Apply(keys, keysValue, e)
		if err != nil {
			t.Fatal(err)
		}
		if calls != 2 {
			t.Fatalf("expected 2 calls (all keys), got %d", calls)
		}
	})
}

func TestMode_MatchAnyN_Prepare(t *testing.T) {
	if err := mode.NewMatchAnyN(0).Prepare(); err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
	if err := mode.NewMatchAnyN(-1).Prepare(); err == nil {
		t.Fatal("expected error for n=-1, got nil")
	}
	if err := mode.NewMatchAnyN(1).Prepare(); err != nil {
		t.Fatal(err)
	}
}
