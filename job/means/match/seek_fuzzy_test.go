package match

import (
	"strings"
	"testing"
)

func Test_fuzzy(t *testing.T) {
	config := seekConfig{debug: true}

	// lowercase
	t.Run("lowercase", func(t *testing.T) {
		_fuzzy := newFuzzy(1, "a_c", "a_c", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2}, config)
		originSc := seekContent{1, "acabcabcabbcabbbc", "acabcabcabbcabbbc"}
		sr, sc, ok := _fuzzy.seeking(originSc)
		_outputSeekResults(sr, sc)
		_assertSeekResults(t, ok, sc, sr, 4, 3)

		_assertSeekResult(t, sr[0], "a_c", "ac")
		_assertSeekResult(t, sr[1], "a_c", "abc")
		_assertSeekResult(t, sr[2], "a_c", "abc")
		_assertSeekResult(t, sr[3], "a_c", "abbc")

		if strings.Index(sc.origin, "abbbc") < 0 || strings.Index(sc.content, "abbbc") < 0 {
			t.Fatal()
		}
	})

	// uppercase
	t.Run("uppercase", func(t *testing.T) {
		_fuzzy := newFuzzy(1, "A_c", "A_c", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2}, config)
		originSc := seekContent{1, "acAbcabbCabbbcACABC", "acAbcabbCabbbcACABC"}
		sr, sc, ok := _fuzzy.seeking(originSc)
		_outputSeekResults(sr, sc)
		_assertSeekResults(t, ok, sc, sr, 1, 1)

		_assertSeekResult(t, sr[0], "A_c", "Abc")

		if strings.Index(sc.origin, "abbbc") < 0 || strings.Index(sc.content, "abbbc") < 0 {
			t.Fatal()
		}
	})

	// ignore case
	t.Run("ignore case", func(t *testing.T) {
		_fuzzy := newFuzzy(1, "A_c", "a_c", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2}, config)
		originSc := seekContent{1, "acA一ca二二Ca三三三cACcAAcac", "aca一ca二二ca三三三caccaacac"}
		sr, sc, ok := _fuzzy.seeking(originSc)
		_outputSeekResults(sr, sc)
		_assertSeekResults(t, ok, sc, sr, 6, 5)

		_assertSeekResult(t, sr[0], "A_c", "ac")
		_assertSeekResult(t, sr[1], "A_c", "A一c")
		_assertSeekResult(t, sr[2], "A_c", "a二二C")
		_assertSeekResult(t, sr[3], "A_c", "AC")
		_assertSeekResult(t, sr[4], "A_c", "AAc")
		_assertSeekResult(t, sr[5], "A_c", "ac")
	})

	t.Run("underline", func(t *testing.T) {
		_fuzzy := newFuzzy(1, "A_c", "a_c_", map[string]string{"b": "b1"}, FuzzyConfig{Sep: "_", Window: 2}, config)
		originSc := seekContent{1, "aca一ca二二ca三三三cACcac", "aca一ca二二ca三三三cACcac"}
		sr, sc, ok := _fuzzy.seeking(originSc)
		_outputSeekResults(sr, sc)
		_assertSeekResults(t, ok, sc, sr, 4, 3)

		_assertSeekResult(t, sr[0], "A_c", "ac")
		_assertSeekResult(t, sr[1], "A_c", "a一c")
		_assertSeekResult(t, sr[2], "A_c", "a二二c")
		_assertSeekResult(t, sr[3], "A_c", "ac")
	})
}
