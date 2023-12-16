package match

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

func _assertSeekResults(t *testing.T, ok bool, sc seekContent, sr seekResults, expectAmount, expectTextsNum int) {
	originPlWidth := 0
	contentPlWidth := 0

	if len(sr) != expectAmount {
		t.Fatal(fmt.Sprintf("amount[%d != %d]", expectAmount, len(sr)), t.Name())
	}

	textsIndex := make(map[string]int)
	for _, _r := range sr {
		if _, ok1 := textsIndex[_r.text]; !ok1 {
			textsIndex[_r.text] += 1
		}

		contentPlWidth += 1

		_len := len(_r.text)
		_runeLen := utf8.RuneCountInString(_r.text)
		_zhLen := (_len - _runeLen) / 2

		originPlWidth += len(_r.text) - _zhLen
	}

	if len(textsIndex) != expectTextsNum {
		t.Fatal(fmt.Sprintf("texts num[%d != %d]", expectTextsNum, len(textsIndex)), t.Name())
	}

	oplw := strings.Count(sc.origin, _placeholder)
	cplw := strings.Count(sc.content, _placeholder)
	if oplw != originPlWidth {
		t.Fatal(fmt.Sprintf("placeholder origin[%d != %d]", oplw, originPlWidth), t.Name())
	}

	if cplw != contentPlWidth {
		t.Fatal(fmt.Sprintf("placeholder content[%d != %d]", cplw, contentPlWidth), t.Name())
	}

	if contentPlWidth != len(sr) {
		t.Fatal(fmt.Sprintf("placeholder[%d != %d]", len(sr), contentPlWidth), t.Name())
	}
}

func _assertSeekResult(t *testing.T, sr seekResult, keyword, text string) {
	if keyword != sr.keyword {
		t.Fatal(fmt.Sprintf("keyword[%s != %s]", keyword, sr.keyword), t.Name())
	}

	if text != sr.text {
		t.Fatal(fmt.Sprintf("text[%s != %s]", text, sr.text), t.Name())
	}
}

func _outputSeekResults(sr seekResults, originSc, sc seekContent) {
	for _, _r := range sr {
		fmt.Println(fmt.Sprintf("%+v", _r))
	}

	for i := 5; i < 50; i += 5 {
		fmt.Printf("%5d", i)
	}

	fmt.Println()
	fmt.Println(originSc.origin)
	fmt.Println(originSc.content)
	fmt.Println(sc.origin)
	fmt.Println(sc.content)
}
