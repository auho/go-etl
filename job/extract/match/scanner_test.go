package match

import (
	"fmt"
	"testing"
)

var _scannerItems = []map[string]string{
	{"a": "abcdef", "b": "b1"},
	{"a": "abcd", "b": "b2"},
	{"a": "abc", "b": "b3"},
	{"a": "cba", "b": "b3"},
	{"a": "cb", "b": "b3"},
	{"a": "ca", "b": "b4"},
	{"a": "ab", "b": "b5"},
	{"a": "ba", "b": "b6"},
	{"a": "a", "b": "b6"},
	{"a": "A_c", "b": "b7"},
	{"a": "E_F_G", "b": "b8"},
	{"a": "h_i_j_", "b": "b9"},
	{"a": "hij", "b": "b10"},
}

var _corpus = []string{
	"abcdef-abcd-abc-ab-a",
	"a-ab-abc-abcd-abcdef",
	"a-abcdef",
	"abcba",
	"babcdefa",
}

func TestScanner(t *testing.T) {
	_m := newScanner("a", _scannerItems, ScannerConfig{})
	_m.Scan(_corpus)
	_m.ScanInTextOrder(_corpus)
	_m.ScanText(_corpus)
	_m.ScanFirstText(_corpus)
	_m.ScanLastText(_corpus)
	_m.ScanMostText(_corpus)
	_m.ScanKey(_corpus)
	_m.ScanFirstKey(_corpus)
	_m.ScanLastKey(_corpus)
	_m.ScanMostKey(_corpus)
	_m.ScanLabel(_corpus)
	_m.ScanLabelMostText(_corpus)
}

func TestScanner_ScanKey_Accurate(t *testing.T) {
	_m := newScanner("a", _scannerItems, ScannerConfig{
		IgnoreCase: false,
		Mode:       modePriorityAccurate,
		Debug:      true,
		Fuzzy: FuzzyConfig{
			enabled: true,
			Window:  3,
			Sep:     "_",
		},
	})

	rets := _m.ScanKey([]string{"ABCDABcAbabacabABBaAc_aE_F_G_e_f_g_h_i_j_H_I_J_iefgAxxciaB"})
	_outputResults(rets)
	_assertResults(t, rets, 6, 11)

	_assertResult(t, rets[0], "ca", 1, 1, map[string]int{"ca": 1})
	_assertResult(t, rets[1], "ab", 1, 1, map[string]int{"ab": 1})
	_assertResult(t, rets[2], "a", 4, 1, map[string]int{"a": 4})
	_assertResult(t, rets[3], "A_c", 3, 3, map[string]int{
		"ABc":  1,
		"Ac":   1,
		"Axxc": 1,
	})
	_assertResult(t, rets[4], "E_F_G", 1, 1, map[string]int{"E_F_G": 1})
	_assertResult(t, rets[5], "h_i_j_", 1, 1, map[string]int{"h_i_j": 1})
}

func TestScanner_ScanKey_Accurate_IgnoreCase(t *testing.T) {
	_m := newScanner("a", _scannerItems, ScannerConfig{
		IgnoreCase: true,
		Mode:       modePriorityFuzzy,
		Debug:      true,
		Fuzzy: FuzzyConfig{
			enabled: true,
			Window:  3,
			Sep:     "_",
		},
	})

	rets := _m.ScanKey([]string{"ABCDABcAbabacabABBaAc_aE_F_G_e_f_g_h_i_j_H_I_J_iefgAxxciaBAc_aabacaE_F_G_aB"})
	_outputResults(rets)
	_assertResults(t, rets, 5, 19)

	_assertResult(t, rets[0], "A_c", 7, 7, map[string]int{
		"ABC":   1,
		"ABc":   1,
		"abac":  1,
		"aAc":   1,
		"Axxc":  1,
		"aBAc":  1,
		"aabac": 1,
	})
	_assertResult(t, rets[1], "E_F_G", 4, 3, map[string]int{
		"E_F_G": 2,
		"e_f_g": 1,
		"efg":   1,
	})
	_assertResult(t, rets[2], "h_i_j_", 2, 2, map[string]int{
		"h_i_j": 1,
		"H_I_J": 1,
	})
	_assertResult(t, rets[3], "ab", 4, 4, map[string]int{
		"Ab": 1,
		"ab": 1,
		"AB": 1,
		"aB": 1,
	})
	_assertResult(t, rets[4], "a", 2, 1, map[string]int{"a": 2})
}

func TestScanner_ScanKey_Fuzzy(t *testing.T) {
	var rets results
	_m := newScanner("a", _scannerItems, ScannerConfig{
		IgnoreCase: true,
		Mode:       modePriorityFuzzy,
		Debug:      true,
		Fuzzy: FuzzyConfig{
			enabled: true,
			Window:  3,
			Sep:     "_",
		},
	})

	rets = _m.ScanKey([]string{"acAbcabbCAbbbCABbbBc"})
	_outputResults(rets)
	_assertResults(t, rets, 2, 5)

	_assertResult(t, rets[0], "A_c", 4, 4, map[string]int{
		"ac":    1,
		"Abc":   1,
		"abbC":  1,
		"AbbbC": 1,
	})
	_assertResult(t, rets[1], "ab", 1, 1, map[string]int{"AB": 1})

	rets = _m.ScanKey([]string{"efgE一f一gE一f一gE一二三FgeF一GE一FGEF一二三四G"})
	_outputResults(rets)
	_assertResults(t, rets, 1, 6)

	_assertResult(t, rets[0], "E_F_G", 6, 5, map[string]int{
		"efg":    1,
		"E一f一g":  2,
		"E一二三Fg": 1,
		"eF一G":   1,
		"E一FG":   1,
	})

	rets = _m.ScanKey([]string{"hijH1ijxH二二IjxxH三三三I123Jh三三三I333JxxxHiJHIJHI四四四四J"})
	_outputResults(rets)

	_assertResults(t, rets, 1, 7)

	_assertResult(t, rets[0], "h_i_j_", 7, 7, map[string]int{
		"H1ij":      1,
		"HIJ":       1,
		"HiJ":       1,
		"H三三三I123J": 1,
		"H二二Ij":     1,
		"hij":       1,
		"h三三三I333J": 1,
	})
}

func TestScanner_ScanText(t *testing.T) {
	var rets results
	_m := newScanner("a", _scannerItems, ScannerConfig{
		Debug: true,
	})

	rets = _m.ScanText([]string{
		"efgE一f一gE一f一gE一二三FgeF一GE一FGEF一二三四G",
		"acAbcabbCAbbbCABbbBc",
		"hijH1ijxH二二IjxxH三三三I123Jh三三三I333JxxxHiJHIJHI四四四四J",
	})
	_assertResults(t, rets, 3, 3)

	_assertResult(t, rets[0], "a", 1, 1, map[string]int{"a": 1})
	_assertResult(t, rets[1], "ca", 1, 1, map[string]int{"ca": 1})
	_assertResult(t, rets[2], "hij", 1, 1, map[string]int{"hij": 1})

	_m = newScanner("a", _scannerItems, ScannerConfig{
		IgnoreCase: true,
		Fuzzy:      FuzzyConfig{enabled: true},
		Debug:      true,
	})

	rets = _m.ScanText([]string{
		"efgE一f一gE一f一gE一二三FgeF一GE一FGEF一二三四G",
		"acAbcabbCAbbbCABbbBc",
		"hijH1ijxH二二IjxxH三三三I123Jh三三三I333JxxxHiJHIJHI四四四四J",
	})
	_assertResults(t, rets, 16, 18)

	_assertResult(t, rets[0], "E_F_G", 1, 1, map[string]int{"efg": 1})
	_assertResult(t, rets[1], "E_F_G", 2, 1, map[string]int{"E一f一g": 2})
	_assertResult(t, rets[8], "ca", 2, 1, map[string]int{"CA": 2})
	_assertResult(t, rets[15], "h_i_j_", 1, 1, map[string]int{"HIJ": 1})
}

func TestScanner_ScanFirstText(t *testing.T) {
	var rets results
	_m := newScanner("a", _scannerItems, ScannerConfig{
		Debug: true,
	})

	rets = _m.ScanFirstText([]string{"abcdef-abcd-abc-ab-a"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abcdef", 1, 1, map[string]int{"abcdef": 1})

	rets = _m.ScanFirstText([]string{"ABCDEF-abCd-Abc-ab-a"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "ab", 1, 1, map[string]int{"ab": 1})

	rets = _m.ScanFirstText([]string{"aBcdef-aBcd-abc-ab-a"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "a", 1, 1, map[string]int{"a": 1})

	rets = _m.ScanFirstText([]string{"babcdefa"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abcdef", 1, 1, map[string]int{"abcdef": 1})

	rets = _m.ScanFirstText([]string{"abcba"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abc", 1, 1, map[string]int{"abc": 1})

	rets = _m.ScanFirstText([]string{"caabcdefaababcabcdabcdef"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "ca", 1, 1, map[string]int{"ca": 1})

	rets = _m.ScanFirstText([]string{"abcba", "caabcdef"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abc", 1, 1, map[string]int{"abc": 1})
}

func TestScanner_ScanLastText(t *testing.T) {
	var rets results
	_m := newScanner("a", _scannerItems, ScannerConfig{
		Debug: true,
	})

	rets = _m.ScanLastText([]string{"abcdef-abc-ab-a-abcd"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abcd", 1, 1, map[string]int{"abcd": 1})

	rets = _m.ScanLastText([]string{"ABCDEF-ab-a-abCd-Abc"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "ab", 1, 1, map[string]int{"ab": 1})

	rets = _m.ScanLastText([]string{"aBcdef-ab-a-aBcd-abc"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abc", 1, 1, map[string]int{"abc": 1})

	rets = _m.ScanLastText([]string{"babcdefa"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "a", 1, 1, map[string]int{"a": 1})

	rets = _m.ScanLastText([]string{"abcba"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "ba", 1, 1, map[string]int{"ba": 1})

	rets = _m.ScanLastText([]string{"caabcdef"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abcdef", 1, 1, map[string]int{"abcdef": 1})

	rets = _m.ScanLastText([]string{"abcba", "caabcdef"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abcdef", 1, 1, map[string]int{"abcdef": 1})
}

func TestScanner_ScanMostText(t *testing.T) {
	var rets results
	_m := newScanner("a", _scannerItems, ScannerConfig{
		Debug: true,
	})

	rets = _m.ScanMostText([]string{
		"acAbcabbCAbbbCABbbBc",
		"efgE一f一gE一f一gE一二三FgeF一GE一FGEF一二三四G",
		"acAbcabbCAbbbCABbbBc",
		"acAbcabbCAbbbCABbbBc",
	})
	_assertResults(t, rets, 1, 3)

	_assertResult(t, rets[0], "a", 3, 1, map[string]int{"a": 3})

	_m = newScanner("a", _scannerItems, ScannerConfig{
		IgnoreCase: true,
		Fuzzy:      FuzzyConfig{enabled: true},
		Debug:      true,
	})

	rets = _m.ScanMostText([]string{
		"acAbcabbCAbbbCABbbBc",
		"acAbcabbCAbbbCABbbBc",
		"efgE一f一gE一f一gE一二三FgeF一GE一FGEF一二三四G",
		"acAbcabbCAbbbCABbbBc",
	})
	_outputResults(rets)
	_assertResults(t, rets, 1, 6)

	_assertResult(t, rets[0], "ca", 6, 1, map[string]int{"CA": 6})
}

func TestScanner_ScanKey(t *testing.T) {
	var rets results
	_m := newScanner("a", _scannerItems, ScannerConfig{})

	rets = _m.ScanKey([]string{"abc-abcdef-abcd-ab-abcdef-abcd-ab-a"})
	_outputResults(rets)

	_assertResults(t, rets, 5, 8)
	_assertResult(t, rets[0], "abcdef", 2, 1, nil)
	_assertResult(t, rets[1], "abcd", 2, 1, nil)
	_assertResult(t, rets[2], "abc", 1, 1, nil)
	_assertResult(t, rets[3], "ab", 2, 1, nil)
	_assertResult(t, rets[4], "a", 1, 1, nil)

	rets = _m.ScanKey([]string{"babcdefa"})
	_outputResults(rets)

	_assertResults(t, rets, 2, 2)

	rets = _m.ScanKey([]string{"abcba"})
	_outputResults(rets)

	_assertResults(t, rets, 2, 2)

	rets = _m.ScanKey([]string{"babcdefa"})
	_outputResults(rets)

	_assertResults(t, rets, 2, 2)

	_m = newScanner("a", _scannerItems, ScannerConfig{IgnoreCase: true})

	rets = _m.ScanKey([]string{"ABCDEF-ABCD-abc-ab-A"})
	_outputResults(rets)

	_assertResults(t, rets, 5, 5)

	rets = _m.ScanKey([]string{"BaBcDeFa"})
	_outputResults(rets)

	_assertResults(t, rets, 2, 2)
}

func TestScanner_ScanFirstKey(t *testing.T) {
	var rets results
	_m := newScanner("a", _scannerItems, ScannerConfig{
		Debug: true,
	})

	rets = _m.ScanFirstKey([]string{"abc-abcdef-abcd-ab-a"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abcdef", 1, 1, map[string]int{"abcdef": 1})

	rets = _m.ScanFirstKey([]string{"ABCDEF-abCd-Abc-ab-a"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "ab", 1, 1, map[string]int{"ab": 1})

	rets = _m.ScanFirstKey([]string{"aBcdef-aBcd-abc-ab-a"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abc", 1, 1, map[string]int{"abc": 1})

	rets = _m.ScanFirstKey([]string{"babcdefa"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abcdef", 1, 1, map[string]int{"abcdef": 1})

	rets = _m.ScanFirstKey([]string{"abcba"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abc", 1, 1, map[string]int{"abc": 1})

	rets = _m.ScanFirstKey([]string{"caabcdef"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abcdef", 1, 1, map[string]int{"abcdef": 1})
}

func TestScanner_ScanLastKey(t *testing.T) {
	var rets results
	_m := newScanner("a", _scannerItems, ScannerConfig{
		Debug: true,
	})

	rets = _m.ScanLastKey([]string{"abcdef-abcd-abc-ab-a"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "a", 1, 1, map[string]int{"a": 1})

	rets = _m.ScanLastKey([]string{"ABCDEF-ab-a-abCd-Abc"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "a", 1, 1, map[string]int{"a": 1})

	rets = _m.ScanLastKey([]string{"aBcdef-ab-a-aBcd-abc"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "a", 1, 1, map[string]int{"a": 1})

	rets = _m.ScanLastKey([]string{"babcdefa"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "a", 1, 1, map[string]int{"a": 1})

	rets = _m.ScanLastKey([]string{"abcba"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "ba", 1, 1, map[string]int{"ba": 1})

	rets = _m.ScanLastKey([]string{"caabcdef"})
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "ca", 1, 1, map[string]int{"ca": 1})
}

func TestScanner_ScanMostKey(t *testing.T) {
	var rets results
	_m := newScanner("a", _scannerItems, ScannerConfig{
		Debug: true,
	})

	rets = _m.ScanMostKey([]string{"abcdef-abcd-abc-ab-aabcafasbabcdabcdabefabacabdabadabcdd"})
	_assertResults(t, rets, 1, 5)
	_assertResult(t, rets[0], "a", 5, 1, map[string]int{"a": 5})

	_m = newScanner("a", _scannerItems, ScannerConfig{
		IgnoreCase: true,
		Fuzzy: FuzzyConfig{
			enabled: true,
			Window:  3,
			Sep:     "_",
		},
		Debug: true,
	})

	rets = _m.ScanMostKey([]string{
		"hijH1ijxH二二IjxxH三三三I123Jh三三三I333JxxxHiJHIJHI四四四四J",
		"abcdef",
		"hijH1ijxH二二IjxxH三三三I123Jh三三三I333JxxxHiJHIJHI四四四四J",
	})
	_assertResults(t, rets, 1, 14)
	_assertResult(t, rets[0], "h_i_j_", 14, 7, map[string]int{
		"hij":       2,
		"H1ij":      2,
		"H二二Ij":     2,
		"H三三三I123J": 2,
		"h三三三I333J": 2,
		"HiJ":       2,
		"HIJ":       2,
	})
}

func TestScanner_ScanLabel(t *testing.T) {
	var rets labelResults
	_m := newScanner("a", _scannerItems, ScannerConfig{
		Debug: true,
	})

	rets = _m.ScanLabel([]string{"abcdef-abcd-abc-ab-a"})
	_assertLabelResults(t, rets, 5, 5)

	_assertLabelResult(t, rets[0], "-b1", 1, 1, 1, map[string]int{"abcdef": 1})
	_assertLabelResult(t, rets[1], "-b2", 1, 1, 1, map[string]int{"abcd": 1})
	_assertLabelResult(t, rets[2], "-b3", 1, 1, 1, map[string]int{"abc": 1})
	_assertLabelResult(t, rets[3], "-b5", 1, 1, 1, map[string]int{"ab": 1})
	_assertLabelResult(t, rets[4], "-b6", 1, 1, 1, map[string]int{"a": 1})

	_m = newScanner("a", _scannerItems, ScannerConfig{
		IgnoreCase: true,
		Fuzzy: FuzzyConfig{
			enabled: true,
			Window:  3,
			Sep:     "_",
		},
		Debug: true,
	})

	rets = _m.ScanLabel([]string{
		"hijH1ijxH二二IjxxH三三三I123Jh三三三I333JxxxHiJHIJHI四四四四J",
		"abcdef",
		"hijH1ijxH二二IjxxH三三三I123Jh三三三I333JxxxHiJHIJHI四四四四J",
	})
	_outputResults(rets)
	_assertLabelResults(t, rets, 2, 15)

	_assertLabelResult(t, rets[0], "-b9", 14, 1, 7, map[string]int{
		"hij":       2,
		"H1ij":      2,
		"H二二Ij":     2,
		"H三三三I123J": 2,
		"h三三三I333J": 2,
		"HiJ":       2,
		"HIJ":       2,
	})
	_assertLabelResult(t, rets[1], "-b1", 1, 1, 1, map[string]int{
		"abcdef": 1,
	})
}

func TestScanner_ScanLabelMostText(t *testing.T) {
	var rets labelResults
	_m := newScanner("a", _scannerItems, ScannerConfig{Debug: true})

	rets = _m.ScanLabelMostText([]string{"abcdef-abcd-abc-ab-a"})
	_assertLabelResults(t, rets, 1, 1)

	_assertLabelResult(t, rets[0], "-b1", 1, 1, 1, nil)

	rets = _m.ScanLabelMostText([]string{"acbabcabcabc"})
	_assertLabelResults(t, rets, 1, 4)

	_assertLabelResult(t, rets[0], "-b3", 4, 2, 2, map[string]int{
		"abc": 3,
		"cb":  1,
	})
}

func _assertResult(t *testing.T, ret result, keyword string, expectAmount, expectTextsNum int, expectTextsAmount map[string]int) {
	if ret.keyword != keyword {
		t.Fatal(fmt.Sprintf("result[%s != %s]", keyword, ret.keyword), t.Name())
	}

	if ret.amount != expectAmount {
		t.Fatal(fmt.Sprintf("result[%s] amount[%d != %d]", keyword, expectAmount, ret.amount), t.Name())
	}

	if expectTextsNum != len(ret.texts) {
		t.Fatal(fmt.Sprintf("result texts[%s] num[%d != %d]", keyword, expectTextsNum, len(ret.texts)), t.Name())
	}

	for _t, _a := range expectTextsAmount {
		if ret.texts[_t] != _a {
			t.Fatal(fmt.Sprintf("result text[%s:%s] amont[%d != %d]", keyword, _t, _a, ret.texts[_t]), t.Name())
		}
	}
}

func _assertResults(t *testing.T, rets results, expectResultsAmount, expectTextsAmount int) {
	if expectResultsAmount != len(rets) {
		t.Fatal(fmt.Sprintf("results len[%d != %d]", expectResultsAmount, len(rets)), t.Name())
	}

	amount := 0
	textsAmount := 0
	for _, ret := range rets {
		for _, _n := range ret.texts {
			textsAmount += _n
		}
		amount += ret.amount
	}

	if amount != textsAmount {
		t.Fatal(fmt.Sprintf("results amount[%d] != texts[%d]", amount, textsAmount), t.Name())
	}

	if amount != expectTextsAmount {
		t.Fatal(fmt.Sprintf("results amount[%d!= %d]", expectTextsAmount, amount), t.Name())
	}

	for _, ret := range rets {
		_a := ret.amount
		for _, _n := range ret.texts {
			_a -= _n
		}

		if _a != 0 {
			t.Fatal(fmt.Sprintf("%s amount", ret.keyword), t.Name())
		}
	}
}

func _assertLabelResult(t *testing.T, ret labelResult, id string, expectAmount, expectKeysNum, expectTextsNum int, expectTextsAmount map[string]int) {
	if ret.identity != id {
		t.Fatal(fmt.Sprintf("result[%s != %s]", id, ret.identity), t.Name())
	}

	if ret.amount != expectAmount {
		t.Fatal(fmt.Sprintf("result[%s] amount", id), t.Name())
	}

	if expectKeysNum != len(ret.match) {
		t.Fatal(fmt.Sprintf("result keyword[%s] num[%d != %d]", id, expectKeysNum, len(ret.match)), t.Name())
	}

	tn := 0
	for _, _kt := range ret.match {
		tn += len(_kt)
	}
	if tn != expectTextsNum {
		t.Fatal(fmt.Sprintf("result texts[%s] num[%d != %d]", id, expectTextsNum, tn), t.Name())
	}

	for _t, _a := range expectTextsAmount {
		var _ok bool
		for _, _tn := range ret.match {
			if _n, ok := _tn[_t]; ok {
				if _a == _n {
					_ok = true
				}
			}
		}

		if _ok == false {
			if _a == 0 {
				_ok = true
			}
		}

		if !_ok {
			t.Fatal(fmt.Sprintf("result text[%s:%s] amont[%d]", id, _t, _a), t.Name())
		}
	}
}

func _assertLabelResults(t *testing.T, rets labelResults, expectResultsAmount, expectTextsAmount int) {
	if expectResultsAmount != len(rets) {
		t.Fatal(fmt.Sprintf("label results len[%d != %d]", expectResultsAmount, len(rets)), t.Name())
	}

	amount := 0
	textsAmount := 0
	for _, ret := range rets {
		for _, texts := range ret.match {
			for _, _n := range texts {
				textsAmount += _n
			}
		}
		amount += ret.amount
	}

	if amount != textsAmount {
		t.Fatal(fmt.Sprintf("label results amount[%d] != texts[%d]", amount, textsAmount), t.Name())
	}

	if amount != expectTextsAmount {
		t.Fatal(fmt.Sprintf("label results amount[%d != %d]", expectTextsAmount, amount), t.Name())
	}

	for _, ret := range rets {
		_a := ret.amount
		for _, texts := range ret.match {
			for _, _n := range texts {
				_a -= _n
			}
		}

		if _a != 0 {
			t.Fatal(fmt.Sprintf("%s amount", ret.identity), t.Name())
		}
	}
}
