package collect

import (
	"testing"
)

func TestContentReader_Content(t *testing.T) {
	c := &ContentReader{}

	t.Run("string", func(t *testing.T) {
		s, err := c.Content("k", map[string]any{"k": "abc"})
		if err != nil {
			t.Fatal(err)
		}
		if s != "abc" {
			t.Fatalf("expected %q, got %q", "abc", s)
		}
	})

	t.Run("int", func(t *testing.T) {
		s, err := c.Content("k", map[string]any{"k": 123})
		if err != nil {
			t.Fatal(err)
		}
		if s != "123" {
			t.Fatalf("expected %q, got %q", "123", s)
		}
	})

	t.Run("int64", func(t *testing.T) {
		s, err := c.Content("k", map[string]any{"k": int64(456)})
		if err != nil {
			t.Fatal(err)
		}
		if s != "456" {
			t.Fatalf("expected %q, got %q", "456", s)
		}
	})

	t.Run("float64", func(t *testing.T) {
		s, err := c.Content("k", map[string]any{"k": 1.5})
		if err != nil {
			t.Fatal(err)
		}
		if s != "1.5" {
			t.Fatalf("expected %q, got %q", "1.5", s)
		}
	})

	t.Run("bytes", func(t *testing.T) {
		s, err := c.Content("k", map[string]any{"k": []byte("abc")})
		if err != nil {
			t.Fatal(err)
		}
		if s != "abc" {
			t.Fatalf("expected %q, got %q", "abc", s)
		}
	})

	t.Run("key not exist", func(t *testing.T) {
		s, err := c.Content("missing", map[string]any{"k": "v"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if s != "" {
			t.Fatalf("expected empty, got %q", s)
		}
	})

	t.Run("unsupported type", func(t *testing.T) {
		s, err := c.Content("k", map[string]any{"k": true})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if s != "" {
			t.Fatalf("expected empty, got %q", s)
		}
	})
}

func TestContentReader_Contents(t *testing.T) {
	c := &ContentReader{}

	t.Run("all convertible", func(t *testing.T) {
		got, err := c.Contents([]string{"a", "b"}, map[string]any{"a": "x", "b": 2})
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 2 {
			t.Fatalf("expected 2, got %d", len(got))
		}
		if got[0] != "x" || got[1] != "2" {
			t.Fatalf("expected [x 2], got %v", got)
		}
	})

	t.Run("middle key unsupported", func(t *testing.T) {
		got, err := c.Contents([]string{"a", "b", "c"}, map[string]any{"a": "x", "b": true, "c": "z"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
	})

	t.Run("empty keys", func(t *testing.T) {
		got, err := c.Contents([]string{}, map[string]any{"a": "x"})
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Fatalf("expected empty slice, got %v", got)
		}
	})
}
