package model

import (
	"sort"
	"testing"
)

func TestRule_LabelsName_Sorted(t *testing.T) {
	r := NewRule("test", 30, 30, map[string]int{
		"zebra":  30,
		"alpha":  30,
		"middle": 30,
	}, nil)

	labels := r.LabelsName()

	if !sort.StringsAreSorted(labels) {
		t.Errorf("LabelsName() = %v, want sorted", labels)
	}
}
