package slices

import (
	"testing"
)

func TestDropDuplicates(t *testing.T) {
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
		{"4", "2", "2"},
		{"5", "2", "3"},
		{"6", "3", "3"},
		{"7", "3", "3"},
		{"8", "3", "4"},
		{"9", "4", "4"},
		{"10", "4", "4"},
	}

	newS2 := SliceSliceDropDuplicates(s2, []int{1, 2})

	_expect = 7
	_actual = len(newS2)
	if _actual != _expect {
		t.Errorf("s2 error %d != %d", _expect, _actual)
	}
}
