package match

import (
	"testing"
)

func Test_accurate(t *testing.T) {
	config := seekConfig{debug: true}

	// lowercase
	t.Run("lowercase", func(t *testing.T) {
		_accurate := newAccurate(1, "ac", "ac", map[string]string{"b": "b1"}, config)
		originSc := seekContent{1, "acabcabbcabbbcac", "acabcabbcabbbcac"}
		sr, sc, ok := _accurate.seeking(originSc)
		_outputSeekResults(sr, sc)
		_assertSeekResults(t, ok, sc, sr, 2, 1)

		_assertSeekResult(t, sr[0], "ac", "ac")
		_assertSeekResult(t, sr[1], "ac", "ac")
	})

	// uppercase
	t.Run("uppercase", func(t *testing.T) {
		_accurate := newAccurate(1, "Ac", "AC", map[string]string{"b": "b1"}, config)
		originSc := seekContent{1, "AcaCACaccabbCabbbcAC", "AcaCACaccabbCabbbcAC"}
		sr, sc, ok := _accurate.seeking(originSc)
		_outputSeekResults(sr, sc)
		_assertSeekResults(t, ok, sc, sr, 2, 1)

		_assertSeekResult(t, sr[0], "Ac", "AC")
		_assertSeekResult(t, sr[1], "Ac", "AC")
	})

	// ignore case
	t.Run("ignore case", func(t *testing.T) {
		_accurate := newAccurate(1, "Ac", "ac", map[string]string{"b": "b1"}, config)
		originSc := seekContent{1, "AcaCACaccabbCabbbcAC", "acacacaccabbcabbbcac"}
		sr, sc, ok := _accurate.seeking(originSc)
		_outputSeekResults(sr, sc)
		_assertSeekResults(t, ok, sc, sr, 5, 4)

		_assertSeekResult(t, sr[0], "Ac", "Ac")
		_assertSeekResult(t, sr[1], "Ac", "aC")
		_assertSeekResult(t, sr[2], "Ac", "AC")
		_assertSeekResult(t, sr[3], "Ac", "ac")
		_assertSeekResult(t, sr[4], "Ac", "AC")
	})
}
