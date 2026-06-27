package transform

import (
	"testing"
)

func TestNewTransfer(t *testing.T) {
	tm := NewTransfer([]string{"a", "b"}, map[string]string{"a": "x"}, map[string]any{"c": "d"})
	if tm == nil {
		t.Fatal("expected non-nil Transfer")
	}
}

func TestTransfer_GetFields(t *testing.T) {
	tm := NewTransfer([]string{"a", "b"}, nil, nil)
	fields := tm.GetFields()
	if len(fields) != 2 || fields[0] != "a" || fields[1] != "b" {
		t.Fatalf("expected [a b], got %v", fields)
	}
}

func TestTransfer_Apply(t *testing.T) {
	t.Run("alias and fixed", func(t *testing.T) {
		tm := NewTransfer(
			[]string{"a", "b"},
			map[string]string{"a": "x"},
			map[string]any{"c": "d"},
		)

		item := map[string]any{"a": "1", "b": "2"}
		result, err := tm.Apply(item)
		if err != nil {
			t.Fatal(err)
		}

		// a is aliased to x
		if result["x"] != "1" {
			t.Fatalf("expected x=1, got %v", result["x"])
		}
		// b stays as b
		if result["b"] != "2" {
			t.Fatalf("expected b=2, got %v", result["b"])
		}
		// fixed key c is added
		if result["c"] != "d" {
			t.Fatalf("expected c=d, got %v", result["c"])
		}
		// original a should not exist
		if _, ok := result["a"]; ok {
			t.Fatal("expected key 'a' to be absent")
		}
	})

	t.Run("alias overlaps fixed", func(t *testing.T) {
		tm := NewTransfer(
			[]string{"a"},
			map[string]string{"a": "x"},
			map[string]any{"x": "y"},
		)

		item := map[string]any{"a": "1"}
		result, err := tm.Apply(item)
		if err != nil {
			t.Fatal(err)
		}

		// fixed x takes from alias target x
		if result["x"] != "1" {
			t.Fatalf("expected x=1 (from alias), got %v", result["x"])
		}
	})

	t.Run("no alias or fixed", func(t *testing.T) {
		tm := NewTransfer([]string{"a", "b"}, nil, nil)

		item := map[string]any{"a": "1", "b": "2"}
		result, err := tm.Apply(item)
		if err != nil {
			t.Fatal(err)
		}

		if result["a"] != "1" || result["b"] != "2" {
			t.Fatalf("expected [a=1 b=2], got %v", result)
		}
	})
}

func TestTransfer_Title(t *testing.T) {
	tm := NewTransfer([]string{"a"}, nil, nil)
	title := tm.Title()
	if title == "" {
		t.Fatal("expected non-empty title")
	}
}

func TestTransfer_Prepare(t *testing.T) {
	tm := NewTransfer([]string{"a"}, nil, nil)
	err := tm.Prepare()
	if err != nil {
		t.Fatal(err)
	}
}

func TestTransfer_Close(t *testing.T) {
	tm := NewTransfer([]string{"a"}, nil, nil)
	err := tm.Close()
	if err != nil {
		t.Fatal(err)
	}
}