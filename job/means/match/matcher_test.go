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
	_m.MatchFirstKey(_corpus)
	_m.MatchLabel(_corpus)
	_m.MatchFirstLabel(_corpus)
}

func TestMatcher_MatchKey_Fuzzy(t *testing.T) {
	var _rts Results
	_m := newMatcher("a", _matcherItems, &matcherConfig{
		ignoreCase:  true,
		mode:        modePriorityFuzzy,
		enableFuzzy: true,
		fuzzyConfig: FuzzyConfig{
			Window: 3,
			Sep:    "_",
		},
	})

	_rts = _m.MatchKey([]string{"acAbcabbCAbbbCABbbBc"})
	if _rts.Len() != 2 {
		t.Fatal()
	}

	if _rts[0].Key != "A_c" {
		t.Fatal()
	}

	if _rts[0].Num != 4 {
		t.Fatal()
	}

	if _rts[1].Key != "ab" {
		t.Fatal()
	}

	if _rts[1].Num != 1 {
		t.Fatal()
	}

	_rts = _m.MatchKey([]string{"efgE一f一二gE一二三FgeF一GE一FGEF一二三四G"})
	if _rts.Len() != 1 {
		t.Fatal()
	}

	if _rts[0].Key != "E_F_G" {
		t.Fatal()
	}

	_rts = _m.MatchKey([]string{"efgE一f一二gE一二三FgeF一GE一FGEF一二三四G"})
	if _rts.Len() != 1 {
		t.Fatal()
	}

	if _rts[0].Key != "E_F_G" {
		t.Fatal()
	}

	if _rts[0].Num != 5 {
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

	_rts = _m.MatchFirstKey([]string{"abcdef-abcd-abc-ab-a"})
	if _rts[0].Key != "abcdef" {
		t.Fatal(1)
	}

	_rts = _m.MatchFirstKey([]string{"ABCDEF-abCd-Abc-ab-a"})
	if _rts[0].Key != "ab" {
		t.Fatal(12)
	}

	_rts = _m.MatchFirstKey([]string{"aBcdef-aBcd-abc-ab-a"})
	if _rts[0].Key != "abc" {
		t.Fatal(13)
	}

	_rts = _m.MatchFirstKey([]string{"babcdefa"})
	if _rts[0].Key != "abcdef" {
		t.Fatal(2)
	}

	_rts = _m.MatchFirstKey([]string{"abcba"})
	if _rts[0].Key != "abc" {
		t.Fatal(3)
	}

	_rts = _m.MatchFirstKey([]string{"babcdefa"})
	if _rts[0].Key != "abcdef" {
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
	if len(_rts) != 1 || _rts[0].Labels["b"] != "b3" || _rts[0].Match["abc"] != 3 || _rts[0].MatchAmount != 3 {
		t.Fatal(2)
	}

	_rts = _m.MatchLabel([]string{"abccbcbaabccbaabc"})
	if len(_rts) != 1 || _rts[0].Labels["b"] != "b3" || _rts[0].Match["abc"] != 3 || _rts[0].Match["cba"] != 2 || _rts[0].Match["cb"] != 1 || _rts[0].MatchAmount != 6 {
		t.Fatal(31)
	}

	_m = newMatcher("a", _matcherItems, &matcherConfig{ignoreCase: true})

	_rts = _m.MatchLabel([]string{"ABCDEF-aBcd-Abc-aB-a"})
	if len(_rts) != 5 {
		t.Fatal(1)
	}

	_rts = _m.MatchLabel([]string{"aBcABCabC"})
	if len(_rts) != 1 || _rts[0].Labels["b"] != "b3" || _rts[0].Match["abc"] != 3 || _rts[0].MatchAmount != 3 {
		t.Fatal(2)
	}

	_rts = _m.MatchLabel([]string{"abCCbcBAABCCBAabc"})
	if len(_rts) != 1 || _rts[0].Labels["b"] != "b3" || _rts[0].Match["abc"] != 3 || _rts[0].Match["cba"] != 2 || _rts[0].Match["cb"] != 1 || _rts[0].MatchAmount != 6 {
		t.Fatal(31)
	}

}

func TestMatcher_MatchFirstLabel(t *testing.T) {
	var _rts LabelResults
	_m := newMatcher("a", _matcherItems, nil)

	_rts = _m.MatchFirstLabel([]string{"abcdef-abcd-abc-ab-a"})
	if len(_rts) != 1 || _rts[0].Match["abcdef"] <= 0 {
		t.Fatal(1)
	}

	_rts = _m.MatchFirstLabel([]string{"abcabcabc"})
	if len(_rts) != 1 || _rts[0].Match["abc"] <= 0 {
		t.Fatal(2)
	}

	_rts = _m.MatchFirstLabel([]string{"cbcbaabcdcbaabc"})
	if len(_rts) != 1 || _rts[0].Match["abcd"] <= 0 {
		t.Fatal(3)
	}
}

func _outputResults(rts Results) {
	for _, rt := range rts {
		fmt.Println(rt)
	}
}
