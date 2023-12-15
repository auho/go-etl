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
	var amount, textAmount int
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

	rets = _m.MatchKey([]string{"ABCDABcAbabacabABBaAc_aE_F_G_e_f_g_h_i_j_H_I_J_iefgAxxciaB"})
	_outputResults(rets)
	if len(rets) != 6 {
		t.Fatal()
	}

	amount = 0
	textAmount = 0
	for _, ret := range rets {
		amount += ret.Amount
		for _, _n := range ret.TextsAmount {
			textAmount += _n
		}
	}
	if amount != textAmount {
		t.Fatal()
	}

	if rets[0].Keyword != "ca" || rets[0].Amount != 1 || rets[0].TextsAmount["ca"] != 1 || rets[0].Tags["b"] != "b4" {
		t.Fatal()
	}

	amount = 0
	for _, _n := range rets[0].TextsAmount {
		amount += _n
	}
	if amount != rets[0].Amount {
		t.Fatal()
	}

	if rets[1].Keyword != "ab" || rets[1].Amount != 1 {
		t.Fatal()
	}

	if rets[2].Keyword != "a" || rets[2].Amount != 4 || rets[2].TextsAmount["a"] != 4 || rets[2].Tags["b"] != "b6" {
		t.Fatal()
	}

	if rets[3].Keyword != "A_c" || rets[3].Amount != 3 || rets[3].TextsAmount["ABc"] != 1 || rets[3].TextsAmount["Ac"] != 1 || rets[3].TextsAmount["Axxc"] != 1 {
		t.Fatal()
	}

	amount = 0
	for _, _n := range rets[3].TextsAmount {
		amount += _n
	}
	if amount != rets[3].Amount {
		t.Fatal()
	}

	if rets[4].Keyword != "E_F_G" || rets[4].Amount != 1 {
		t.Fatal()
	}

	if rets[5].Keyword != "h_i_j_" || rets[5].Amount != 1 {
		t.Fatal()
	}
}

func TestMatcher_MatchKey_Accurate_IgnoreCase(t *testing.T) {
	var amount, textAmount int
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

	rets = _m.MatchKey([]string{"ABCDABcAbabacabABBaAc_aE_F_G_e_f_g_h_i_j_H_I_J_iefgAxxciaB"})
	if len(rets) != 5 {
		t.Fatal()
	}

	amount = 0
	textAmount = 0
	for _, ret := range rets {
		amount += ret.Amount
		for _, _n := range ret.TextsAmount {
			textAmount += _n
		}
	}
	if amount != textAmount {
		t.Fatal()
	}

	if rets[0].Keyword != "A_c" || rets[0].Amount != 5 || rets[0].Tags["b"] != "b7" {
		t.Fatal()
	}

	if rets[0].TextsAmount["ABC"] != 1 || rets[0].TextsAmount["ABc"] != 1 || rets[0].TextsAmount["abac"] != 1 || rets[0].TextsAmount["aAc"] != 1 || rets[0].TextsAmount["Axxc"] != 1 {
		t.Fatal()
	}

	amount = 0
	for _, _n := range rets[0].TextsAmount {
		amount += _n
	}
	if amount != rets[0].Amount {
		t.Fatal()
	}

	if rets[1].Keyword != "E_F_G" || rets[1].Amount != 3 {
		t.Fatal()
	}

	if rets[2].Keyword != "h_i_j_" || rets[2].Amount != 2 || rets[2].Tags["b"] != "b9" || rets[2].TextsAmount["h_i_j"] != 1 || rets[2].TextsAmount["H_I_J"] != 1 {
		t.Fatal()
	}

	if rets[3].Keyword != "ab" || rets[3].Amount != 4 || rets[3].TextsAmount["Ab"] != 1 || rets[3].TextsAmount["ab"] != 1 || rets[3].TextsAmount["AB"] != 1 || rets[3].TextsAmount["aB"] != 1 {
		t.Fatal()
	}

	amount = 0
	for _, _n := range rets[3].TextsAmount {
		amount += _n
	}
	if amount != rets[3].Amount {
		t.Fatal()
	}

	if rets[4].Keyword != "a" || rets[4].Amount != 1 {
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

	amount := 0
	textAmount := 0
	for _, ret := range rets {
		amount += ret.Amount
		for _, _n := range ret.TextsAmount {
			textAmount += _n
		}
	}
	if amount != textAmount {
		t.Fatal()
	}

	if rets[0].Keyword != "A_c" || rets[0].Amount != 4 || rets[0].Tags["b"] != "b7" {
		t.Fatal()
	}

	if rets[0].TextsAmount["AbbbC"] != 1 || rets[0].TextsAmount["Abc"] != 1 || rets[0].TextsAmount["abbC"] != 1 || rets[0].TextsAmount["ac"] != 1 {
		t.Fatal()
	}

	amount = 0
	for _, _n := range rets[0].TextsAmount {
		amount += _n
	}
	if amount != rets[0].Amount {
		t.Fatal()
	}

	if rets[1].Keyword != "ab" || rets[1].Amount != 1 {
		t.Fatal()
	}

	if rets[1].TextsAmount["AB"] != 1 {
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

	if rets[0].TextsAmount["H1ij"] != 1 || rets[0].TextsAmount["HIJ"] != 1 || rets[0].TextsAmount["HiJ"] != 1 || rets[0].TextsAmount["H三三三I123J"] != 1 ||
		rets[0].TextsAmount["H二二Ij"] != 1 || rets[0].TextsAmount["hij"] != 1 || rets[0].TextsAmount["h三三三I333J"] != 1 {
		t.Fatal()
	}
}

func TestMatcher_MatchKey(t *testing.T) {
	var rets Results
	_m := newMatcher("a", _matcherItems, nil)

	rets = _m.MatchKey([]string{"abcdef-abcd-abc-ab-a"})
	if len(rets) != 5 {
		t.Fatal()
	}

	rets = _m.MatchKey([]string{"babcdefa"})
	if len(rets) != 2 {
		t.Fatal()
	}

	rets = _m.MatchKey([]string{"abcba"})
	if len(rets) != 2 {
		t.Fatal()
	}

	rets = _m.MatchKey([]string{"babcdefa"})
	if len(rets) != 2 {
		t.Fatal()
	}

	_m = newMatcher("a", _matcherItems, &matcherConfig{ignoreCase: true})

	rets = _m.MatchKey([]string{"ABCDEF-ABCD-abc-ab-A"})
	if len(rets) != 5 {
		t.Fatal()
	}

	rets = _m.MatchKey([]string{"BaBcDeFa"})
	if len(rets) != 2 {
		t.Fatal()
	}
}

func TestMatcher_MatchFirstKey(t *testing.T) {
	var rets Results
	_m := newMatcher("a", _matcherItems, &matcherConfig{
		debug: true,
	})

	rets = _m.MatchFirstKey([]string{"abcdef-abcd-abc-ab-a"})
	if rets[0].Keyword != "abcdef" {
		t.Fatal()
	}

	rets = _m.MatchFirstKey([]string{"ABCDEF-abCd-Abc-ab-a"})
	if rets[0].Keyword != "ab" {
		t.Fatal()
	}

	rets = _m.MatchFirstKey([]string{"aBcdef-aBcd-abc-ab-a"})
	if rets[0].Keyword != "abc" {
		t.Fatal()
	}

	rets = _m.MatchFirstKey([]string{"babcdefa"})
	if rets[0].Keyword != "abcdef" {
		t.Fatal()
	}

	rets = _m.MatchFirstKey([]string{"abcba"})
	if rets[0].Keyword != "abc" {
		t.Fatal()
	}

	rets = _m.MatchFirstKey([]string{"caabcdef"})
	if rets[0].Keyword != "abcdef" {
		t.Fatal()
	}
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

func TestMatcher_MatchLabel(t *testing.T) {
	var amount int
	var rets LabelResults
	_m := newMatcher("a", _matcherItems, &matcherConfig{
		debug: true,
	})

	rets = _m.MatchLabel([]string{"abcdef-abcd-abc-ab-a"})
	if len(rets) != 5 {
		t.Fatal()
	}

	rets = _m.MatchLabel([]string{"ABCDEF-ABCD-ABC-AB-A"})
	if len(rets) != 0 {
		t.Fatal()
	}

	rets = _m.MatchLabel([]string{"aBcdef-aBcd-aBc-ab-a"})
	if len(rets) != 2 {
		t.Fatal()
	}

	rets = _m.MatchLabel([]string{"abcabcabc"})
	if len(rets) != 1 || rets[0].Tags["b"] != "b3" || rets[0].Match["abc"]["abc"] != 3 || rets[0].Amount != 3 {
		t.Fatal()
	}

	rets = _m.MatchLabel([]string{"abccbcbaabccbaabc"})
	if len(rets) != 1 {
		t.Fatal()
	}

	amount = 0
	textAmount := 0
	for _, ret := range rets {
		amount += ret.Amount
		for _, __m := range ret.Match {
			for _, _n := range __m {
				textAmount += _n
			}
		}
	}

	if amount != textAmount {
		t.Fatal()
	}

	if rets[0].Tags["b"] != "b3" || rets[0].Match["abc"]["abc"] != 3 || rets[0].Match["cba"]["cba"] != 2 || rets[0].Match["cb"]["cb"] != 1 || rets[0].Amount != 6 {
		t.Fatal()
	}

	_m = newMatcher("a", _matcherItems, &matcherConfig{ignoreCase: true})

	rets = _m.MatchLabel([]string{"ABCDEF-aBcd-Abc-aB-a"})
	if len(rets) != 5 {
		t.Fatal()
	}

	rets = _m.MatchLabel([]string{"aBcABCabC"})
	if len(rets) != 1 || rets[0].Tags["b"] != "b3" || len(rets[0].Match["abc"]) != 3 || rets[0].Amount != 3 {
		t.Fatal()
	}

	rets = _m.MatchLabel([]string{"abCCbcBAABCCBAabc"})
	if len(rets) != 1 || rets[0].Tags["b"] != "b3" || len(rets[0].Match["abc"]) != 3 || len(rets[0].Match["cba"]) != 2 || len(rets[0].Match["cb"]) != 1 || rets[0].Amount != 6 {
		t.Fatal()
	}

}

func TestMatcher_MatchFirstLabel(t *testing.T) {
	var rets LabelResults
	_m := newMatcher("a", _matcherItems, &matcherConfig{debug: true})

	rets = _m.MatchFirstLabel([]string{"abcdef-abcd-abc-ab-a"})
	if len(rets) != 1 || len(rets[0].Match["abcdef"]) <= 0 {
		t.Fatal()
	}

	rets = _m.MatchFirstLabel([]string{"acbabcabcabc"})
	if len(rets) != 1 || len(rets[0].Match["abc"]) <= 0 {
		t.Fatal()
	}

	rets = _m.MatchFirstLabel([]string{"cbcbaabcdcbaabc"})
	_outputLabelResults(rets)
	if len(rets) != 1 || len(rets[0].Match["abcd"]) <= 0 {
		t.Fatal()
	}
}

func _outputResults(rts Results) {
	for _, rt := range rts {
		fmt.Println(fmt.Sprintf("%+v", rt))
	}
}

func _outputLabelResults(rts LabelResults) {
	for _, rt := range rts {
		fmt.Println(fmt.Sprintf("%+v", rt))
	}
}
