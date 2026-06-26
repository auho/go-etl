package filter

import (
	"testing"
)

func TestNewContainsAll(t *testing.T) {
	pred := NewContainsAll("a", []string{"a1", "a2"})

	// all subs present
	if !pred(map[string]any{"a": "a1a2"}) {
		t.Fatal("expected true for all subs present")
	}

	// partial sub present
	if pred(map[string]any{"a": "a1"}) {
		t.Fatal("expected false for partial subs")
	}

	// no sub present
	if pred(map[string]any{"a": "b"}) {
		t.Fatal("expected false for no subs")
	}

	// empty string value
	if pred(map[string]any{"a": ""}) {
		t.Fatal("expected false for empty string")
	}
}

func TestNewContainsAll_EmptySubs(t *testing.T) {
	// empty subs -> loop does not execute -> true
	pred := NewContainsAll("a", []string{})
	if !pred(map[string]any{"a": "anything"}) {
		t.Fatal("expected true for empty subs")
	}
}

func TestNewContainsAll_NonStringType(t *testing.T) {
	// int value is converted to string then checked
	pred := NewContainsAll("a", []string{"12"})
	if !pred(map[string]any{"a": 123}) {
		t.Fatal("expected true for int value containing sub")
	}

	// int value not containing sub
	if pred(map[string]any{"a": 456}) {
		t.Fatal("expected false for int value not containing sub")
	}
}

func TestNewContainsAll_PanicOnNilMap(t *testing.T) {
	// nil map -> m[key] is nil -> FromAny(nil) errors -> panic
	pred := NewContainsAll("a", []string{"a1"})

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil map")
		}
	}()
	pred(nil)
}

func TestNewContainsAll_PanicOnMissingKey(t *testing.T) {
	// missing key -> m[key] is nil -> FromAny(nil) errors -> panic
	pred := NewContainsAll("a", []string{"a1"})

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for missing key")
		}
	}()
	pred(map[string]any{"b": "x"})
}

func TestNewContainsAll_PanicOnUnsupportedType(t *testing.T) {
	// bool is not supported by FromAny -> panic
	pred := NewContainsAll("a", []string{"a1"})

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unsupported type")
		}
	}()
	pred(map[string]any{"a": true})
}

func TestNewContainsAny(t *testing.T) {
	pred := NewContainsAny("a", []string{"a1", "a2"})

	// first sub present
	if !pred(map[string]any{"a": "a1b"}) {
		t.Fatal("expected true for first sub present")
	}

	// second sub present
	if !pred(map[string]any{"a": "ba2"}) {
		t.Fatal("expected true for second sub present")
	}

	// both subs present
	if !pred(map[string]any{"a": "a1a2"}) {
		t.Fatal("expected true for both subs present")
	}

	// no sub present
	if pred(map[string]any{"a": "b"}) {
		t.Fatal("expected false for no subs present")
	}

	// empty string value
	if pred(map[string]any{"a": ""}) {
		t.Fatal("expected false for empty string")
	}
}

func TestNewContainsAny_EmptySubs(t *testing.T) {
	// empty subs -> loop does not execute -> false
	pred := NewContainsAny("a", []string{})
	if pred(map[string]any{"a": "anything"}) {
		t.Fatal("expected false for empty subs")
	}
}

func TestNewContainsAny_NonStringType(t *testing.T) {
	// int value is converted to string then checked
	pred := NewContainsAny("a", []string{"12"})
	if !pred(map[string]any{"a": 123}) {
		t.Fatal("expected true for int value containing sub")
	}

	// int value not containing sub
	if pred(map[string]any{"a": 456}) {
		t.Fatal("expected false for int value not containing sub")
	}
}

func TestNewContainsAny_PanicOnNilMap(t *testing.T) {
	// nil map -> m[key] is nil -> FromAny(nil) errors -> panic
	pred := NewContainsAny("a", []string{"a1"})

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil map")
		}
	}()
	pred(nil)
}

func TestNewContainsAny_PanicOnMissingKey(t *testing.T) {
	// missing key -> m[key] is nil -> FromAny(nil) errors -> panic
	pred := NewContainsAny("a", []string{"a1"})

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for missing key")
		}
	}()
	pred(map[string]any{"b": "x"})
}
