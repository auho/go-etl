package match

import (
	"fmt"
	"strings"
	"testing"

	"github.com/auho/go-etl/v2/job/means"
)

func TestMeans(t *testing.T) {
	_means := means.NewMeans(NewSearchKey(_rule, NewExportKeywordAll(_rule)))
	err := _means.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	keys := _means.GetKeys()
	if len(keys) < 3 {
		t.Fatal()
	}
}

func _genMeans[T ResultsEntity](t *testing.T, fn func(means.Ruler) *Search[T]) *means.Means {
	_means := means.NewMeans(
		fn(_rule).
			WithIgnoreCase().
			WithPriorityFuzzy().
			WithFuzzy(FuzzyConfig{
				Window: 3,
				Sep:    "_",
			}).
			WithDebug(),
	)
	err := _means.Prepare()
	if err != nil {
		t.Fatal(err, t.Name())
	}

	return _means
}

func TestNewFirstText(t *testing.T) {
	_means := _genMeans(t, NewFirstText)
	rets := _means.Insert(_contents)
	_outputResults(rets)

	_assertTags(t, _rule, rets, 1, 1)

	_assertTag(t, _rule, rets[0], "b", 1)
}

func TestNewMostText(t *testing.T) {
	_means := _genMeans(t, NewMostText)
	rets := _means.Insert(_contents)
	_outputResults(rets)

	_assertTags(t, _rule, rets, 1, 12)

	_assertTag(t, _rule, rets[0], "123", 12)
}

func TestNewKey(t *testing.T) {
	_means := _genMeans(t, NewKey)
	rets := _means.Insert(_contents)
	_outputResults(rets)

	_assertTags(t, _rule, rets, 3, 41)

	_assertTag(t, _rule, rets[0], "123", 12)
	_assertTag(t, _rule, rets[1], "b", 8)
	_assertTag(t, _rule, rets[2], "中_文", 21)
}

func TestNewFirstKey(t *testing.T) {
	_means := _genMeans(t, NewFirstKey)
	rets := _means.Insert(_contents)
	_outputResults(rets)

	_assertTags(t, _rule, rets, 1, 1)

	_assertTag(t, _rule, rets[0], "123", 1)
}

func TestNewMostKey(t *testing.T) {
	_means := _genMeans(t, NewMostKey)
	rets := _means.Insert(_contents)
	_outputResults(rets)

	_assertTags(t, _rule, rets, 1, 21)

	_assertTag(t, _rule, rets[0], "中_文", 21)
}

func TestNewWholeLabels(t *testing.T) {
	_means := _genMeans(t, NewWholeLabels)
	rets := _means.Insert(_contents)
	_outputResults(rets)

	_assertTags(t, _rule, rets, 1, 41)

	_assertTag(t, _rule, rets[0], "123|b|中_文", 41)

	if rets[0][_rule.LabelNumNameAlias()] != 3 || rets[0][_rule.KeywordNumNameAlias()] != 3 {
		t.Fatal()
	}

	if rets[0][_rule.LabelNumNameAlias()] != len(strings.Split(rets[0]["a"].(string), "|")) {
		t.Fatal("label num")
	}

	amount := 0
	keyword := rets[0][_rule.KeywordNameAlias()].(string)
	for _, _s := range strings.Split(keyword, "|") {
		amount += len(strings.Split(_s, ","))
	}

	if rets[0][_rule.KeywordNumNameAlias()] != amount {
		t.Fatal("keyword num")
	}
}

func TestNewLabel(t *testing.T) {
	_means := _genMeans(t, NewLabel)
	rets := _means.Insert(_contents)
	_outputResults(rets)

	_assertTags(t, _rule, rets, 3, 41)

	_assertTagLabel(t, _rule, rets[0], "123", 12)
	_assertTagLabel(t, _rule, rets[1], "b", 8)
	_assertTagLabel(t, _rule, rets[2], "中_文", 21)
}

func _assertTag(t *testing.T, rule means.Ruler, m map[string]any, keyword string, expectAmount int) {
	if m[rule.KeywordNameAlias()] != keyword {
		t.Fatal(fmt.Sprintf("keyword[%s != %s]", keyword, m[rule.KeywordNameAlias()]), t.Name())
	}

	if m[rule.KeywordAmountNameAlias()] != expectAmount {
		t.Fatal(fmt.Sprintf("keyword[%s] amount[%d != %d]", keyword, expectAmount, m[rule.KeywordAmountNameAlias()]), t.Name())
	}
}

func _assertTagLabel(t *testing.T, rule means.Ruler, m map[string]any, keyword string, expectAmount int) {
	_ky := fmt.Sprintf("%s %d", keyword, expectAmount)
	if m[rule.KeywordNameAlias()] != _ky {
		t.Fatal(fmt.Sprintf("keyword[%s != %s]", keyword, m[rule.KeywordNameAlias()]), t.Name())
	}

	if m[rule.KeywordAmountNameAlias()] != expectAmount {
		t.Fatal(fmt.Sprintf("keyword[%s] amount[%d != %d]", keyword, expectAmount, m[rule.KeywordAmountNameAlias()]), t.Name())
	}
}

func _assertTags(t *testing.T, rule means.Ruler, sm []map[string]any, expectNum, expectAmount int) {
	amount := 0
	for _, m := range sm {
		amount += m[rule.KeywordAmountNameAlias()].(int)
	}

	if len(sm) != expectNum {
		t.Fatal(fmt.Sprintf("num[%d != %d]", expectNum, len(sm)), t.Name())
	}

	if amount != expectAmount {
		t.Fatal(fmt.Sprintf("num[%d != %d]", expectAmount, amount), t.Name())
	}
}
