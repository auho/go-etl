package slices

import (
	"testing"
)

func TestSliceStringDropDuplicates(t *testing.T) {
	s1 := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 5}
	newS1 := SliceDropDuplicates(s1)

	_expect := 5
	_actual := len(newS1)
	if _actual != _expect {
		t.Errorf("s1 error %d != %d", _expect, _actual)
	}

	s2 := []string{"1", "2", "2", "3", "3", "3", "4", "4", "4", "4", "5", "5", "5", "5", "5"}
	newS2 := SliceDropDuplicates(s2)
	_expect = 5
	_actual = len(newS2)
	if _actual != _expect {
		t.Errorf("s2 error %d != %d", _expect, _actual)
	}
}

func TestSliceSliceDropDuplicates(t *testing.T) {
	s1 := [][]int{
		{1, 2, 3, 4, 1},
		{1, 2, 3, 4, 2},
		{1, 2, 3, 4, 2},
		{1, 2, 3, 4, 3},
		{1, 2, 3, 4, 3},
		{1, 2, 3, 4, 4},
	}

	newS1 := SliceSliceDropDuplicates(s1, []int{4})

	_expect := 4
	_actual := len(newS1)
	if _actual != _expect {
		t.Errorf("s1 error %d != %d", _expect, _actual)
	}

	s2 := [][]string{
		{"1", "1", "1"},
		{"2", "2", "1"},
		{"3", "2", "2"},
		{"4", "2", "2"}, // 1
		{"5", "2", "3"},
		{"6", "3", "3"},
		{"7", "3", "3"}, // 2
		{"8", "3", "4"},
		{"9", "4", "4"},
		{"10", "4", "4"}, // 3
		{"11", "4", "5"},
		{"12", "4", "5"}, // 4
		{"13", "5", "4"},
		{"14", "5", "4"}, // 5
		{"15", "5", "3"},
		{"16", "3", "5"},
	}

	newS2 := SliceSliceDropDuplicates(s2, []int{1, 2})

	_expect = 11
	_actual = len(newS2)
	if _actual != _expect {
		t.Errorf("s2 error %d != %d", _expect, _actual)
	}
}
