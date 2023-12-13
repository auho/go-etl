package match

import (
	"fmt"
	"strings"
	"testing"
)

func Test_accurate(t *testing.T) {
	_accurate := newAccurate(1, "ac", "ac", map[string]string{"b": "b1"})

	// lowercase
	sr, c, ok := _accurate.seeking("acabcabbcabbbc", "acabcabbcabbbc")
	fmt.Println(c)
	fmt.Println(sr)
	if !ok {
		t.Fatal()
	}

	if sr.amount != 1 {
		t.Fatal()
	}

	if strings.Count(c, _placeholder) != 1 {
		t.Fatal()
	}

	if sr.keyword != "ac" {
		t.Fatal()
	}

	// uppercase
	_accurate = newAccurate(1, "Ac", "AC", map[string]string{"b": "b1"})

	sr, c, ok = _accurate.seeking("AcaCACaccabbCabbbcAC", "AcaCACaccabbCabbbcAC")
	fmt.Println(c)
	fmt.Println(sr)
	if !ok {
		t.Fatal()
	}

	if sr.amount != 2 {
		t.Fatal()
	}

	if strings.Count(c, _placeholder) != 2 {
		t.Fatal()
	}

	if sr.keyword != "Ac" {
		t.Fatal()
	}

	// ignore case
	_accurate = newAccurate(1, "Ac", "ac", map[string]string{"b": "b1"})

	sr, c, ok = _accurate.seeking("AcaCACaccabbCabbbcAC", "acacacaccabbcabbbcac")
	fmt.Println(c)
	fmt.Println(sr)
	if !ok {
		t.Fatal()
	}

	if sr.amount != 5 {
		t.Fatal()
	}

	if strings.Count(c, _placeholder) != 5 {
		t.Fatal()
	}

	if sr.keyword != "Ac" {
		t.Fatal()
	}

	if sr.textsAmount["AC"] != 2 || sr.textsAmount["Ac"] != 1 || sr.textsAmount["aC"] != 1 || sr.textsAmount["ac"] != 1 {
		t.Fatal()
	}
}
