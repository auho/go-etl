package match

import (
	"fmt"
	"sort"
	"strings"

	"github.com/auho/go-etl/v3/job/extract"
	maps "github.com/auho/go-etl/v3/tool/mapx"
)

// defaultFormat is the default output format: keyword amount appended, comma-separated.
var defaultFormat = Format{
	withKeywordAmount: true,
	sep:               ",",
}

// Format controls how matched keywords are serialized.
type Format struct {
	withKeywordAmount bool   // if true, appends " <amount>" to each keyword
	sep               string // separator between multiple keyword values
}

// result holds the match outcome for a single keyword.
type result struct {
	amount  int               // matched amount
	keyword string            // keyword
	tags    map[string]string // tags map[tag name]tag value
	texts   map[string]int    // matched text map[matched text]amount
}

func newResult() result {
	m := result{}
	m.tags = make(map[string]string)
	m.texts = make(map[string]int)

	return m
}

func (r *result) toTag(rule extract.Rule) map[string]any {
	item := make(map[string]any)

	for _k, _v := range r.tags {
		item[_k] = _v
	}

	item[rule.KeywordNameAlias()] = r.keyword
	item[rule.KeywordNumNameAlias()] = 1
	item[rule.KeywordAmountNameAlias()] = r.amount

	return item
}

// results is a slice of result, used for keyword-oriented extraction.
type results []result

func (rs results) toAll(rule extract.Rule) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias(), rule.KeywordAmountNameAlias())

	items := make([]map[string]any, 0, len(rs))
	for _, _r := range rs {
		items = append(items, maps.PluckMap(_r.toTag(rule), keys))
	}

	return items
}

func (rs results) toLine(rule extract.Rule, format Format) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias(), rule.KeywordNumNameAlias())
	m := rs.mergeKeysToWhole(rule, format)

	return []map[string]any{maps.PluckMap(m, keys)}
}

func (rs results) toFlag(rule extract.Rule, format Format) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias())
	m := rs.mergeKeysToWhole(rule, format)
	m[rule.NameAlias()] = 1

	return []map[string]any{maps.PluckMap(m, keys)}
}

func (rs results) mergeKeysToWhole(rule extract.Rule, format Format) map[string]any {
	keyNum := 0
	keyAmount := 0
	tagsValues := make(map[string][]string)
	for _, _r := range rs {
		for _ta, _tv := range _r.tags {
			tagsValues[_ta] = append(tagsValues[_ta], _tv)
		}

		keyNum += 1
		keyAmount += _r.amount

		var keywordText string
		if format.withKeywordAmount {
			keywordText = fmt.Sprintf("%s %d", _r.keyword, _r.amount)
		} else {
			keywordText = _r.keyword
		}

		tagsValues[rule.KeywordNameAlias()] = append(tagsValues[rule.KeywordNameAlias()], keywordText)
	}

	m := make(map[string]any)
	for _tn, _tv := range tagsValues {
		m[_tn] = strings.Join(_tv, format.sep)
	}

	m[rule.KeywordNumNameAlias()] = keyNum
	m[rule.KeywordAmountNameAlias()] = keyAmount

	return m
}

// labelResult holds the match outcome for a label-oriented extraction,
// grouping matches by tag identity.
type labelResult struct {
	identity string
	amount   int                       // match amount
	tags     map[string]string         // tags map[tag name]tag value
	match    map[string]map[string]int // keyword and match text map[keyword]map[matched text]num
	keywords []string                  // []keyword
}

func newLabelResult() labelResult {
	l := labelResult{}
	l.tags = make(map[string]string)
	l.match = make(map[string]map[string]int)

	return l
}

func (lr *labelResult) toTag(rule extract.Rule, format Format) map[string]any {
	m := make(map[string]any)

	for _tn, _tv := range lr.tags {
		m[_tn] = _tv
	}

	keyNum := 0
	keyAmount := 0
	var keysValue []string
	for _, _key := range lr.keywords {
		_textAmount := 0
		for _, _a := range lr.match[_key] {
			_textAmount += _a
		}

		keyNum += 1
		keyAmount += _textAmount

		var keyText string
		if format.withKeywordAmount {
			keyText = fmt.Sprintf("%s %d", _key, _textAmount)
		} else {
			keyText = _key
		}

		keysValue = append(keysValue, keyText)
	}

	m[rule.KeywordNameAlias()] = strings.Join(keysValue, format.sep)
	m[rule.KeywordNumNameAlias()] = keyNum
	m[rule.KeywordAmountNameAlias()] = keyAmount

	return m
}

// labelResults is a slice of labelResult, used for label-oriented extraction.
type labelResults []labelResult

func (lrs labelResults) toAll(rule extract.Rule, format Format) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias(), rule.KeywordAmountNameAlias())

	items := make([]map[string]any, 0, len(lrs))
	for _, _r := range lrs {
		items = append(items, maps.PluckMap(_r.toTag(rule, format), keys))
	}

	return items
}

func (lrs labelResults) toLine(rule extract.Rule, format Format) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias(), rule.LabelNumNameAlias(), rule.KeywordNumNameAlias(), rule.KeywordAmountNameAlias())
	m := lrs.mergeLabelsToWhole(rule, format)

	return []map[string]any{maps.PluckMap(m, keys)}
}

func (lrs labelResults) toFlag(rule extract.Rule, format Format) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias())
	m := lrs.mergeLabelsToWhole(rule, format)
	m[rule.NameAlias()] = 1

	return []map[string]any{maps.PluckMap(m, keys)}
}

func (lrs labelResults) mergeLabelsToWhole(rule extract.Rule, format Format) map[string]any {
	sort.SliceStable(lrs, func(i, j int) bool {
		return lrs[i].identity < lrs[j].identity
	})

	labelNum := 0
	labelAmount := 0
	keywordNum := 0
	keywordAmount := 0

	tagsValues := make(map[string][]string)

	for _, _lr := range lrs {
		for _ta, _tv := range _lr.tags {
			tagsValues[_ta] = append(tagsValues[_ta], _tv)
		}

		var keysValue []string
		for _, _key := range _lr.keywords {
			_keyAmount := 0
			for _, _a := range _lr.match[_key] {
				_keyAmount += _a
			}

			keywordNum += 1
			keywordAmount += _keyAmount
			keysValue = append(keysValue, _key)
		}

		tagsValues[rule.KeywordNameAlias()] = append(tagsValues[rule.KeywordNameAlias()], strings.Join(keysValue, format.sep))

		labelNum += 1
		labelAmount += _lr.amount
	}

	m := make(map[string]any)
	for _tn, _tv := range tagsValues {
		m[_tn] = strings.Join(_tv, "|")
	}

	m[rule.LabelNumNameAlias()] = labelNum
	m[rule.KeywordNumNameAlias()] = keywordNum
	m[rule.KeywordAmountNameAlias()] = keywordAmount

	return m
}
