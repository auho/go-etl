package match

import (
	"strings"
	"testing"
)

func Test_fuzzy(t *testing.T) {
	_fuzzy := newFuzzy(1, "a_c", "a_c", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2})

	sr, c, ok := _fuzzy.seeking("acabcabbcabbbc")
	if !ok {
		t.Fatal("")
	}

	if sr.amount != 3 {
		t.Fatal("")
	}

	if strings.Index(c, "abbbc") < 0 {
		t.Fatal("")
	}

	if sr.keyword != "a_c" {
		t.Fatal("")
	}

	_fuzzy = newFuzzy(1, "A_c", "A_c", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2})

	sr, c, ok = _fuzzy.seeking("acAbcabbCabbbcAC")
	if !ok {
		t.Fatal("")
	}

	if sr.amount != 1 {
		t.Fatal("")
	}

	if strings.Index(c, "abbbc") < 0 {
		t.Fatal("")
	}

	if sr.keyword != "A_c" {
		t.Fatal("")
	}

	_fuzzy = newFuzzy(1, "A_c", "a_c", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2})

	sr, c, ok = _fuzzy.seeking("aca一ca二二ca三三三cACcA")
	if !ok {
		t.Fatal("")
	}

	if sr.amount != 3 {
		t.Fatal("")
	}

	if sr.keyword != "A_c" {
		t.Fatal("")
	}

	_fuzzy = newFuzzy(1, "A_c", "a_c_", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2})

	sr, c, ok = _fuzzy.seeking("aca一ca二二ca三三三cACcA")
	if !ok {
		t.Fatal("")
	}

	if sr.amount != 3 {
		t.Fatal("")
	}

	if sr.keyword != "A_c" {
		t.Fatal("")
	}
}
