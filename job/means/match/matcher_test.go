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
	_m.MatchKey(_corpus)
	_m.MatchMostKey(_corpus)
	_m.MatchText(_corpus)
	_m.MatchMostText(_corpus)
	_m.MatchFirstText(_corpus)
	_m.MatchLastText(_corpus)
	_m.MatchLabel(_corpus)
	_m.MatchLabelMostText(_corpus)
	_m.MatchFirstLabel(_corpus)
}

func TestMatcher_MatchKey_Accurate(t *testing.T) {
	var rets Results
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

	rets = _m.MatchKey([]string{"ABCDABcAbabacabABBaAc_aE_F_G_e_f_g_h_i_j_H_I_J_efgAxxc"})
	if len(rets) != 5 {
		t.Fatal()
	}

	if rets[0].Keyword != "ca" || rets[0].Amount != 1 || rets[0].Texts["ca"] != 1 || rets[0].Tags["b"] != "b4" {
		t.Fatal()
	}

	if rets[1].Keyword != "ab" || rets[1].Amount != 1 {
		t.Fatal()
	}

	if rets[2].Keyword != "a" || rets[2].Amount != 3 || rets[0].Texts["ca"] != 3 || rets[0].Tags["b"] != "b6" {
		t.Fatal()
	}

	if rets[3].Keyword != "E_F_G" || rets[3].Amount != 1 {
		t.Fatal()
	}

	if rets[4].Keyword != "h_i_j_" || rets[4].Amount != 1 {
		t.Fatal()
	}
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
	if len(rets) != 2 {
		t.Fatal()
	}

	if rets[0].Keyword != "A_c" || rets[0].Amount != 4 || rets[0].Tags["b"] != "b7" {
		t.Fatal()
	}

	if rets[0].Texts["Abc"] != 1 {
		t.Fatal()
	}

	amount := 0
	for _, _n := range rets[0].Texts {
		amount += _n
	}
	if amount != rets[0].Amount {
		t.Fatal()
	}

	if rets[1].Keyword != "ab" || rets[1].Amount != 1 {
		t.Fatal()
	}

	// double check
	rets = _m.MatchKey([]string{"acAbcabbCAbbbCABbbBc"})
	if len(rets) != 2 {
		t.Fatal()
	}

	if rets[0].Keyword != "A_c" || rets[0].Amount != 4 {
		t.Fatal()
	}

	if rets[1].Keyword != "ab" || rets[1].Amount != 1 {
		t.Fatal()
	}

	rets = _m.MatchKey([]string{"efgE一f一二gE一二三FgeF一GE一FGEF一二三四G"})
	if len(rets) != 1 {
		t.Fatal()
	}

	if rets[0].Keyword != "E_F_G" || rets[0].Amount != 5 {
		t.Fatal()
	}

	rets = _m.MatchKey([]string{"hijH1ijxH二二IjxxH三三三I123Jh三三三I333JxxxHiJHIJHI四四四四J"})
	if len(rets) != 1 {
		t.Fatal()
	}

	if rets[0].Keyword != "h_i_j_" || rets[0].Amount != 7 {
		t.Fatal()
	}
}

func TestMatcher_MatchKey(t *testing.T) {
	var rets Results
	_m := newMatcher("a", _matcherItems, nil)

	rets = _m.MatchKey([]string{"abcdef-abcd-abc-ab-a"})
	if len(rets) != 5 {
		t.Fatal(1)
	}

	rets = _m.MatchKey([]string{"babcdefa"})
	if len(rets) != 2 {
		t.Fatal(2)
	}

	rets = _m.MatchKey([]string{"abcba"})
	if len(rets) != 2 {
		t.Fatal(3)
	}

	rets = _m.MatchKey([]string{"babcdefa"})
	if len(rets) != 2 {
		t.Fatal(4)
	}

	_m = newMatcher("a", _matcherItems, &matcherConfig{ignoreCase: true})

	rets = _m.MatchKey([]string{"ABCDEF-ABCD-abc-ab-A"})
	if len(rets) != 5 {
		t.Fatal(1)
	}

	rets = _m.MatchKey([]string{"BaBcDeFa"})
	if len(rets) != 2 {
		t.Fatal(4)
	}
}

func TestMatcher_MatchFirstKey(t *testing.T) {
	var rets Results
	_m := newMatcher("a", _matcherItems, nil)

	rets = _m.MatchFirstText([]string{"abcdef-abcd-abc-ab-a"})
	if rets[0].Keyword != "abcdef" {
		t.Fatal(1)
	}

	rets = _m.MatchFirstText([]string{"ABCDEF-abCd-Abc-ab-a"})
	if rets[0].Keyword != "ab" {
		t.Fatal(12)
	}

	rets = _m.MatchFirstText([]string{"aBcdef-aBcd-abc-ab-a"})
	if rets[0].Keyword != "abc" {
		t.Fatal(13)
	}

	rets = _m.MatchFirstText([]string{"babcdefa"})
	if rets[0].Keyword != "abcdef" {
		t.Fatal(2)
	}

	rets = _m.MatchFirstText([]string{"abcba"})
	if rets[0].Keyword != "abc" {
		t.Fatal(3)
	}

	rets = _m.MatchFirstText([]string{"babcdefa"})
	if rets[0].Keyword != "abcdef" {
		t.Fatal(4)
	}
}

func TestMatcher_MatchLabel(t *testing.T) {
	var rets LabelResults
	_m := newMatcher("a", _matcherItems, nil)

	rets = _m.MatchLabel([]string{"abcdef-abcd-abc-ab-a"})
	if len(rets) != 5 {
		t.Fatal(1)
	}

	rets = _m.MatchLabel([]string{"ABCDEF-ABCD-ABC-AB-A"})
	if len(rets) != 0 {
		t.Fatal(11)
	}

	rets = _m.MatchLabel([]string{"aBcdef-aBcd-aBc-ab-a"})
	if len(rets) != 2 {
		t.Fatal(11)
	}

	rets = _m.MatchLabel([]string{"abcabcabc"})
	if len(rets) != 1 || rets[0].Tags["b"] != "b3" || len(rets[0].Match["abc"]) != 3 || rets[0].Amount != 3 {
		t.Fatal(2)
	}

	rets = _m.MatchLabel([]string{"abccbcbaabccbaabc"})
	if len(rets) != 1 || rets[0].Tags["b"] != "b3" || len(rets[0].Match["abc"]) != 3 || len(rets[0].Match["cba"]) != 2 || len(rets[0].Match["cb"]) != 1 || rets[0].Amount != 6 {
		t.Fatal(31)
	}

	_m = newMatcher("a", _matcherItems, &matcherConfig{ignoreCase: true})

	rets = _m.MatchLabel([]string{"ABCDEF-aBcd-Abc-aB-a"})
	if len(rets) != 5 {
		t.Fatal(1)
	}

	rets = _m.MatchLabel([]string{"aBcABCabC"})
	if len(rets) != 1 || rets[0].Tags["b"] != "b3" || len(rets[0].Match["abc"]) != 3 || rets[0].Amount != 3 {
		t.Fatal(2)
	}

	rets = _m.MatchLabel([]string{"abCCbcBAABCCBAabc"})
	if len(rets) != 1 || rets[0].Tags["b"] != "b3" || len(rets[0].Match["abc"]) != 3 || len(rets[0].Match["cba"]) != 2 || len(rets[0].Match["cb"]) != 1 || rets[0].Amount != 6 {
		t.Fatal(31)
	}

}

func TestMatcher_MatchFirstLabel(t *testing.T) {
	var rets LabelResults
	_m := newMatcher("a", _matcherItems, nil)

	rets = _m.MatchFirstLabel([]string{"abcdef-abcd-abc-ab-a"})
	if len(rets) != 1 || len(rets[0].Match["abcdef"]) <= 0 {
		t.Fatal(1)
	}

	rets = _m.MatchFirstLabel([]string{"abcabcabc"})
	if len(rets) != 1 || len(rets[0].Match["abc"]) <= 0 {
		t.Fatal(2)
	}

	rets = _m.MatchFirstLabel([]string{"cbcbaabcdcbaabc"})
	if len(rets) != 1 || len(rets[0].Match["abcd"]) <= 0 {
		t.Fatal(3)
	}
}

func _outputResults(rts Results) {
	for _, rt := range rts {
		fmt.Println(rt)
	}
}
