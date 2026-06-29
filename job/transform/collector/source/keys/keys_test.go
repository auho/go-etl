package keys

import (
	"testing"
)

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

	t.Run("key missing defaults to empty", func(t *testing.T) {
		keys, keysValue, err := k.Contents(map[string]any{"a": "x"})
		if err != nil {
			t.Fatal(err)
		}
		if len(keys) != 2 {
			t.Fatalf("expected 2 keys, got %d", len(keys))
		}
		if keysValue["a"] != "x" {
			t.Fatalf("expected a=x, got %q", keysValue["a"])
		}
		if keysValue["b"] != "" {
			t.Fatalf("expected b empty, got %q", keysValue["b"])
		}
	})

	t.Run("type conversion failure", func(t *testing.T) {
		_, _, err := k.Contents(map[string]any{"a": []string{"not", "scalar"}, "b": "y"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
