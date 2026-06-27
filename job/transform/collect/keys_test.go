package collect

import (
	"testing"

	"github.com/auho/go-etl/v3/job/extract"
)

// mockExtractor implements extract.Extractor for testing.
// searchFn controls the Result returned by Search and lets tests capture inputs.
type mockExtractor struct {
	searchFn func(contents []string) extract.Result
}

func (m *mockExtractor) Title() string                        { return "mock" }
func (m *mockExtractor) Prepare() error                       { return nil }
func (m *mockExtractor) NewExport() extract.FieldSpec         { return nil }
func (m *mockExtractor) Search(contents []string) extract.Result {
	return m.searchFn(contents)
}
func (m *mockExtractor) Close() error { return nil }

func TestNewKeys(t *testing.T) {
	keys := []string{"a", "b"}
	k := NewKeys(keys)

	if !k.IsAll() {
		t.Fatal("expected IsAll true")
	}
	if k.IsAny() {
		t.Fatal("expected IsAny false")
	}
	got := k.SourceKeys()
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("expected %v, got %v", keys, got)
	}
}

func TestNewKeysAny(t *testing.T) {
	keys := []string{"a", "b"}
	k := NewKeysAny(keys)

	if !k.IsAny() {
		t.Fatal("expected IsAny true")
	}
	if k.IsAll() {
		t.Fatal("expected IsAll false")
	}
}

func TestKeys_Title(t *testing.T) {
	t.Run("multi", func(t *testing.T) {
		k := NewKeys([]string{"a", "b", "c"})
		if k.Title() != "keys{a,b,c}" {
			t.Fatalf("expected %q, got %q", "keys{a,b,c}", k.Title())
		}
	})

	t.Run("single", func(t *testing.T) {
		k := NewKeys([]string{"a"})
		if k.Title() != "keys{a}" {
			t.Fatalf("expected %q, got %q", "keys{a}", k.Title())
		}
	})
}

func TestKeys_Extract_All(t *testing.T) {
	t.Run("all keys present", func(t *testing.T) {
		var got []string
		e := &mockExtractor{searchFn: func(contents []string) extract.Result {
			got = contents
			r := extract.Result{}
			r.SetOK()
			return r
		}}

		k := NewKeys([]string{"a", "b"})
		res, err := k.Extract(map[string]any{"a": "x", "b": "y"}, e)
		if err != nil {
			t.Fatal(err)
		}
		if !res.IsOK() {
			t.Fatal("expected ok")
		}
		if len(got) != 2 || got[0] != "x" || got[1] != "y" {
			t.Fatalf("expected [x y], got %v", got)
		}
	})

	t.Run("key missing", func(t *testing.T) {
		e := &mockExtractor{searchFn: func(contents []string) extract.Result {
			t.Fatal("Search should not be called")
			return extract.Result{}
		}}

		k := NewKeys([]string{"a", "b"})
		_, err := k.Extract(map[string]any{"a": "x"}, e)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestKeys_Extract_Any(t *testing.T) {
	t.Run("first key ok", func(t *testing.T) {
		calls := 0
		e := &mockExtractor{searchFn: func(contents []string) extract.Result {
			calls++
			r := extract.Result{}
			r.SetOK()
			return r
		}}

		k := NewKeysAny([]string{"a", "b"})
		res, err := k.Extract(map[string]any{"a": "x", "b": "y"}, e)
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

		k := NewKeysAny([]string{"a", "b"})
		res, err := k.Extract(map[string]any{"a": "x", "b": "y"}, e)
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

		k := NewKeysAny([]string{"a", "b"})
		res, err := k.Extract(map[string]any{"a": "x", "b": "y"}, e)
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

	t.Run("key missing", func(t *testing.T) {
		e := &mockExtractor{searchFn: func(contents []string) extract.Result {
			t.Fatal("Search should not be called")
			return extract.Result{}
		}}

		k := NewKeysAny([]string{"a", "b"})
		_, err := k.Extract(map[string]any{}, e)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
