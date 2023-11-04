package tag

import (
	"fmt"
	"strings"
)

// TagMeans
// tag means
type TagMeans struct {
	rule    Ruler
	matcher *Matcher
	fn      func(*Matcher, []string) []*Result
}

func NewTagMeans(rule Ruler, fn func(*Matcher, []string) []*Result) *TagMeans {
	tm := &TagMeans{
		rule: rule,
		fn:   fn,
	}

	return tm
}

func (t *TagMeans) GetTitle() string {
	return fmt.Sprintf("Tag:%s{%s}", t.rule.Name(), strings.Join(t.rule.Labels(), ", "))
}

func (t *TagMeans) GetKeys() []string {
	return append([]string{t.rule.KeywordName(), t.rule.KeywordNumName()}, append(t.rule.Labels(), t.rule.FixedKeys()...)...)
}

func (t *TagMeans) Close() {}

func (t *TagMeans) Insert(contents []string) []map[string]any {
	results := t.fn(t.matcher, contents)
	if results == nil {
		return nil
	}

	return t.resultsToSliceMap(results)
}

func (t *TagMeans) Update(contents []string) map[string]any {
	results := t.fn(t.matcher, contents)
	if results == nil {
		return nil
	}

	return t.resultToMap(results[0])
}

func (t *TagMeans) resultsToSliceMap(results []*Result) []map[string]any {
	items := make([]map[string]any, 0, len(results))
	for _, result := range results {
		items = append(items, t.resultToMap(result))
	}

	return items
}

func (t *TagMeans) resultToSliceMap(result *Result) []map[string]any {
	return []map[string]any{t.resultToMap(result)}
}

func (t *TagMeans) resultToMap(result *Result) map[string]any {
	item := make(map[string]any)
	item[t.rule.KeywordNameAlias()] = result.Key
	item[t.rule.KeywordNumNameAlias()] = result.Num

	fixed := t.rule.FixedAlias()
	for _, key := range t.rule.FixedKeysAlias() {
		item[key] = fixed[key]
	}

	return item
}

func NewKey(rule Ruler) *TagMeans {
	t := NewTagMeans(rule, func(m *Matcher, c []string) []*Result {
		return m.MatchKey(c)
	})

	return t
}

func NewMostText(rule Ruler) *TagMeans {
	t := NewTagMeans(rule, func(m *Matcher, c []string) []*Result {
		return m.MatchMostText(c)
	})

	return t
}

func NewMostKey(rule Ruler) *TagMeans {
	t := NewTagMeans(rule, func(m *Matcher, c []string) []*Result {
		return m.MatchMostKey(c)
	})

	return t
}
