package match

import (
	"fmt"
	"testing"
)

var _matcherItems = []map[string]string{
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

func TestMatcher(t *testing.T) {
	_m := newMatcher("a", _matcherItems, nil)
	_m.Match(_corpus)
	_m.MatchInTextOrder(_corpus)
	_m.MatchText(_corpus)
	_m.MatchFirstText(_corpus)
	_m.MatchLastText(_corpus)
	_m.MatchMostText(_corpus)
	_m.MatchKey(_corpus)
	_m.MatchFirstKey(_corpus)
	_m.MatchLastKey(_corpus)
	_m.MatchMostKey(_corpus)
	_m.MatchLabel(_corpus)
	_m.MatchLabelMostText(_corpus)
}

func TestMatcher_MatchKey_Accurate(t *testing.T) {
	_m := newMatcher("a", _matcherItems, &matcherConfig{
		ignoreCase:  false,
		mode:        modePriorityAccurate,
		enableFuzzy: true,
		debug:       true,
		fuzzyConfig: FuzzyConfig{
			Window: 3,
			Sep:    "_",
		},
	})

	rets := _m.MatchKey([]string{"ABCDABcAbabacabABBaAc_aE_F_G_e_f_g_h_i_j_H_I_J_iefgAxxciaB"})
	_outputResults(rets)

	_assertResults(t, rets, 6, 12)

	_assertResult(t, rets[0], "ca", 1, 1, map[string]int{"b": 1})
	_assertResult(t, rets[1], "ab", 1, 1, map[string]int{"ab": 1})
	_assertResult(t, rets[2], "a", 4, 4, map[string]int{"a": 1})
	_assertResult(t, rets[3], "A_c", 4, 4, map[string]int{
		"ABc":  1,
		"Ac":   1,
		"Axxc": 1,
	})
	_assertResult(t, rets[4], "E_F_G", 1, 1, map[string]int{"E_F_G": 1})
	_assertResult(t, rets[5], "h_i_j_", 1, 1, map[string]int{"h_i_j_": 1})
}

func TestMatcher_MatchKey_Accurate_IgnoreCase(t *testing.T) {
	_m := newMatcher("a", _matcherItems, &matcherConfig{
		ignoreCase:  true,
		mode:        modePriorityFuzzy,
		enableFuzzy: true,
		debug:       true,
		fuzzyConfig: FuzzyConfig{
			Window: 3,
			Sep:    "_",
		},
	})

	rets := _m.MatchKey([]string{"ABCDABcAbabacabABBaAc_aE_F_G_e_f_g_h_i_j_H_I_J_iefgAxxciaB"})
	_outputResults(rets)

	_assertResults(t, rets, 1, 1)

	_assertResult(t, rets[0], "A_c", 5, 5, map[string]int{
		"ABC":  1,
		"ABc":  1,
		"abac": 1,
		"Ac":   1,
		"aAc":  1,
		"Axxc": 1,
	})
	_assertResult(t, rets[1], "E_F_G", 3, 1, map[string]int{"E_F_G": 3})
	_assertResult(t, rets[2], "h_i_j_", 2, 2, map[string]int{
		"h_i_j_": 1,
		"H_I_J":  1,
	})
	_assertResult(t, rets[3], "ab", 4, 4, map[string]int{
		"Ab": 1,
		"ab": 1,
		"AB": 1,
		"aB": 1,
	})
	_assertResult(t, rets[4], "a", 1, 1, map[string]int{"a": 1})
}

func TestMatcher_MatchKey_Fuzzy(t *testing.T) {
	var rets Results
	_m := newMatcher("a", _matcherItems, &matcherConfig{
		ignoreCase:  true,
		mode:        modePriorityFuzzy,
		enableFuzzy: true,
		debug:       true,
		fuzzyConfig: FuzzyConfig{
			Window: 3,
			Sep:    "_",
		},
	})

	rets = _m.MatchKey([]string{"acAbcabbCAbbbCABbbBc"})
	_outputResults(rets)

	_assertResults(t, rets, 1, 1)

	_assertResult(t, rets[0], "A_c", 4, 4, map[string]int{
		"AbbbC": 1,
		"Abc":   1,
		"abbC":  1,
		"ac":    1,
	})
	_assertResult(t, rets[1], "ab", 1, 1, map[string]int{"AB": 1})

	// double check
	rets = _m.MatchKey([]string{"acAbcabbCAbbbCABbbBc"})
	_outputResults(rets)

	_assertResults(t, rets, 1, 1)

	_assertResult(t, rets[0], "A_c", 4, 4, map[string]int{
		"AbbbC": 1,
		"Abc":   1,
		"abbC":  1,
		"ac":    1,
	})
	_assertResult(t, rets[1], "ab", 1, 1, map[string]int{"AB": 1})

	rets = _m.MatchKey([]string{"efgE一f一二gE一二三FgeF一GE一FGEF一二三四G"})
	_outputResults(rets)

	_assertResults(t, rets, 5, 1)

	_assertResult(t, rets[0], "E_F_G", 5, 5, map[string]int{
		// TODO
	})

	rets = _m.MatchKey([]string{"hijH1ijxH二二IjxxH三三三I123Jh三三三I333JxxxHiJHIJHI四四四四J"})
	_outputResults(rets)

	_assertResults(t, rets, 7, 1)

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

func TestMatcher_MatchFirstText(t *testing.T) {
	var rets Results
	_m := newMatcher("a", _matcherItems, &matcherConfig{
		debug: true,
	})

	rets = _m.MatchFirstText([]string{"abcdef-abcd-abc-ab-a"})
	if rets[0].Keyword != "abcdef" {
		t.Fatal()
	}

	rets = _m.MatchFirstText([]string{"ABCDEF-abCd-Abc-ab-a"})
	if rets[0].Keyword != "ab" {
		t.Fatal()
	}

	rets = _m.MatchFirstText([]string{"aBcdef-aBcd-abc-ab-a"})
	if rets[0].Keyword != "a" {
		t.Fatal()
	}

	rets = _m.MatchFirstText([]string{"babcdefa"})
	if rets[0].Keyword != "abcdef" {
		t.Fatal()
	}

	rets = _m.MatchFirstText([]string{"abcba"})
	if rets[0].Keyword != "abc" {
		t.Fatal()
	}

	rets = _m.MatchFirstText([]string{"caabcdef"})
	if rets[0].Keyword != "ca" {
		t.Fatal()
	}
}
func TestMatcher_MatchKey(t *testing.T) {
	var rets Results
	_m := newMatcher("a", _matcherItems, nil)

	rets = _m.MatchKey([]string{"abcdef-abcd-abc-ab-a"})
	_outputResults(rets)

	_assertResults(t, rets, 5, 5)

	rets = _m.MatchKey([]string{"babcdefa"})
	_outputResults(rets)

	_assertResults(t, rets, 2, 2)

	rets = _m.MatchKey([]string{"abcba"})
	_outputResults(rets)

	_assertResults(t, rets, 2, 2)

	rets = _m.MatchKey([]string{"babcdefa"})
	_outputResults(rets)

	_assertResults(t, rets, 2, 2)

	_m = newMatcher("a", _matcherItems, &matcherConfig{ignoreCase: true})

	rets = _m.MatchKey([]string{"ABCDEF-ABCD-abc-ab-A"})
	_outputResults(rets)

	_assertResults(t, rets, 5, 5)

	rets = _m.MatchKey([]string{"BaBcDeFa"})
	_outputResults(rets)

	_assertResults(t, rets, 2, 2)
}

func TestMatcher_MatchFirstKey(t *testing.T) {
	var rets Results
	_m := newMatcher("a", _matcherItems, &matcherConfig{
		debug: true,
	})

	rets = _m.MatchFirstKey([]string{"abcdef-abcd-abc-ab-a"})
	_outputResults(rets)
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abcdef", 1, 1, map[string]int{"abcdef": 1})

	rets = _m.MatchFirstKey([]string{"ABCDEF-abCd-Abc-ab-a"})
	_outputResults(rets)
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "ab", 1, 1, map[string]int{"ab": 1})

	rets = _m.MatchFirstKey([]string{"aBcdef-aBcd-abc-ab-a"})
	_outputResults(rets)
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abc", 1, 1, map[string]int{"abc": 1})

	rets = _m.MatchFirstKey([]string{"babcdefa"})
	_outputResults(rets)
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abcdef", 1, 1, map[string]int{"abcdef": 1})

	rets = _m.MatchFirstKey([]string{"abcba"})
	_outputResults(rets)
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abc", 1, 1, map[string]int{"abc": 1})

	rets = _m.MatchFirstKey([]string{"caabcdef"})
	_outputResults(rets)
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abcdef", 1, 1, map[string]int{"abcdef": 1})
}

func TestMatcher_MatchLastKey(t *testing.T) {
	var rets Results
	_m := newMatcher("a", _matcherItems, &matcherConfig{
		debug: true,
	})

	rets = _m.MatchLastKey([]string{"abcdef-abcd-abc-ab-a"})
	_outputResults(rets)
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "a", 1, 1, map[string]int{"a": 1})

	rets = _m.MatchLastKey([]string{"ABCDEF-abCd-Abc-ab-a"})
	_outputResults(rets)
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "a", 1, 1, map[string]int{"a": 1})

	rets = _m.MatchLastKey([]string{"aBcdef-aBcd-abc-ab-a"})
	_outputResults(rets)
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "a", 1, 1, map[string]int{"a": 1})

	rets = _m.MatchLastKey([]string{"babcdefa"})
	_outputResults(rets)
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "a", 1, 1, map[string]int{"a": 1})

	rets = _m.MatchLastKey([]string{"abcba"})
	_outputResults(rets)
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "ba", 1, 1, map[string]int{"ba": 1})

	rets = _m.MatchLastKey([]string{"caabcdef"})
	_outputResults(rets)
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "abcdef", 1, 1, map[string]int{"abcdef": 1})
}

func TestMatcher_MatchMostKey(t *testing.T) {
	var rets Results
	_m := newMatcher("a", _matcherItems, &matcherConfig{
		debug: true,
	})

	rets = _m.MatchMostKey([]string{"abcdef-abcd-abc-ab-aabcafasbabcdabcdabefabacabdabadabcdd"})
	_outputResults(rets)
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "a", 1, 1, map[string]int{"a": 1})

	_m = newMatcher("a", _matcherItems, &matcherConfig{
		enableFuzzy: true,
		fuzzyConfig: FuzzyConfig{
			Window: 3,
			Sep:    "_",
		},
		debug: true,
	})

	rets = _m.MatchMostKey([]string{
		"abcdef-abcd-abc-ab-aabcafasbabcdabcdabefabacabdabadabcdd",
		"hijH1ijxH二二IjxxH三三三I123Jh三三三I333JxxxHiJHIJHI四四四四J",
	})
	_outputResults(rets)
	_assertResults(t, rets, 1, 1)
	_assertResult(t, rets[0], "a", 1, 1, map[string]int{"a": 1})
}

func TestMatcher_MatchLabel(t *testing.T) {
	var rets LabelResults
	_m := newMatcher("a", _matcherItems, &matcherConfig{
		debug: true,
	})

	rets = _m.MatchLabel([]string{"abcdef-abcd-abc-ab-a"})
	_outputResults(rets)
	_assertLabelResults(t, rets, 1, 1)

	_assertLabelResult(t, rets[0], "abcdef", 1, 1, 1, nil)

	rets = _m.MatchLabel([]string{"ABCDEF-ABCD-ABC-AB-A"})
	_outputResults(rets)
	_assertLabelResults(t, rets, 0, 0)

	rets = _m.MatchLabel([]string{"aBcdef-aBcd-aBc-ab-a"})
	_outputResults(rets)
	_assertLabelResults(t, rets, 2, 2)

	_assertLabelResult(t, rets[0], "abcdef", 1, 1, 1, nil)

	rets = _m.MatchLabel([]string{"abcabcabc"})
	_outputResults(rets)
	_assertLabelResults(t, rets, 2, 2)

	_assertLabelResult(t, rets[0], "abc", 1, 1, 1, nil)

	rets = _m.MatchLabel([]string{"abccbcbaabccbaabc"})
	_outputResults(rets)
	_assertLabelResults(t, rets, 2, 2)

	_assertLabelResult(t, rets[0], "b", 1, 1, 1, nil)

	_m = newMatcher("a", _matcherItems, &matcherConfig{ignoreCase: true})

	rets = _m.MatchLabel([]string{"ABCDEF-aBcd-Abc-aB-a"})
	_outputResults(rets)
	_assertLabelResults(t, rets, 2, 2)

	_assertLabelResult(t, rets[0], "b", 1, 1, 1, nil)

	rets = _m.MatchLabel([]string{"aBcABCabC"})
	_outputResults(rets)
	_assertLabelResults(t, rets, 2, 2)

	_assertLabelResult(t, rets[0], "b", 1, 1, 1, nil)

	rets = _m.MatchLabel([]string{"abCCbcBAABCCBAabc"})
	_outputResults(rets)
	_assertLabelResults(t, rets, 2, 2)

	_assertLabelResult(t, rets[0], "b", 1, 1, 1, nil)
}

func TestMatcher_MatchLabelMostText(t *testing.T) {
	var rets LabelResults
	_m := newMatcher("a", _matcherItems, &matcherConfig{debug: true})

	rets = _m.MatchLabelMostText([]string{"abcdef-abcd-abc-ab-a"})
	_outputResults(rets)
	_assertLabelResults(t, rets, 2, 2)

	_assertLabelResult(t, rets[0], "abcdef", 1, 1, 1, nil)

	rets = _m.MatchLabelMostText([]string{"acbabcabcabc"})
	_outputResults(rets)
	_assertLabelResults(t, rets, 2, 2)

	_assertLabelResult(t, rets[0], "acb", 1, 1, 1, nil)

	rets = _m.MatchLabelMostText([]string{"cbcbaabcdcbaabc"})
	_outputResults(rets)
	_assertLabelResults(t, rets, 2, 2)

	_assertLabelResult(t, rets[0], "abcd", 1, 1, 1, nil)
}

func _assertResult(t *testing.T, ret Result, keyword string, expectAmount, expectTextsNum int, expectTextsAmount map[string]int) {
	if ret.Keyword != keyword {
		t.Fatal(fmt.Sprintf("result[%s != %s]", keyword, ret.Keyword), t.Name())
	}

	if ret.Amount != expectAmount {
		t.Fatal(fmt.Sprintf("result[%s] amount", keyword), t.Name())
	}

	if expectTextsNum != len(ret.Texts) {
		t.Fatal(fmt.Sprintf("result texts[%s] num[%d != %d]", keyword, expectTextsNum, len(ret.Texts)), t.Name())
	}

	for _t, _a := range expectTextsAmount {
		if ret.Texts[_t] != _a {
			t.Fatal(fmt.Sprintf("result text[%s:%s] amont[%d != %d]", keyword, _t, _a, ret.Texts[_t]), t.Name())
		}
	}
}

func _assertResults(t *testing.T, rets Results, expectTextsAmount int, expectResultsAmount int) {
	if expectResultsAmount != len(rets) {
		t.Fatal("results len", t.Name())
	}

	amount := 0
	textsAmount := 0
	for _, ret := range rets {
		for _, _n := range ret.Texts {
			textsAmount += _n
		}
		amount += ret.Amount
	}

	if amount != textsAmount {
		t.Fatal(fmt.Sprintf("results amount[%d] != texts[%d]", amount, textsAmount), t.Name())
	}

	if amount != expectTextsAmount {
		t.Fatal(fmt.Sprintf("results amount[%d!= %d]", expectTextsAmount, amount), t.Name())
	}

	for _, ret := range rets {
		_a := ret.Amount
		for _, _n := range ret.Texts {
			_a -= _n
		}

		if _a != 0 {
			t.Fatal(fmt.Sprintf("%s amount", ret.Keyword), t.Name())
		}
	}
}

func _assertLabelResult(t *testing.T, ret LabelResult, id string, expectAmount, expectKeysNum, expectTextsNum int, expectTextsAmount map[string]int) {
	if ret.Identity != id {
		t.Fatal(fmt.Sprintf("result[%s != %s]", id, ret.Identity), t.Name())
	}

	if ret.Amount != expectAmount {
		t.Fatal(fmt.Sprintf("result[%s] amount", id), t.Name())
	}

	if expectKeysNum != len(ret.Match) {
		t.Fatal(fmt.Sprintf("result keyword[%s] num[%d != %d]", id, expectKeysNum, len(ret.Match)), t.Name())
	}

	tn := 0
	for _, _kt := range ret.Match {
		tn += len(_kt)
	}
	if tn != expectTextsNum {
		t.Fatal(fmt.Sprintf("result texts[%s] num[%d != %d]", id, expectTextsNum, tn), t.Name())
	}

	for _t, _a := range expectTextsAmount {
		var _ok bool
		for _, _tn := range ret.Match {
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

func _assertLabelResults(t *testing.T, rets LabelResults, expectTextsAmount int, expectResultsAmount int) {
	if expectResultsAmount != len(rets) {
		t.Fatal("label results len", t.Name())
	}

	amount := 0
	textsAmount := 0
	for _, ret := range rets {
		for _, texts := range ret.Match {
			for _, _n := range texts {
				textsAmount += _n
			}
		}
		amount += ret.Amount
	}

	if amount != textsAmount {
		t.Fatal(fmt.Sprintf("label results amount[%d] != texts[%d]", amount, textsAmount), t.Name())
	}

	if amount != expectTextsAmount {
		t.Fatal(fmt.Sprintf("label results amount[%d != %d]", expectTextsAmount, amount), t.Name())
	}

	for _, ret := range rets {
		_a := ret.Amount
		for _, texts := range ret.Match {
			for _, _n := range texts {
				_a -= _n
			}
		}

		if _a != 0 {
			t.Fatal(fmt.Sprintf("%s amount", ret.Identity), t.Name())
		}
	}
}
