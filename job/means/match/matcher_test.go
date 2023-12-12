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
}

func TestMatcher_MatchKey_Accurate(t *testing.T) {
	var _rts Results
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

	_rts = _m.MatchKey([]string{"ABCDABcAbabacabABBaAc_aE_F_G_e_f_g_h_i_j_H_I_J_efgAxxc"})
	if len(_rts) != 5 {
		t.Fatal()
	}

	if _rts[0].Keyword != "ca" || _rts[0].Amount != 1 {
		t.Fatal()
	}

	if _rts[1].Keyword != "ab" || _rts[1].Amount != 1 {
		t.Fatal()
	}

	if _rts[2].Keyword != "a" || _rts[2].Amount != 3 {
		t.Fatal()
	}

	if _rts[3].Keyword != "E_F_G" || _rts[3].Amount != 1 {
		t.Fatal()
	}

	if _rts[4].Keyword != "h_i_j_" || _rts[4].Amount != 1 {
		t.Fatal()
	}
}

func TestMatcher_MatchKey_Fuzzy(t *testing.T) {
	var _rts Results
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

	_rts = _m.MatchKey([]string{"acAbcabbCAbbbCABbbBc"})
	if len(_rts) != 2 {
		t.Fatal()
	}

	if _rts[0].Keyword != "A_c" || _rts[0].Amount != 4 {
		t.Fatal()
	}

	if _rts[1].Keyword != "ab" || _rts[1].Amount != 1 {
		t.Fatal()
	}

	_rts = _m.MatchKey([]string{"acAbcabbCAbbbCABbbBc"})
	if len(_rts) != 2 {
		t.Fatal()
	}

	if _rts[0].Keyword != "A_c" || _rts[0].Amount != 4 {
		t.Fatal()
	}

	if _rts[1].Keyword != "ab" || _rts[1].Amount != 1 {
		t.Fatal()
	}

	_rts = _m.MatchKey([]string{"efgE一f一二gE一二三FgeF一GE一FGEF一二三四G"})
	if len(_rts) != 1 {
		t.Fatal()
	}

	if _rts[0].Keyword != "E_F_G" || _rts[0].Amount != 5 {
		t.Fatal()
	}

	_rts = _m.MatchKey([]string{"hijH1ijxH二二IjxxH三三三I123Jh三三三I333JxxxHiJHIJHI四四四四J"})
	if len(_rts) != 1 {
		t.Fatal()
	}

	if _rts[0].Keyword != "h_i_j_" || _rts[0].Amount != 7 {
		t.Fatal()
	}
}

func TestMatcher_MatchKey(t *testing.T) {
	var _rts Results
	_m := newMatcher("a", _matcherItems, nil)

	_rts = _m.MatchKey([]string{"abcdef-abcd-abc-ab-a"})
	if len(_rts) != 5 {
		t.Fatal(1)
	}

	_rts = _m.MatchKey([]string{"babcdefa"})
	if len(_rts) != 2 {
		t.Fatal(2)
	}

	_rts = _m.MatchKey([]string{"abcba"})
	if len(_rts) != 2 {
		t.Fatal(3)
	}

	_rts = _m.MatchKey([]string{"babcdefa"})
	if len(_rts) != 2 {
		t.Fatal(4)
	}

	_m = newMatcher("a", _matcherItems, &matcherConfig{ignoreCase: true})

	_rts = _m.MatchKey([]string{"ABCDEF-ABCD-abc-ab-A"})
	if len(_rts) != 5 {
		t.Fatal(1)
	}

	_rts = _m.MatchKey([]string{"BaBcDeFa"})
	if len(_rts) != 2 {
		t.Fatal(4)
	}
}

func TestMatcher_MatchFirstKey(t *testing.T) {
	var _rts Results
	_m := newMatcher("a", _matcherItems, nil)

	_rts = _m.MatchFirstText([]string{"abcdef-abcd-abc-ab-a"})
	if _rts[0].Keyword != "abcdef" {
		t.Fatal(1)
	}

	_rts = _m.MatchFirstText([]string{"ABCDEF-abCd-Abc-ab-a"})
	if _rts[0].Keyword != "ab" {
		t.Fatal(12)
	}

	_rts = _m.MatchFirstText([]string{"aBcdef-aBcd-abc-ab-a"})
	if _rts[0].Keyword != "abc" {
		t.Fatal(13)
	}

	_rts = _m.MatchFirstText([]string{"babcdefa"})
	if _rts[0].Keyword != "abcdef" {
		t.Fatal(2)
	}

	_rts = _m.MatchFirstText([]string{"abcba"})
	if _rts[0].Keyword != "abc" {
		t.Fatal(3)
	}

	_rts = _m.MatchFirstText([]string{"babcdefa"})
	if _rts[0].Keyword != "abcdef" {
		t.Fatal(4)
	}
}

func TestMatcher_MatchLabel(t *testing.T) {
	var _rts LabelResults
	_m := newMatcher("a", _matcherItems, nil)

	_rts = _m.MatchLabel([]string{"abcdef-abcd-abc-ab-a"})
	if len(_rts) != 5 {
		t.Fatal(1)
	}

	_rts = _m.MatchLabel([]string{"ABCDEF-ABCD-ABC-AB-A"})
	if len(_rts) != 0 {
		t.Fatal(11)
	}

	_rts = _m.MatchLabel([]string{"aBcdef-aBcd-aBc-ab-a"})
	if len(_rts) != 2 {
		t.Fatal(11)
	}

	_rts = _m.MatchLabel([]string{"abcabcabc"})
	if len(_rts) != 1 || _rts[0].Tags["b"] != "b3" || len(_rts[0].Match["abc"]) != 3 || _rts[0].Amount != 3 {
		t.Fatal(2)
	}

	_rts = _m.MatchLabel([]string{"abccbcbaabccbaabc"})
	if len(_rts) != 1 || _rts[0].Tags["b"] != "b3" || len(_rts[0].Match["abc"]) != 3 || len(_rts[0].Match["cba"]) != 2 || len(_rts[0].Match["cb"]) != 1 || _rts[0].Amount != 6 {
		t.Fatal(31)
	}

	_m = newMatcher("a", _matcherItems, &matcherConfig{ignoreCase: true})

	_rts = _m.MatchLabel([]string{"ABCDEF-aBcd-Abc-aB-a"})
	if len(_rts) != 5 {
		t.Fatal(1)
	}

	_rts = _m.MatchLabel([]string{"aBcABCabC"})
	if len(_rts) != 1 || _rts[0].Tags["b"] != "b3" || len(_rts[0].Match["abc"]) != 3 || _rts[0].Amount != 3 {
		t.Fatal(2)
	}

	_rts = _m.MatchLabel([]string{"abCCbcBAABCCBAabc"})
	if len(_rts) != 1 || _rts[0].Tags["b"] != "b3" || len(_rts[0].Match["abc"]) != 3 || len(_rts[0].Match["cba"]) != 2 || len(_rts[0].Match["cb"]) != 1 || _rts[0].Amount != 6 {
		t.Fatal(31)
	}

}

func TestMatcher_MatchFirstLabel(t *testing.T) {
	var _rts LabelResults
	_m := newMatcher("a", _matcherItems, nil)

	_rts = _m.MatchFirstLabel([]string{"abcdef-abcd-abc-ab-a"})
	if len(_rts) != 1 || len(_rts[0].Match["abcdef"]) <= 0 {
		t.Fatal(1)
	}

	_rts = _m.MatchFirstLabel([]string{"abcabcabc"})
	if len(_rts) != 1 || len(_rts[0].Match["abc"]) <= 0 {
		t.Fatal(2)
	}

	_rts = _m.MatchFirstLabel([]string{"cbcbaabcdcbaabc"})
	if len(_rts) != 1 || len(_rts[0].Match["abcd"]) <= 0 {
		t.Fatal(3)
	}
}

func _outputResults(rts Results) {
	for _, rt := range rts {
		fmt.Println(rt)
	}
}
