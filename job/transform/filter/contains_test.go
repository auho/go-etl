package filter

import (
	"testing"
)

func TestNewContainsAll(t *testing.T) {
	pred := NewContainsAll("a", []string{"a1", "a2"})

	// all subs present
	ok, err := pred(map[string]any{"a": "a1a2"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true for all subs present")
	}

	// partial sub present
	ok, err = pred(map[string]any{"a": "a1"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false for partial subs")
	}

	// no sub present
	ok, err = pred(map[string]any{"a": "b"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false for no subs")
	}

	// empty string value
	ok, err = pred(map[string]any{"a": ""})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false for empty string")
	}
}

func TestNewContainsAll_EmptySubs(t *testing.T) {
	// empty subs -> loop does not execute -> true
	pred := NewContainsAll("a", []string{})
	ok, err := pred(map[string]any{"a": "anything"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true for empty subs")
	}
}

func TestNewContainsAll_NonStringType(t *testing.T) {
	// int value is converted to string then checked
	pred := NewContainsAll("a", []string{"12"})
	ok, err := pred(map[string]any{"a": 123})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true for int value containing sub")
	}

	// int value not containing sub
	ok, err = pred(map[string]any{"a": 456})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false for int value not containing sub")
	}
}

func TestNewContainsAll_ErrorOnNilMap(t *testing.T) {
	// nil map -> m[key] is nil -> FromAny(nil) errors -> return error
	pred := NewContainsAll("a", []string{"a1"})

	_, err := pred(nil)
	if err == nil {
		t.Fatal("expected error for nil map")
	}
}

func TestNewContainsAll_ErrorOnMissingKey(t *testing.T) {
	// missing key -> m[key] is nil -> FromAny(nil) errors -> return error
	pred := NewContainsAll("a", []string{"a1"})

	_, err := pred(map[string]any{"b": "x"})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestNewContainsAll_ErrorOnUnsupportedType(t *testing.T) {
	// bool is not supported by FromAny -> return error
	pred := NewContainsAll("a", []string{"a1"})

	_, err := pred(map[string]any{"a": true})
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestNewContainsAny(t *testing.T) {
	pred := NewContainsAny("a", []string{"a1", "a2"})

	// first sub present
	ok, err := pred(map[string]any{"a": "a1b"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true for first sub present")
	}

	// second sub present
	ok, err = pred(map[string]any{"a": "ba2"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true for second sub present")
	}

	// both subs present
	ok, err = pred(map[string]any{"a": "a1a2"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true for both subs present")
	}

	// no sub present
	ok, err = pred(map[string]any{"a": "b"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false for no subs present")
	}

	// empty string value
	ok, err = pred(map[string]any{"a": ""})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false for empty string")
	}
}

func TestNewContainsAny_EmptySubs(t *testing.T) {
	// empty subs -> loop does not execute -> false
	pred := NewContainsAny("a", []string{})
	ok, err := pred(map[string]any{"a": "anything"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false for empty subs")
	}
}

func TestNewContainsAny_NonStringType(t *testing.T) {
	// int value is converted to string then checked
	pred := NewContainsAny("a", []string{"12"})
	ok, err := pred(map[string]any{"a": 123})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected true for int value containing sub")
	}

	// int value not containing sub
	ok, err = pred(map[string]any{"a": 456})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false for int value not containing sub")
	}
}

func TestNewContainsAny_ErrorOnNilMap(t *testing.T) {
	// nil map -> m[key] is nil -> FromAny(nil) errors -> return error
	pred := NewContainsAny("a", []string{"a1"})

	_, err := pred(nil)
	if err == nil {
		t.Fatal("expected error for nil map")
	}
}

func TestNewContainsAny_ErrorOnMissingKey(t *testing.T) {
	// missing key -> m[key] is nil -> FromAny(nil) errors -> return error
	pred := NewContainsAny("a", []string{"a1"})

	_, err := pred(map[string]any{"b": "x"})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}
