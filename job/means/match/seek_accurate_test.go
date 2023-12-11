package match

import (
	"strings"
	"testing"
)

func Test_accurate(t *testing.T) {
	_accurate := newAccurate(1, "ac", "ac", map[string]string{"b": "b1"})

	sr, c, ok := _accurate.seeking("acabcabbcabbbc")
	if !ok {
		t.Fatal("")
	}

	if sr.amount != 1 {
		t.Fatal("")
	}

	if strings.Count(c, _placeholder) != 1 {
		t.Fatal("")
	}

	if sr.key != "ac" {
		t.Fatal("")
	}

	_accurate = newAccurate(1, "Ac", "AC", map[string]string{"b": "b1"})

	sr, c, ok = _accurate.seeking("AcaCACaccabbCabbbcAC")
	if !ok {
		t.Fatal("")
	}

	if sr.amount != 2 {
		t.Fatal("")
	}

	if strings.Count(c, _placeholder) != 2 {
		t.Fatal("")
	}

	if sr.key != "Ac" {
		t.Fatal("")
	}
}
