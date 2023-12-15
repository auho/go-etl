package match

import (
	"fmt"
	"strings"
	"testing"
)

func Test_fuzzy(t *testing.T) {
	var amount int

	_fuzzy := newFuzzy(1, "a_c", "a_c", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2}, seekConfig{})

	// lowercase
	sr, sc, ok := _fuzzy.seeking("acabcabcabbcabbbc", "acabcabcabbcabbbc")
	fmt.Println(sc)
	fmt.Println(sr)
	if !ok {
		t.Fatal()
	}

	if strings.Count(sc.origin, _placeholder) != strings.Count(sc.content, _placeholder) ||
		strings.Count(sc.content, _placeholder) != sr.amount {
		t.Fatal()
	}

	if sr.amount != 4 {
		t.Fatal()
	}

	if strings.Index(sc.origin, "abbbc") < 0 || strings.Index(sc.content, "abbbc") < 0 {
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
	_fuzzy = newFuzzy(1, "A_c", "A_c", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2}, seekConfig{})

	sr, sc, ok = _fuzzy.seeking("acAbcabbCabbbcACABC", "acAbcabbCabbbcACABC")
	fmt.Println(sc)
	fmt.Println(sr)
	if !ok {
		t.Fatal()
	}

	if strings.Count(sc.origin, _placeholder) != strings.Count(sc.content, _placeholder) ||
		strings.Count(sc.content, _placeholder) != sr.amount {
		t.Fatal()
	}

	if sr.amount != 1 {
		t.Fatal()
	}

	if strings.Index(sc.origin, "abbbc") < 0 || strings.Index(sc.content, "abbbc") < 0 {
		t.Fatal()
	}

	if sr.keyword != "A_c" {
		t.Fatal()
	}

	if sr.textsAmount["Abc"] != 1 {
		t.Fatal()
	}

	// ignore case
	_fuzzy = newFuzzy(1, "A_c", "a_c", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2}, seekConfig{})

	sr, sc, ok = _fuzzy.seeking("acA一ca二二Ca三三三cACcAAcac", "aca一ca二二ca三三三caccaacac")
	fmt.Println(sc)
	fmt.Println(sr)
	if !ok {
		t.Fatal()
	}

	if strings.Count(sc.origin, _placeholder) != strings.Count(sc.content, _placeholder) ||
		strings.Count(sc.content, _placeholder) != sr.amount {
		t.Fatal()
	}

	if sr.amount != 6 {
		t.Fatal()
	}

	if sr.keyword != "A_c" {
		t.Fatal()
	}

	if len(sr.textsAmount) != 5 {
		t.Fatal()
	}

	if sr.textsAmount["AAc"] != 1 || sr.textsAmount["AC"] != 1 || sr.textsAmount["A一c"] != 1 || sr.textsAmount["ac"] != 2 || sr.textsAmount["a二二C"] != 1 {
		t.Fatal()
	}

	_fuzzy = newFuzzy(1, "A_c", "a_c_", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2}, seekConfig{})

	sr, sc, ok = _fuzzy.seeking("aca一ca二二ca三三三cACcac", "aca一ca二二ca三三三cACcac")
	fmt.Println(sc)
	fmt.Println(sr)
	if !ok {
		t.Fatal()
	}

	if strings.Count(sc.origin, _placeholder) != strings.Count(sc.content, _placeholder) ||
		strings.Count(sc.content, _placeholder) != sr.amount {
		t.Fatal()
	}
	if sr.amount != 4 {
		t.Fatal()
	}

	if sr.keyword != "A_c" {
		t.Fatal()
	}

	if len(sr.textsAmount) != 3 {
		t.Fatal()
	}

	if sr.textsAmount["ac"] != 2 || sr.textsAmount["a一c"] != 1 || sr.textsAmount["a二二c"] != 1 {
		t.Fatal()
	}

	amount = 0
	for _, _n := range sr.textsAmount {
		amount += _n
	}
	if amount != sr.amount {
		t.Fatal()
	}
}
