package match

import (
	"fmt"
	"strings"
	"testing"
)

func Test_accurate(t *testing.T) {
	var amount int
	_accurate := newAccurate(1, "ac", "ac", map[string]string{"b": "b1"})

	// lowercase
	sr, sc, ok := _accurate.seeking("acabcabbcabbbcac", "acabcabbcabbbcac")
	fmt.Println(sc)
	fmt.Println(sr)
	if !ok {
		t.Fatal()
	}

	if strings.Count(sc.origin, _placeholder) != strings.Count(sc.content, _placeholder) ||
		strings.Count(sc.content, _placeholder) != sr.amount {
		t.Fatal()
	}

	if sr.amount != 2 {
		t.Fatal()
	}

	if sr.keyword != "ac" {
		t.Fatal()
	}

	if sr.textsAmount["ac"] != 2 {
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
	_accurate = newAccurate(1, "Ac", "AC", map[string]string{"b": "b1"})

	sr, sc, ok = _accurate.seeking("AcaCACaccabbCabbbcAC", "AcaCACaccabbCabbbcAC")
	fmt.Println(sc)
	fmt.Println(sr)
	if !ok {
		t.Fatal()
	}

	if strings.Count(sc.origin, _placeholder) != strings.Count(sc.content, _placeholder) ||
		strings.Count(sc.content, _placeholder) != sr.amount {
		t.Fatal()
	}

	if sr.amount != 2 {
		t.Fatal()
	}

	if sr.keyword != "Ac" {
		t.Fatal()
	}

	if sr.textsAmount["AC"] != 2 {
		t.Fatal()
	}

	// ignore case
	_accurate = newAccurate(1, "Ac", "ac", map[string]string{"b": "b1"})

	sr, sc, ok = _accurate.seeking("AcaCACaccabbCabbbcAC", "acacacaccabbcabbbcac")
	fmt.Println(sc)
	fmt.Println(sr)
	if !ok {
		t.Fatal()
	}

	if strings.Count(sc.origin, _placeholder) != strings.Count(sc.content, _placeholder) ||
		strings.Count(sc.content, _placeholder) != sr.amount {
		t.Fatal()
	}

	if sr.amount != 5 {
		t.Fatal()
	}

	if sr.keyword != "Ac" {
		t.Fatal()
	}

	if len(sr.textsAmount) != 4 {
		t.Fatal()
	}

	if sr.textsAmount["AC"] != 2 || sr.textsAmount["Ac"] != 1 || sr.textsAmount["aC"] != 1 || sr.textsAmount["ac"] != 1 {
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
