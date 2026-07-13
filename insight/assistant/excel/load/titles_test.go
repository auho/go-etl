package load

import (
	"testing"
)

func TestTitlesPrepare_TitlesOnly(t *testing.T) {
	titles := Titles{
		Titles: []string{"a", "b"},
	}
	err := titles.prepare()
	if err != nil {
		t.Fatalf("prepare() error = %v", err)
	}

	names := titles.TitlesName()
	if len(names) != 2 {
		t.Fatalf("TitlesName() len = %d, want 2", len(names))
	}
	if names[0] != "a" || names[1] != "b" {
		t.Errorf("TitlesName() = %v, want [a b]", names)
	}

	indexes := titles.TitlesIndex()
	if len(indexes) != 2 {
		t.Fatalf("TitlesIndex() len = %d, want 2", len(indexes))
	}
	if indexes[0] != 0 || indexes[1] != 1 {
		t.Errorf("TitlesIndex() = %v, want [0 1]", indexes)
	}
}

func TestTitlesPrepare_TitlesWithIndex(t *testing.T) {
	titles := Titles{
		TitlesWithIndex: map[int]string{2: "b", 1: "a"},
	}
	err := titles.prepare()
	if err != nil {
		t.Fatalf("prepare() error = %v", err)
	}

	names := titles.TitlesName()
	if len(names) != 2 {
		t.Fatalf("TitlesName() len = %d, want 2", len(names))
	}
	if names[0] != "a" || names[1] != "b" {
		t.Errorf("TitlesName() = %v, want [a b]", names)
	}

	indexes := titles.TitlesIndex()
	if len(indexes) != 2 {
		t.Fatalf("TitlesIndex() len = %d, want 2", len(indexes))
	}
	if indexes[0] != 0 || indexes[1] != 1 {
		t.Errorf("TitlesIndex() = %v, want [0 1]", indexes)
	}
}

func TestTitlesPrepare_Mixed(t *testing.T) {
	titles := Titles{
		Titles:          []string{"a", "b"},
		TitlesWithIndex: map[int]string{3: "c"},
	}
	err := titles.prepare()
	if err != nil {
		t.Fatalf("prepare() error = %v", err)
	}

	names := titles.TitlesName()
	if len(names) != 3 {
		t.Fatalf("TitlesName() len = %d, want 3", len(names))
	}
	// sorted by index: a(1), b(2), c(3)
	if names[0] != "a" || names[1] != "b" || names[2] != "c" {
		t.Errorf("TitlesName() = %v, want [a b c]", names)
	}

	indexes := titles.TitlesIndex()
	if len(indexes) != 3 {
		t.Fatalf("TitlesIndex() len = %d, want 3", len(indexes))
	}
	if indexes[0] != 0 || indexes[1] != 1 || indexes[2] != 2 {
		t.Errorf("TitlesIndex() = %v, want [0 1 2]", indexes)
	}
}

func TestTitlesPrepare_TitlesWithIndexOverride(t *testing.T) {
	titles := Titles{
		Titles:          []string{"a", "b"},
		TitlesWithIndex: map[int]string{1: "override"},
	}
	err := titles.prepare()
	if err != nil {
		t.Fatalf("prepare() error = %v", err)
	}

	names := titles.TitlesName()
	if len(names) != 2 {
		t.Fatalf("TitlesName() len = %d, want 2", len(names))
	}
	// index 1 is overridden to "override", index 2 stays "b"
	if names[0] != "override" || names[1] != "b" {
		t.Errorf("TitlesName() = %v, want [override b]", names)
	}
}

func TestTitlesPrepare_EmptyError(t *testing.T) {
	titles := Titles{}
	err := titles.prepare()
	if err == nil {
		t.Fatal("prepare() error = nil, want error")
	}
}

func TestTitlesPrepare_ZeroIndexError(t *testing.T) {
	titles := Titles{
		TitlesWithIndex: map[int]string{0: "a"},
	}
	err := titles.prepare()
	if err == nil {
		t.Fatal("prepare() error = nil, want error for zero index")
	}
}
