package tag

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/means"
)

var _ means.InsertMeans = (*Means)(nil)
var _ means.UpdateMeans = (*Means)(nil)

// Means
// tag means
type Means struct {
	rule    means.Ruler
	matcher *Matcher
	fn      func(means.Ruler, *Matcher, []string) []map[string]any

	keys          []string       // output name
	defaultValues map[string]any // output default values
}

func NewMeans(rule means.Ruler, fn func(means.Ruler, *Matcher, []string) []map[string]any) *Means {
	m := &Means{
		rule: rule,
		fn:   fn,
	}

	return m
}

func (m *Means) Prepare() error {
	m.matcher = DefaultMatcher()

	items, err := m.rule.ItemsForRegexp()
	if err != nil {
		return fmt.Errorf("ItemsForRegexp error; %w", err)
	}

	m.matcher.prepare(m.rule.KeywordNameAlias(), items)

	m.keys = m.rule.MeansKeys()
	m.defaultValues = m.rule.MeansDefaultValues()

	return nil
}

func (m *Means) GetTitle() string {
	return fmt.Sprintf("Tag:%s{%s}", m.rule.Name(), strings.Join(m.rule.Labels(), ", "))
}

func (m *Means) GetKeys() []string {
	return m.keys
}

func (m *Means) Close() error {
	return nil
}

func (m *Means) Insert(contents []string) []map[string]any {
	return m.fn(m.rule, m.matcher, contents)
}

func (m *Means) Update(contents []string) map[string]any {
	results := m.fn(m.rule, m.matcher, contents)
	if results == nil {
		return nil
	}

	return results[0]
}

func (m *Means) DefaultValues() map[string]any {
	return m.defaultValues
}

// NewWholeLabels
// merge all labels together
// label1|label2|label3
// keyword1|keyword2|keyword3|
func NewWholeLabels(rule means.Ruler) *Means {
	return NewMeans(rule, func(rule means.Ruler, m *Matcher, c []string) []map[string]any {
		rs := m.MatchLabel(c)
		if rs == nil {
			return nil
		}

		_rts := make(map[string][]string)
		_labelAmount := 0
		_keywordAmount := 0

		for _, _r := range rs {
			for _labelKey, _labelValue := range _r.Labels {
				_rts[_labelKey] = append(_rts[_labelKey], _labelValue)
			}

			for _, _key := range _r.Keys {
				_rts[rule.KeywordNameAlias()] = append(_rts[rule.KeywordNameAlias()], _key)

				_keywordAmount += 1
			}

			_labelAmount += 1
		}

		_rt := make(map[string]any)
		for _rk, _rv := range _rts {
			_rt[_rk] = strings.Join(_rv, "|")
		}

		_rt[rule.LabelNumNameAlias()] = _labelAmount
		_rt[rule.KeywordNumNameAlias()] = _keywordAmount

		return []map[string]any{_rt}
	})
}

// NewKey
// keyword
func NewKey(rule means.Ruler) *Means {
	return NewMeans(rule, func(rule means.Ruler, m *Matcher, c []string) []map[string]any {
		rs := m.MatchKey(c)
		if rs == nil {
			return nil
		}

		return rs.toSliceMapAny(rule)
	})
}

// NewLabel
// label tags
func NewLabel(rule means.Ruler) *Means {
	return NewMeans(rule, func(rule means.Ruler, m *Matcher, c []string) []map[string]any {
		rs := m.MatchLabel(c)
		if rs == nil {
			return nil
		}

		return rs.toSliceMapAny(rule)
	})
}

// NewMostText
// most text
func NewMostText(rule means.Ruler) *Means {
	return NewMeans(rule, func(rule means.Ruler, m *Matcher, c []string) []map[string]any {
		rs := m.MatchMostText(c)
		if rs == nil {
			return nil
		}

		return rs.toSliceMapAny(rule)
	})
}

// NewMostKey
// most key
func NewMostKey(rule means.Ruler) *Means {
	return NewMeans(rule, func(rule means.Ruler, m *Matcher, c []string) []map[string]any {
		rs := m.MatchMostKey(c)
		if rs == nil {
			return nil
		}

		return rs.toSliceMapAny(rule)
	})
}

// NewFirst
// the first part of the text is matched
func NewFirst(rule means.Ruler) *Means {
	return NewMeans(rule, func(rule means.Ruler, m *Matcher, c []string) []map[string]any {
		rs := m.MatchFirstText(c)
		if rs == nil {
			return nil
		}

		return rs.toSliceMapAny(rule)
	})
}
