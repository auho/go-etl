package match

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/means"
)

// Result
// result 匹配结果
type Result struct {
	Key    string            // keyword
	Num    int               // matched num
	Labels map[string]string // labels[label name]label value
}

func NewResult() *Result {
	m := &Result{}
	m.Labels = make(map[string]string)

	return m
}

func (r *Result) toMapAny(rule means.Ruler) map[string]any {
	item := make(map[string]any)

	for _k, _v := range r.Labels {
		item[_k] = _v
	}

	item[rule.KeywordNameAlias()] = r.Key
	item[rule.KeywordNumNameAlias()] = r.Num

	fixed := rule.FixedAlias()
	for _, key := range rule.FixedKeysAlias() {
		item[key] = fixed[key]
	}

	return item
}

// LabelResult
// label result
type LabelResult struct {
	Identity    string
	Labels      map[string]string // labels map[label name]label value
	Match       map[string]int    // keyword num map[keyword]num
	Keys        []string          // []keyword
	MatchAmount int               // match amount
}

func NewLabelResult() *LabelResult {
	l := &LabelResult{}
	l.Labels = make(map[string]string)
	l.Match = make(map[string]int)

	return l
}

func (lr *LabelResult) toMapAny(rule means.Ruler) map[string]any {
	item := make(map[string]any)

	for _name, _value := range lr.Labels {
		item[_name] = _value
	}

	var _keywords []string
	for _, _key := range lr.Keys {
		_keywords = append(_keywords, fmt.Sprintf("%s %d", _key, lr.Match[_key]))
	}

	item[rule.KeywordNameAlias()] = strings.Join(_keywords, ",")
	item[rule.KeywordNumNameAlias()] = lr.MatchAmount

	fixed := rule.FixedAlias()
	for _, key := range rule.FixedKeysAlias() {
		item[key] = fixed[key]
	}

	return item
}

// Results
// result
type Results []*Result

func (rs Results) Len() int {
	return len(rs)
}

func (rs Results) Less(i, j int) bool {
	return rs[i].Num > rs[j].Num
}

func (rs Results) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs Results) toSliceMapAny(rule means.Ruler) []map[string]any {
	items := make([]map[string]any, 0, rs.Len())
	for _, _r := range rs {
		items = append(items, _r.toMapAny(rule))
	}

	return items
}

// LabelResults
// label results
type LabelResults []*LabelResult

func (lrs LabelResults) Len() int {
	return len(lrs)
}

func (lrs LabelResults) Less(i, j int) bool {
	return lrs[i].MatchAmount > lrs[j].MatchAmount
}

func (lrs LabelResults) Swap(i, j int) {
	lrs[i], lrs[j] = lrs[j], lrs[i]
}

func (lrs LabelResults) toSliceMapAny(rule means.Ruler) []map[string]any {
	items := make([]map[string]any, 0, lrs.Len())
	for _, _r := range lrs {
		items = append(items, _r.toMapAny(rule))
	}

	return items
}
