package match

import (
	"fmt"
	"strings"
	"testing"
)

func Test_fuzzy(t *testing.T) {
	var amount int

	_fuzzy := newFuzzy(1, "a_c", "a_c", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2})

	// lowercase
	sr, c, ok := _fuzzy.seeking("acabcabcabbcabbbc", "acabcabcabbcabbbc")
	fmt.Println(c)
	fmt.Println(sr)
	if !ok {
		t.Fatal()
	}

	if sr.amount != 4 {
		t.Fatal()
	}

	if strings.Index(c, "abbbc") < 0 {
		t.Fatal()
	}

	if sr.keyword != "a_c" {
		t.Fatal()
	}

	if sr.textsAmount["abbc"] != 1 || sr.textsAmount["abc"] != 2 || sr.textsAmount["ac"] != 1 {
		t.Fatal()
	}

	amount = 0
	for _, _n := range sr.textsAmount {
		amount += _n
	}
	if amount != sr.amount {
		t.Fatal()
	}

	// uppercase
	_fuzzy = newFuzzy(1, "A_c", "A_c", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2})

	sr, c, ok = _fuzzy.seeking("acAbcabbCabbbcACABC", "acAbcabbCabbbcACABC")
	fmt.Println(c)
	fmt.Println(sr)
	if !ok {
		t.Fatal()
	}

	if sr.amount != 1 {
		t.Fatal()
	}

	if strings.Index(c, "abbbc") < 0 {
		t.Fatal()
	}

	if sr.keyword != "A_c" {
		t.Fatal()
	}

	if sr.textsAmount["Abc"] != 1 {
		t.Fatal()
	}

	// ignore case
	_fuzzy = newFuzzy(1, "A_c", "a_c", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2})

	sr, c, ok = _fuzzy.seeking("acA一ca二二Ca三三三cACcA", "aca一ca二二ca三三三cacca")
	fmt.Println(c)
	fmt.Println(sr)
	if !ok {
		t.Fatal()
	}

	if sr.amount != 3 {
		t.Fatal()
	}

	if sr.keyword != "A_c" {
		t.Fatal()
	}

	if sr.textsAmount["ac"] != 1 || sr.textsAmount["A一c"] != 1 || sr.textsAmount["a二二C"] != 1 {
		t.Fatal()
	}

	_fuzzy = newFuzzy(1, "A_c", "a_c_", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2})

	sr, c, ok = _fuzzy.seeking("aca一ca二二ca三三三cACcA", "aca一ca二二ca三三三cACcA")
	fmt.Println(c)
	fmt.Println(sr)
	if !ok {
		t.Fatal()
	}

	if sr.amount != 3 {
		t.Fatal()
	}

	if sr.keyword != "A_c" {
		t.Fatal()
	}
}
