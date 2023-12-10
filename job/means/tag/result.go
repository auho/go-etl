package tag

import (
	"fmt"
	"sort"
	"strings"

	"github.com/auho/go-etl/v2/job/means"
)

func attachFixesValues(rule means.Ruler, m map[string]any) map[string]any {
	fixed := rule.FixedAlias()
	for _, key := range rule.FixedKeysAlias() {
		m[key] = fixed[key]
	}

	return m
}

func pluckFromSliceMap(sm []map[string]any, keys []string) []map[string]any {
	var nsm []map[string]any
	for _, m := range sm {
		nsm = append(nsm, pluckFromMap(m, keys))
	}

	return nsm
}

func pluckFromMap(m map[string]any, keys []string) map[string]any {
	nm := make(map[string]any)
	for _, key := range keys {
		nm[key] = m[key]
	}

	return nm
}

// Result
// result 匹配结果
type Result struct {
	Amount     int               // matched amount
	Keyword    string            // keyword
	Tags       map[string]string // tags map[tag name]tag
	Texts      map[string]int    // matched text map[matched text]amount
	IsKeyMerge bool
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

	return attachFixesValues(rule, item)
}

// Results
// result
type Results []Result

func (rs Results) Len() int {
	return len(rs)
}

func (rs Results) Less(i, j int) bool {
	return rs[i].Amount > rs[j].Amount
}

func (rs Results) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs Results) ToTags(rule means.Ruler) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias(), rule.KeywordAmountNameAlias())
	items := make([]map[string]any, 0, rs.Len())
	for _, _r := range rs {
		items = append(items, pluckFromMap(_r.ToTag(rule), keys))
	}

	return items
}

func (rs Results) ToLine(rule means.Ruler) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias())
	sm := rs.MergeKeys(rule)

	return pluckFromSliceMap(sm, keys)
}

func (rs Results) ToFlag(rule means.Ruler) []map[string]any {
	sm := rs.ToLine(rule)
	sm[0][rule.NameAlias()] = 1

	return sm
}

func (rs Results) MergeKeys(rule means.Ruler) []map[string]any {
	keyNum := 0
	keyAmount := 0
	tagsValues := make(map[string][]string)
	for _, _r := range rs {
		for _, _ta := range rule.TagsAlias() {
			tagsValues[_ta] = append(tagsValues[_ta], _r.Tags[_ta])
		}

		keyNum += 1
		keyAmount += _r.Amount
		tagsValues[rule.KeywordNameAlias()] = append(tagsValues[rule.KeywordNameAlias()], fmt.Sprintf("%s %d", _r.Keyword, _r.Amount))
	}

	m := make(map[string]any)
	for _tn, _tv := range tagsValues {
		m[_tn] = strings.Join(_tv, "|")
	}

	m[rule.KeywordNumNameAlias()] = keyNum
	m[rule.KeywordAmountNameAlias()] = keyAmount

	m = attachFixesValues(rule, m)

	return []map[string]any{m}
}

// LabelResult
// label result
type LabelResult struct {
	Identity string
	Amount   int                       // match amount
	Tags     map[string]string         // tags map[tag name]tag
	Match    map[string]map[string]int // keyword and match text map[keyword]map[matched text]num
	Keys     []string                  // []keyword
}

func NewLabelResult() LabelResult {
	l := LabelResult{}
	l.Tags = make(map[string]string)
	l.Match = make(map[string]map[string]int)

	return l
}

func (lr *LabelResult) ToTag(rule means.Ruler) map[string]any {
	m := make(map[string]any)

	for _tn, _tv := range lr.Tags {
		m[_tn] = _tv
	}

	keyNum := 0
	keyAmount := 0
	var keysValue []string
	for _, _key := range lr.Keys {
		_textAmount := 0
		for _, _a := range lr.Match[_key] {
			_textAmount += _a
		}

		keyNum += 1
		keyAmount += _textAmount
		keysValue = append(keysValue, fmt.Sprintf("%s %d", _key, _textAmount))
	}

	m[rule.KeywordNameAlias()] = strings.Join(keysValue, ",")
	m[rule.KeywordNumNameAlias()] = keyNum
	m[rule.KeywordAmountNameAlias()] = keyAmount

	return attachFixesValues(rule, m)
}

// LabelResults
// label results
type LabelResults []LabelResult

func (lrs LabelResults) Len() int {
	return len(lrs)
}

func (lrs LabelResults) Less(i, j int) bool {
	return lrs[i].Amount > lrs[j].Amount
}

func (lrs LabelResults) Swap(i, j int) {
	lrs[i], lrs[j] = lrs[j], lrs[i]
}

func (lrs LabelResults) ToTags(rule means.Ruler) []map[string]any {
	keys := append(rule.Tags(), rule.KeywordNameAlias(), rule.KeywordAmountNameAlias())

	items := make([]map[string]any, 0, lrs.Len())
	for _, _r := range lrs {
		items = append(items, pluckFromMap(_r.ToTag(rule), keys))
	}

	return items
}

func (lrs LabelResults) ToLine(rule means.Ruler) []map[string]any {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias())
	m := lrs.MergeLabels(rule)

	return []map[string]any{pluckFromMap(m, keys)}
}

func (lrs LabelResults) ToFlag(rule means.Ruler) []map[string]any {
	sm := lrs.ToLine(rule)
	sm[0][rule.NameAlias()] = 1

	return sm
}

func (lrs LabelResults) MergeLabels(rule means.Ruler) map[string]any {
	sort.SliceStable(lrs, func(i, j int) bool {
		return lrs[i].Identity < lrs[j].Identity
	})

	labelNum := 0
	labelAmount := 0
	keywordNum := 0
	keywordAmount := 0

	tagsValues := make(map[string][]string)

	for _, _lr := range lrs {
		for _, _ta := range rule.TagsAlias() {
			tagsValues[_ta] = append(tagsValues[_ta], _lr.Tags[_ta])
		}

		var keysValue []string
		for _, _key := range _lr.Keys {
			_keyAmount := 0
			for _, _a := range _lr.Match[_key] {
				_keyAmount += _a
			}

			keywordNum += 1
			keywordAmount += _keyAmount
			keysValue = append(keysValue, _key)
		}

		tagsValues[rule.KeywordNameAlias()] = append(tagsValues[rule.KeywordNameAlias()], strings.Join(keysValue, ","))

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

	return attachFixesValues(rule, m)
}
