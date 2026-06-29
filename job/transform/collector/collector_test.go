package collector

import (
	"errors"
	"testing"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector/mode"
	"github.com/auho/go-etl/v3/job/transform/collector/source"
)

// --- mocks ---

type mockSource struct {
	title      string
	keys       []string
	prepareErr error
	contentsFn func(item map[string]any) ([]string, map[string]string, error)
}

func (s *mockSource) Title() string { return s.title }
func (s *mockSource) Keys() []string { return s.keys }
func (s *mockSource) Prepare() error { return s.prepareErr }
func (s *mockSource) Contents(item map[string]any) ([]string, map[string]string, error) {
	return s.contentsFn(item)
}

type mockExtractor struct {
	title         string
	prepareErr    error
	closeErr      error
	export        extract.FieldSpec
	searchFn      func(contents []string) extract.Result
	searchCalled  int
}

func (e *mockExtractor) Title() string                   { return e.title }
func (e *mockExtractor) Prepare() error                  { return e.prepareErr }
func (e *mockExtractor) NewExport() extract.FieldSpec    { return e.export }
func (e *mockExtractor) Search(contents []string) extract.Result {
	e.searchCalled++
	return e.searchFn(contents)
}
func (e *mockExtractor) Close() error { return e.closeErr }

type FieldSpecMock struct {
	keys       []string
	defaults   map[string]any
}

func (f FieldSpecMock) Keys() []string           { return f.keys }
func (f FieldSpecMock) DefaultValues() map[string]any { return f.defaults }

// helper: always-ok extractor that records contents
func okExtractor(captured *[]string) *mockExtractor {
	return &mockExtractor{
		searchFn: func(contents []string) extract.Result {
			*captured = contents
			r := extract.Result{}
			r.SetOK()
			return r
		},
	}
}

// helper: build a Collector with mock source/extractor
func newTestCollector(keys []string, item map[string]any, m mode.Mode, ext *mockExtractor) *Collector {
	src := &mockSource{
		title: "testSource",
		keys:  keys,
		contentsFn: func(it map[string]any) ([]string, map[string]string, error) {
			kv := make(map[string]string, len(keys))
			for _, k := range keys {
				if v, ok := it[k]; ok {
					if s, ok := v.(string); ok {
						kv[k] = s
					}
				}
			}
			return keys, kv, nil
		},
	}
	return NewCollector(src, m, ext)
}

// --- NewCollector ---

func TestNewCollector(t *testing.T) {
	src := &mockSource{title: "s", keys: []string{"a"}, contentsFn: func(map[string]any) ([]string, map[string]string, error) {
		return nil, nil, nil
	}}
	ext := &mockExtractor{title: "e"}
	c := NewCollector(src, mode.NewAll(), ext)
	if c == nil {
		t.Fatal("expected non-nil Collector")
	}
}

// --- Title ---

func TestCollector_Title(t *testing.T) {
	src := &mockSource{title: "src", keys: []string{"a"}, contentsFn: func(map[string]any) ([]string, map[string]string, error) {
		return nil, nil, nil
	}}
	ext := &mockExtractor{title: "ext"}
	c := NewCollector(src, mode.NewAll(), ext)

	got := c.Title()
	want := "src | ext"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

// --- Fields ---

func TestCollector_Fields(t *testing.T) {
	keys := []string{"name", "email"}
	src := &mockSource{title: "s", keys: keys, contentsFn: func(map[string]any) ([]string, map[string]string, error) {
		return nil, nil, nil
	}}
	ext := &mockExtractor{}
	c := NewCollector(src, mode.NewAll(), ext)

	got := c.Fields()
	if len(got) != 2 || got[0] != "name" || got[1] != "email" {
		t.Fatalf("expected [name email], got %v", got)
	}
}

// --- Keys ---

func TestCollector_Keys(t *testing.T) {
	t.Run("with export", func(t *testing.T) {
		src := &mockSource{title: "s", keys: []string{"a"}, contentsFn: func(map[string]any) ([]string, map[string]string, error) {
			return nil, nil, nil
		}}
		ext := &mockExtractor{export: FieldSpecMock{keys: []string{"out1", "out2"}}}
		c := NewCollector(src, mode.NewAll(), ext)

		got := c.Keys()
		if len(got) != 2 || got[0] != "out1" || got[1] != "out2" {
			t.Fatalf("expected [out1 out2], got %v", got)
		}
	})

	t.Run("nil export", func(t *testing.T) {
		src := &mockSource{title: "s", keys: []string{"a"}, contentsFn: func(map[string]any) ([]string, map[string]string, error) {
			return nil, nil, nil
		}}
		ext := &mockExtractor{export: nil}
		c := NewCollector(src, mode.NewAll(), ext)

		got := c.Keys()
		if got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
	})
}

// --- DefaultValues ---

func TestCollector_DefaultValues(t *testing.T) {
	defaults := map[string]any{"out1": "", "out2": 0}
	src := &mockSource{title: "s", keys: []string{"a"}, contentsFn: func(map[string]any) ([]string, map[string]string, error) {
		return nil, nil, nil
	}}
	ext := &mockExtractor{export: FieldSpecMock{defaults: defaults}}
	c := NewCollector(src, mode.NewAll(), ext)

	got := c.DefaultValues()
	if len(got) != 2 {
		t.Fatalf("expected 2 defaults, got %d", len(got))
	}
}

// --- Prepare ---

func TestCollector_Prepare(t *testing.T) {
	t.Run("all ok", func(t *testing.T) {
		src := &mockSource{title: "s", keys: []string{"a"}, contentsFn: func(map[string]any) ([]string, map[string]string, error) {
			return nil, nil, nil
		}}
		ext := &mockExtractor{}
		c := NewCollector(src, mode.NewAll(), ext)
		if err := c.Prepare(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("source prepare error", func(t *testing.T) {
		srcErr := errors.New("source fail")
		src := &mockSource{title: "s", keys: []string{"a"}, prepareErr: srcErr, contentsFn: func(map[string]any) ([]string, map[string]string, error) {
			return nil, nil, nil
		}}
		ext := &mockExtractor{}
		c := NewCollector(src, mode.NewAll(), ext)
		err := c.Prepare()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, srcErr) {
			t.Fatalf("expected source error wrapped, got %v", err)
		}
	})

	t.Run("mode prepare error", func(t *testing.T) {
		src := &mockSource{title: "s", keys: []string{"a"}, contentsFn: func(map[string]any) ([]string, map[string]string, error) {
			return nil, nil, nil
		}}
		ext := &mockExtractor{}
		c := NewCollector(src, mode.NewFirstN(0), ext) // n=0 triggers Prepare error
		err := c.Prepare()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("extractor prepare error", func(t *testing.T) {
		extErr := errors.New("extractor fail")
		src := &mockSource{title: "s", keys: []string{"a"}, contentsFn: func(map[string]any) ([]string, map[string]string, error) {
			return nil, nil, nil
		}}
		ext := &mockExtractor{prepareErr: extErr}
		c := NewCollector(src, mode.NewAll(), ext)
		err := c.Prepare()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, extErr) {
			t.Fatalf("expected extractor error wrapped, got %v", err)
		}
	})
}

// --- Close ---

func TestCollector_Close(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		src := &mockSource{title: "s", keys: []string{"a"}, contentsFn: func(map[string]any) ([]string, map[string]string, error) {
			return nil, nil, nil
		}}
		ext := &mockExtractor{}
		c := NewCollector(src, mode.NewAll(), ext)
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error", func(t *testing.T) {
		closeErr := errors.New("close fail")
		src := &mockSource{title: "s", keys: []string{"a"}, contentsFn: func(map[string]any) ([]string, map[string]string, error) {
			return nil, nil, nil
		}}
		ext := &mockExtractor{closeErr: closeErr}
		c := NewCollector(src, mode.NewAll(), ext)
		err := c.Close()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, closeErr) {
			t.Fatalf("expected close error, got %v", err)
		}
	})
}

// --- Extract ---

func TestCollector_Extract(t *testing.T) {
	t.Run("all mode", func(t *testing.T) {
		var got []string
		ext := okExtractor(&got)
		c := newTestCollector([]string{"a", "b"}, map[string]any{"a": "x", "b": "y"}, mode.NewAll(), ext)

		res, err := c.Extract(map[string]any{"a": "x", "b": "y"})
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

	t.Run("contents error", func(t *testing.T) {
		contentsErr := errors.New("contents fail")
		src := &mockSource{
			title: "s",
			keys:  []string{"a"},
			contentsFn: func(map[string]any) ([]string, map[string]string, error) {
				return nil, nil, contentsErr
			},
		}
		ext := &mockExtractor{}
		c := NewCollector(src, mode.NewAll(), ext)

		_, err := c.Extract(map[string]any{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, contentsErr) {
			t.Fatalf("expected contents error wrapped, got %v", err)
		}
	})
}

// --- Compile-time interface checks ---

var _ source.Source = (*mockSource)(nil)

// --- Entry (factory) functions: verify each creates a working Collector ---

func TestEntry_All(t *testing.T) {
	var got []string
	ext := okExtractor(&got)
	c := NewKeysAll([]string{"a", "b"}, ext)
	if err := c.Prepare(); err != nil {
		t.Fatal(err)
	}
	res, err := c.Extract(map[string]any{"a": "x", "b": "y"})
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

func TestEntry_First(t *testing.T) {
	var got []string
	ext := okExtractor(&got)
	c := NewKeysFirst([]string{"a", "b"}, ext)
	res, _ := c.Extract(map[string]any{"a": "x", "b": "y"})
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 1 || got[0] != "x" {
		t.Fatalf("expected [x], got %v", got)
	}
}

func TestEntry_Last(t *testing.T) {
	var got []string
	ext := okExtractor(&got)
	c := NewKeysLast([]string{"a", "b"}, ext)
	res, _ := c.Extract(map[string]any{"a": "x", "b": "y"})
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 1 || got[0] != "y" {
		t.Fatalf("expected [y], got %v", got)
	}
}

func TestEntry_FirstN(t *testing.T) {
	var got []string
	ext := okExtractor(&got)
	c := NewKeysFirstN([]string{"a", "b", "c"}, 2, ext)
	res, _ := c.Extract(map[string]any{"a": "x", "b": "y", "c": "z"})
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 2 || got[0] != "x" || got[1] != "y" {
		t.Fatalf("expected [x y], got %v", got)
	}
}

func TestEntry_LastN(t *testing.T) {
	var got []string
	ext := okExtractor(&got)
	c := NewKeysLastN([]string{"a", "b", "c"}, 2, ext)
	res, _ := c.Extract(map[string]any{"a": "x", "b": "y", "c": "z"})
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 2 || got[0] != "y" || got[1] != "z" {
		t.Fatalf("expected [y z], got %v", got)
	}
}

func TestEntry_MatchAny(t *testing.T) {
	ext := &mockExtractor{searchFn: func(contents []string) extract.Result {
		r := extract.Result{}
		if len(contents) > 0 && contents[0] == "y" {
			r.SetOK()
		}
		return r
	}}
	c := NewKeysMatchAny([]string{"a", "b"}, ext)
	res, _ := c.Extract(map[string]any{"a": "x", "b": "y"})
	if !res.IsOK() {
		t.Fatal("expected ok on second key")
	}
}

func TestEntry_MatchAnyN(t *testing.T) {
	ext := &mockExtractor{searchFn: func(contents []string) extract.Result {
		r := extract.Result{}
		r.SetOK()
		return r
	}}
	c := NewKeysMatchAnyN([]string{"a", "b", "c"}, 2, ext)
	res, _ := c.Extract(map[string]any{"a": "x", "b": "y", "c": "z"})
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if ext.searchCalled != 1 {
		t.Fatalf("expected 1 search call (first key ok), got %d", ext.searchCalled)
	}
}

// --- Entry: Valued modes ---

func TestEntry_AllValued(t *testing.T) {
	var got []string
	ext := okExtractor(&got)
	c := NewKeysAllValued([]string{"a", "b", "c"}, ext)
	res, _ := c.Extract(map[string]any{"b": "y", "c": "z"}) // a missing → empty
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 2 || got[0] != "y" || got[1] != "z" {
		t.Fatalf("expected [y z] (skips empty a), got %v", got)
	}
}

func TestEntry_FirstValued(t *testing.T) {
	var got []string
	ext := okExtractor(&got)
	c := NewKeysFirstValued([]string{"a", "b", "c"}, ext)
	res, _ := c.Extract(map[string]any{"b": "y", "c": "z"})
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 1 || got[0] != "y" {
		t.Fatalf("expected [y] (first valued), got %v", got)
	}
}

func TestEntry_LastValued(t *testing.T) {
	var got []string
	ext := okExtractor(&got)
	c := NewKeysLastValued([]string{"a", "b", "c"}, ext)
	res, _ := c.Extract(map[string]any{"a": "x", "b": "y"})
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 1 || got[0] != "y" {
		t.Fatalf("expected [y] (last valued), got %v", got)
	}
}

func TestEntry_FirstValuedN(t *testing.T) {
	var got []string
	ext := okExtractor(&got)
	c := NewKeysFirstValuedN([]string{"a", "b", "c", "d"}, 2, ext)
	res, _ := c.Extract(map[string]any{"b": "y", "c": "z", "d": "w"})
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 2 || got[0] != "y" || got[1] != "z" {
		t.Fatalf("expected [y z], got %v", got)
	}
}

func TestEntry_LastValuedN(t *testing.T) {
	var got []string
	ext := okExtractor(&got)
	c := NewKeysLastValuedN([]string{"a", "b", "c", "d"}, 2, ext)
	res, _ := c.Extract(map[string]any{"a": "x", "b": "y", "d": "w"})
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if len(got) != 2 || got[0] != "y" || got[1] != "w" {
		t.Fatalf("expected [y w], got %v", got)
	}
}

func TestEntry_MatchAnyValued(t *testing.T) {
	ext := &mockExtractor{searchFn: func(contents []string) extract.Result {
		r := extract.Result{}
		r.SetOK()
		return r
	}}
	c := NewKeysMatchAnyValued([]string{"a", "b"}, ext)
	res, _ := c.Extract(map[string]any{"b": "y"}) // a missing → skipped
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if ext.searchCalled != 1 {
		t.Fatalf("expected 1 search call (skipped empty a), got %d", ext.searchCalled)
	}
}

func TestEntry_MatchAnyValuedN(t *testing.T) {
	ext := &mockExtractor{searchFn: func(contents []string) extract.Result {
		r := extract.Result{}
		r.SetOK()
		return r
	}}
	c := NewKeysMatchAnyValuedN([]string{"a", "b", "c", "d"}, 2, ext)
	res, _ := c.Extract(map[string]any{"b": "y", "c": "z", "d": "w"})
	if !res.IsOK() {
		t.Fatal("expected ok")
	}
	if ext.searchCalled != 1 {
		t.Fatalf("expected 1 search call (first valued key ok), got %d", ext.searchCalled)
	}
}

// --- Entry: Prepare validation for N modes ---

func TestEntry_FirstN_PrepareError(t *testing.T) {
	c := NewKeysFirstN([]string{"a"}, 0, &mockExtractor{})
	if err := c.Prepare(); err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestEntry_LastN_PrepareError(t *testing.T) {
	c := NewKeysLastN([]string{"a"}, 0, &mockExtractor{})
	if err := c.Prepare(); err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestEntry_MatchAnyN_PrepareError(t *testing.T) {
	c := NewKeysMatchAnyN([]string{"a"}, 0, &mockExtractor{})
	if err := c.Prepare(); err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestEntry_FirstValuedN_PrepareError(t *testing.T) {
	c := NewKeysFirstValuedN([]string{"a"}, 0, &mockExtractor{})
	if err := c.Prepare(); err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestEntry_LastValuedN_PrepareError(t *testing.T) {
	c := NewKeysLastValuedN([]string{"a"}, 0, &mockExtractor{})
	if err := c.Prepare(); err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestEntry_MatchAnyValuedN_PrepareError(t *testing.T) {
	c := NewKeysMatchAnyValuedN([]string{"a"}, 0, &mockExtractor{})
	if err := c.Prepare(); err == nil {
		t.Fatal("expected error for n=0")
	}
}
