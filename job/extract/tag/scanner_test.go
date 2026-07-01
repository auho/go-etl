package tag

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var _expectItemsAmount = 7
var _expectMatcherAmount = 31
var _expectTextAmount = 17
var _expectKeyAmount = 4
var _expectLabelAmount = 3

var _expect_b_Amount = 3
var _expect_123_Amount = 7
var _expect_中文_Amount = 4
var _expect_中_文_Amount = 17
var _expect_中_文_Num = 14

func _genMatcher() *scanner {
	// keyword: a
	// labels: b c
	items := []map[string]string{
		{"a": "123", "b": "b", "c": "c"},
		{"a": "b", "b": "b", "c": "c"},
		{"a": "中文", "b": "b2", "c": "c2"},
		{"a": "中_文", "b": "b3", "c": "c3"},
		{"a": "_中1_a文", "b": "b4", "c": "c4"},
		{"a": "，。【】", "b": "b5", "c": "c5"},
		{"a": ".+*?()|[]{}^$`))", "b": "b6", "c": "c6"},
	}

	_scanner := newScanner("a", items, withScannerKeyFormatFunc(func(s string) string {
		res, err := regexp.MatchString(`^[\w+._\s()]+$`, s)
		if err != nil {
			return s
		}

		if res {
			return fmt.Sprintf(`\b%s\b`, s)
		} else {
			return strings.ReplaceAll(s, "_", `.{1,3}`)
		}
	}))

	return _scanner
}

func TestMatcher(t *testing.T) {
	_scanner := _genMatcher()
	if len(_scanner.regexpItems) != _expectItemsAmount {
		t.Fatal()
	}

	var rets results
	var labelRets labelResults

	t.Run("Scan", func(t *testing.T) {
		rets = _scanner.Scan(_contents)
		_outputResults(rets)

		_assertResults(t, rets, _expectMatcherAmount, _expectMatcherAmount)

		_assertResult(t, rets[0], "b", 1, 1, map[string]int{"b": 1})
		_assertResult(t, rets[1], "123", 1, 1, map[string]int{"123": 1})
		_assertResult(t, rets[7], "中文", 1, 1, map[string]int{"中文": 1})
		_assertResult(t, rets[30], "中_文", 1, 1, map[string]int{"中123文": 1})
	})

	t.Run("ScanInKeyOrder", func(t *testing.T) {
		rets = _scanner.ScanInKeyOrder(_contents)
		_outputResults(rets)

		_assertResults(t, rets, _expectMatcherAmount, _expectMatcherAmount)

		_assertResult(t, rets[0], "123", 1, 1, map[string]int{"123": 1})
		_assertResult(t, rets[_expect_123_Amount], "b", 1, 1, map[string]int{"b": 1})
		_assertResult(t, rets[_expect_b_Amount+_expect_123_Amount], "中文", 1, 1, map[string]int{"中文": 1})
		_assertResult(t, rets[_expect_b_Amount+_expect_123_Amount+_expect_中文_Amount], "中_文", 1, 1, map[string]int{"中bb文": 1})
		_assertResult(t, rets[_expectMatcherAmount-1], "中_文", 1, 1, map[string]int{"中123文": 1})
	})

	t.Run("ScanText", func(t *testing.T) {
		rets = _scanner.ScanText(_contents)
		_outputResults(rets)

		_assertResults(t, rets, _expectMatcherAmount, _expectTextAmount)

		_assertResult(t, rets[0], "b", _expect_b_Amount, 1, map[string]int{"b": _expect_b_Amount})
		_assertResult(t, rets[1], "123", _expect_123_Amount, 1, map[string]int{"123": _expect_123_Amount})
		_assertResult(t, rets[2], "中文", _expect_中文_Amount, 1, map[string]int{"中文": _expect_中文_Amount})
		_assertResult(t, rets[3], "中_文", 1, 1, map[string]int{"中bb文": 1})
		_assertResult(t, rets[6], "中_文", 2, 1, map[string]int{"中aa文": 2})
		_assertResult(t, rets[16], "中_文", 1, 1, map[string]int{"中23文": 1})

	})

	t.Run("ScanFirstText", func(t *testing.T) {
		rets = _scanner.ScanFirstText(_contents)
		_outputResults(rets)

		_assertResults(t, rets, 1, 1)

		_assertResult(t, rets[0], "b", 1, 1, map[string]int{"b": 1})
	})

	t.Run("ScanLastText", func(t *testing.T) {
		rets = _scanner.ScanLastText(_contents)
		_outputResults(rets)

		_assertResults(t, rets, 1, 1)

		_assertResult(t, rets[0], "中_文", 1, 1, map[string]int{"中123文": 1})
	})

	t.Run("ScanMostText", func(t *testing.T) {
		rets = _scanner.ScanMostText(_contents)
		_outputResults(rets)

		_assertResults(t, rets, _expect_123_Amount, 1)

		_assertResult(t, rets[0], "123", _expect_123_Amount, 1, map[string]int{"123": _expect_123_Amount})
	})

	t.Run("ScanKey", func(t *testing.T) {
		rets = _scanner.ScanKey(_contents)
		_outputResults(rets)

		_assertResults(t, rets, _expectMatcherAmount, _expectKeyAmount)

		_assertResult(t, rets[0], "123", _expect_123_Amount, 1, map[string]int{"123": _expect_123_Amount})
		_assertResult(t, rets[1], "b", _expect_b_Amount, 1, map[string]int{"b": _expect_b_Amount})
		_assertResult(t, rets[2], "中文", _expect_中文_Amount, 1, map[string]int{"中文": _expect_中文_Amount})
		_assertResult(t, rets[3], "中_文", _expect_中_文_Amount, _expect_中_文_Num, map[string]int{
			"中00文": 2,
			"中aa文": 2,
			"中12文": 1,
			"中ab文": 1,
			"中二二文": 1,

			"b": 0,
		})
	})

	t.Run("ScanFirstKey", func(t *testing.T) {
		rets = _scanner.ScanFirstKey(_contents)
		_outputResults(rets)

		_assertResults(t, rets, 1, 1)
		_assertResult(t, rets[0], "123", 1, 1, map[string]int{"123": 1})
	})

	t.Run("ScanLastKey", func(t *testing.T) {
		rets = _scanner.ScanLastKey(_contents)
		_outputResults(rets)

		_assertResults(t, rets, 1, 1)

		_assertResult(t, rets[0], "中_文", 1, 1, map[string]int{"中123文": 1})
	})

	t.Run("ScanMostKey", func(t *testing.T) {
		rets = _scanner.ScanMostKey(_contents)
		_outputResults(rets)

		_assertResult(t, rets[0], "中_文", _expect_中_文_Amount, _expect_中_文_Num, map[string]int{
			"中00文": 2,
			"中aa文": 2,
			"中12文": 1,
			"中ab文": 1,
			"中二二文": 1,

			"b": 0,
		})
	})

	t.Run("ScanLabel", func(t *testing.T) {
		labelRets = _scanner.ScanLabel(_contents)
		_outputResults(labelRets)

		_assertLabelResults(t, labelRets, _expectMatcherAmount, _expectLabelAmount)

		_assertLabelResult(t, labelRets[0], "-b-c", _expect_123_Amount+_expect_b_Amount, 2, 2, map[string]int{
			"123": _expect_123_Amount,
			"b":   _expect_b_Amount,

			"中文": 0,
		})

		_assertLabelResult(t, labelRets[1], "-b2-c2", _expect_中文_Amount, 1, 1, map[string]int{
			"中文": _expect_中文_Amount,

			"b": 0,
		})

		_assertLabelResult(t, labelRets[2], "-b3-c3", _expect_中_文_Amount, 1, _expect_中_文_Num, map[string]int{
			"中00文": 2,
			"中aa文": 2,
			"中12文": 1,
			"中ab文": 1,
			"中二二文": 1,

			"b": 0,
		})
	})

	t.Run("ScanLabelMostText", func(t *testing.T) {
		labelRets = _scanner.ScanLabelMostText(_contents)
		_outputResults(labelRets)

		_assertLabelResults(t, labelRets, _expect_中_文_Amount, 1)
		_assertLabelResult(t, labelRets[0], "-b3-c3", _expect_中_文_Amount, 1, _expect_中_文_Num, map[string]int{
			"中00文": 2,
			"中aa文": 2,
			"中12文": 1,
			"中ab文": 1,
			"中二二文": 1,

			"b": 0,
		})
	})
}

func _assertResult(t *testing.T, ret result, keyword string, amount, textsNum int, textsAmount map[string]int) {
	if ret.keyword != keyword {
		t.Fatal(fmt.Sprintf("result[%s != %s]", keyword, ret.keyword), t.Name())
	}

	if ret.amount != amount {
		t.Fatal(fmt.Sprintf("result[%s] amount", keyword), t.Name())
	}

	if textsNum != len(ret.texts) {
		t.Fatal(fmt.Sprintf("result texts[%s] num[%d != %d]", keyword, textsNum, len(ret.texts)), t.Name())
	}

	for _t, _a := range textsAmount {
		if ret.texts[_t] != _a {
			t.Fatal(fmt.Sprintf("result text[%s:%s] amont[%d != %d]", keyword, _t, _a, ret.texts[_t]), t.Name())
		}
	}
}

func _assertResults(t *testing.T, rets results, allTextsAmount int, resultsAmount int) {
	if resultsAmount != len(rets) {
		t.Fatal("results len", t.Name())
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

	if amount != allTextsAmount {
		t.Fatal(fmt.Sprintf("results amount[%d!= %d]", allTextsAmount, amount), t.Name())
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

func _assertLabelResult(t *testing.T, ret labelResult, id string, amount, keysNum, textsNum int, textsAmount map[string]int) {
	if ret.identity != id {
		t.Fatal(fmt.Sprintf("result[%s != %s]", id, ret.identity), t.Name())
	}

	if ret.amount != amount {
		t.Fatal(fmt.Sprintf("result[%s] amount", id), t.Name())
	}

	if keysNum != len(ret.match) {
		t.Fatal(fmt.Sprintf("result keyword[%s] num[%d != %d]", id, keysNum, len(ret.match)), t.Name())
	}

	tn := 0
	for _, _kt := range ret.match {
		tn += len(_kt)
	}
	if tn != textsNum {
		t.Fatal(fmt.Sprintf("result texts[%s] num[%d != %d]", id, textsNum, tn), t.Name())
	}

	for _t, _a := range textsAmount {
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

func _assertLabelResults(t *testing.T, rets labelResults, allTextsAmount int, resultsAmount int) {
	if resultsAmount != len(rets) {
		t.Fatal("label results len", t.Name())
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

	if amount != allTextsAmount {
		t.Fatal(fmt.Sprintf("label results amount[%d != %d]", allTextsAmount, amount), t.Name())
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
