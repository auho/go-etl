package match

import (
	"fmt"
	"sort"
	"strings"

	"github.com/auho/go-etl/v2/job/means"
	"github.com/auho/go-etl/v2/tool/maps"
)

var DefaultFormat = Format{
	WithKeyWordAmount: true,
	Sep:               ",",
}

type Format struct {
	WithKeyWordAmount bool
	Sep               string
}

// Result
// result 匹配结果
type Result struct {
	Amount  int               // matched amount
	Keyword string            // keyword
	Tags    map[string]string // tags map[tag name]tag value
	Texts   map[string]int    // matched text map[matched text]amount
}

func NewResult() Result {
	m := Result{}
	m.Tags = make(map[string]string)
	m.Texts = make(map[string]int)

	return m
}

func (r *Result) ToTag(rule means.Ruler) map[string]any {
	item := make(map[string]any)

	for _k, _v := range r.Tags {
		item[_k] = _v
	}

	item[rule.KeywordNameAlias()] = r.Keyword
	item[rule.KeywordNumNameAlias()] = 1
	item[rule.KeywordAmountNameAlias()] = r.Amount

	return item
}

// Results
// result
type Results []Result

func (rs Results) ToAll(rule means.Ruler) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias(), rule.KeywordAmountNameAlias())

	items := make([]map[string]any, 0, len(rs))
	for _, _r := range rs {
		items = append(items, maps.PluckMap(_r.ToTag(rule), keys))
	}

	return items
}

func (rs Results) ToLine(rule means.Ruler, format Format) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias(), rule.KeywordNumNameAlias())
	m := rs.MergeKeysToWhole(rule, format)

	return []map[string]any{maps.PluckMap(m, keys)}
}

func (rs Results) ToFlag(rule means.Ruler, format Format) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias())
	m := rs.MergeKeysToWhole(rule, format)
	m[rule.NameAlias()] = 1

	return []map[string]any{maps.PluckMap(m, keys)}
}

func (rs Results) MergeKeysToWhole(rule means.Ruler, format Format) map[string]any {
	keyNum := 0
	keyAmount := 0
	tagsValues := make(map[string][]string)
	for _, _r := range rs {
		for _ta, _tv := range _r.Tags {
			tagsValues[_ta] = append(tagsValues[_ta], _tv)
		}

		keyNum += 1
		keyAmount += _r.Amount

		var keywordText string
		if format.WithKeyWordAmount {
			keywordText = fmt.Sprintf("%s %d", _r.Keyword, _r.Amount)
		} else {
			keywordText = _r.Keyword
		}

		tagsValues[rule.KeywordNameAlias()] = append(tagsValues[rule.KeywordNameAlias()], keywordText)
	}

	m := make(map[string]any)
	for _tn, _tv := range tagsValues {
		m[_tn] = strings.Join(_tv, format.Sep)
	}

	m[rule.KeywordNumNameAlias()] = keyNum
	m[rule.KeywordAmountNameAlias()] = keyAmount

	return m
}

// LabelResult
// label result
type LabelResult struct {
	Identity string
	Amount   int                       // match amount
	Tags     map[string]string         // tags map[tag name]tag value
	Match    map[string]map[string]int // keyword and match text map[keyword]map[matched text]num
	Keywords []string                  // []keyword
}

func NewLabelResult() LabelResult {
	l := LabelResult{}
	l.Tags = make(map[string]string)
	l.Match = make(map[string]map[string]int)

	return l
}

func (lr *LabelResult) ToTag(rule means.Ruler, format Format) map[string]any {
	m := make(map[string]any)

	for _tn, _tv := range lr.Tags {
		m[_tn] = _tv
	}

	keyNum := 0
	keyAmount := 0
	var keysValue []string
	for _, _key := range lr.Keywords {
		_textAmount := 0
		for _, _a := range lr.Match[_key] {
			_textAmount += _a
		}

		keyNum += 1
		keyAmount += _textAmount

		var keyText string
		if format.WithKeyWordAmount {
			keyText = fmt.Sprintf("%s %d", _key, _textAmount)
		} else {
			keyText = _key
		}

		keysValue = append(keysValue, keyText)
	}

	m[rule.KeywordNameAlias()] = strings.Join(keysValue, format.Sep)
	m[rule.KeywordNumNameAlias()] = keyNum
	m[rule.KeywordAmountNameAlias()] = keyAmount

	return m
}

// LabelResults
// label results
type LabelResults []LabelResult

func (lrs LabelResults) ToAll(rule means.Ruler, format Format) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias(), rule.KeywordAmountNameAlias())

	items := make([]map[string]any, 0, len(lrs))
	for _, _r := range lrs {
		items = append(items, maps.PluckMap(_r.ToTag(rule, format), keys))
	}

	return items
}

func (lrs LabelResults) ToLine(rule means.Ruler, format Format) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias(), rule.LabelNumNameAlias(), rule.KeywordNumNameAlias(), rule.KeywordAmountNameAlias())
	m := lrs.MergeLabelsToWhole(rule, format)

	return []map[string]any{maps.PluckMap(m, keys)}
}

func (lrs LabelResults) ToFlag(rule means.Ruler, format Format) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias())
	m := lrs.MergeLabelsToWhole(rule, format)
	m[rule.NameAlias()] = 1

	return []map[string]any{maps.PluckMap(m, keys)}
}

func (lrs LabelResults) MergeLabelsToWhole(rule means.Ruler, format Format) map[string]any {
	sort.SliceStable(lrs, func(i, j int) bool {
		return lrs[i].Identity < lrs[j].Identity
	})

	labelNum := 0
	labelAmount := 0
	keywordNum := 0
	keywordAmount := 0

	tagsValues := make(map[string][]string)

	for _, _lr := range lrs {
		for _ta, _tv := range _lr.Tags {
			tagsValues[_ta] = append(tagsValues[_ta], _tv)
		}

		var keysValue []string
		for _, _key := range _lr.Keywords {
			_keyAmount := 0
			for _, _a := range _lr.Match[_key] {
				_keyAmount += _a
			}

			keywordNum += 1
			keywordAmount += _keyAmount
			keysValue = append(keysValue, _key)
		}

		tagsValues[rule.KeywordNameAlias()] = append(tagsValues[rule.KeywordNameAlias()], strings.Join(keysValue, format.Sep))

		labelNum += 1
		labelAmount += _lr.Amount
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
