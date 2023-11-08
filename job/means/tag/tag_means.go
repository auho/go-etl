package tag

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/means"
)

var _ means.InsertMeans = (*TagMeans)(nil)

// TagMeans
// tag means
type TagMeans struct {
	rule    Ruler
	matcher *Matcher
	fn      func(*Matcher, []string) []*Result

	keys          []string
	defaultValues map[string]any
}

func NewTagMeans(rule Ruler, fn func(*Matcher, []string) []*Result) *TagMeans {
	tm := &TagMeans{
		rule: rule,
		fn:   fn,
	}

	return tm
}

func (t *TagMeans) Prepare() error {
	if t.matcher == nil {
		t.matcher = DefaultMatcher()
	}

	items, err := t.rule.ItemsForRegexp()
	if err != nil {
		return fmt.Errorf("ItemsForRegexp error; %w", err)
	}

	t.matcher.prepare(t.rule.KeywordNameAlias(), items)

	t.keys = []string{
		t.rule.NameAlias(),
		t.rule.KeywordNameAlias(),
		t.rule.KeywordNumNameAlias(),
	}
	t.keys = append(t.keys, t.rule.LabelsAlias()...)
	t.keys = append(t.keys, t.rule.FixedKeysAlias()...)

	t.defaultValues = map[string]any{
		t.rule.NameAlias():           "",
		t.rule.KeywordNameAlias():    "",
		t.rule.KeywordNumNameAlias(): 0,
	}
	for _, _la := range t.rule.LabelsAlias() {
		t.defaultValues[_la] = ""
	}
	for _, _fka := range t.rule.FixedKeysAlias() {
		t.defaultValues[_fka] = ""
	}

	return nil
}

func (t *TagMeans) GetTitle() string {
	return fmt.Sprintf("Tag:%s{%s}", t.rule.Name(), strings.Join(t.rule.Labels(), ", "))
}

func (t *TagMeans) GetKeys() []string {
	return t.keys
}

func (t *TagMeans) Close() error {
	return nil
}

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

func (t *TagMeans) DefaultValues() map[string]any {
	return t.defaultValues
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

	for _k, _v := range result.Tags {
		item[_k] = _v
	}

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
