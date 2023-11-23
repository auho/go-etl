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

	// TODO extract func
	m.keys = []string{
		m.rule.NameAlias(),
		m.rule.KeywordNameAlias(),
		m.rule.KeywordNumNameAlias(),
	}
	m.keys = append(m.keys, m.rule.LabelsAlias()...)
	m.keys = append(m.keys, m.rule.FixedKeysAlias()...)

	// TODO extract func
	m.defaultValues = map[string]any{
		m.rule.NameAlias():           "",
		m.rule.KeywordNameAlias():    "",
		m.rule.KeywordNumNameAlias(): 0,
	}
	for _, _la := range m.rule.LabelsAlias() {
		m.defaultValues[_la] = ""
	}
	for _, _fka := range m.rule.FixedKeysAlias() {
		m.defaultValues[_fka] = ""
	}

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

// NewKey
// keyword
func NewKey(rule means.Ruler) *Means {
	t := NewMeans(rule, func(rule means.Ruler, m *Matcher, c []string) []map[string]any {
		rs := m.MatchKey(c)
		if rs == nil {
			return nil
		}

		return rs.toSliceMapAny(rule)
	})

	return t
}

// NewLabel
// label tags
func NewLabel(rule means.Ruler) *Means {
	t := NewMeans(rule, func(rule means.Ruler, m *Matcher, c []string) []map[string]any {
		rs := m.MatchLabel(c)
		if rs == nil {
			return nil
		}

		return rs.toSliceMapAny(rule)
	})

	return t
}

func NewMostText(rule means.Ruler) *Means {
	t := NewMeans(rule, func(rule means.Ruler, m *Matcher, c []string) []map[string]any {
		rs := m.MatchMostText(c)
		if rs == nil {
			return nil
		}

		return rs.toSliceMapAny(rule)
	})

	return t
}

func NewMostKey(rule means.Ruler) *Means {
	t := NewMeans(rule, func(rule means.Ruler, m *Matcher, c []string) []map[string]any {
		rs := m.MatchMostKey(c)
		if rs == nil {
			return nil
		}

		return rs.toSliceMapAny(rule)
	})

	return t
}

func NewFirst(rule means.Ruler) *Means {
	t := NewMeans(rule, func(rule means.Ruler, m *Matcher, c []string) []map[string]any {
		rs := m.MatchFirstText(c)
		if rs == nil {
			return nil
		}

		return rs.toSliceMapAny(rule)
	})

	return t
}
