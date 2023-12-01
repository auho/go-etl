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
}

var _corpus = []string{
	"abcdef-abcd-abc-ab-a",
	"a-ab-abc-abcd-abcdef",
	"a-abcdef",
	"abcba",
	"babcdefa",
}

func TestMatcher(t *testing.T) {
	_m := newMatcher("a", _matcherItems)
	_m.MatchKey(_corpus)
	_m.MatchFirstKey(_corpus)
	_m.MatchLabel(_corpus)
	_m.MatchFirstLabel(_corpus)
}

func TestMatcher_MatchKey(t *testing.T) {
	var _rts Results
	_m := newMatcher("a", _matcherItems)

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
}

func TestMatcher_MatchFirstKey(t *testing.T) {
	var _rts Results
	_m := newMatcher("a", _matcherItems)

	_rts = _m.MatchFirstKey([]string{"abcdef-abcd-abc-ab-a"})
	if _rts[0].Key != "abcdef" {
		t.Fatal(1)
	}

	_rts = _m.MatchFirstKey([]string{"babcdefa"})
	if _rts[0].Key != "abcdef" {
		t.Fatal(2)
	}

	_rts = _m.MatchFirstKey([]string{"abcba"})
	if _rts[0].Key != "abcd" {
		t.Fatal(3)
	}

	_rts = _m.MatchFirstKey([]string{"babcdefa"})
	if _rts[0].Key != "abcdef" {
		t.Fatal(4)
	}
}

func TestMatcher_MatchLabel(t *testing.T) {
	var _rts LabelResults
	_m := newMatcher("a", _matcherItems)

	_rts = _m.MatchLabel([]string{"abcdef-abcd-abc-ab-a"})
	if len(_rts) != 5 {
		t.Fatal(1)
	}

	_rts = _m.MatchLabel([]string{"abcabcabc"})
	if len(_rts) != 1 || _rts[0].Labels["b"] != "b3" || _rts[0].Match["abc"] != 3 || _rts[0].MatchAmount != 3 {
		t.Fatal(2)
	}

	_rts = _m.MatchLabel([]string{"abccbcbaabccbaabc"})
	if len(_rts) != 1 || _rts[0].Labels["b"] != "b3" || _rts[0].Match["abc"] != 3 || _rts[0].Match["cba"] != 2 || _rts[0].Match["cb"] != 1 || _rts[0].MatchAmount != 6 {
		t.Fatal(31)
	}
}

func TestMatcher_MatchFirstLabel(t *testing.T) {
	var _rts LabelResults
	_m := newMatcher("a", _matcherItems)

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
