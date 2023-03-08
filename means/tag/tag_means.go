package tag

import (
	"fmt"
	"strings"
)

// TagMeans
// tag means
type TagMeans struct {
	tagMatcher *ruleMatcher
	fn         func(*Matcher, []string) []*Result
}

func NewTagMeans(ruleName string, fn func(*Matcher, []string) []*Result, opts ...RuleMatcherOption) *TagMeans {
	t := &TagMeans{
		tagMatcher: newRuleMatcher(ruleName, opts...),
		fn:         fn,
	}

	return t
}

func (t *TagMeans) GetTitle() string {
	return fmt.Sprintf("Tag:%s{%s}", t.tagMatcher.key, strings.Join(t.tagMatcher.tagsName, ", "))
}

func (t *TagMeans) GetKeys() []string {
	return t.tagMatcher.getResultInsertKeys()
}

func (t *TagMeans) Close() {}

func (t *TagMeans) Insert(contents []string) []map[string]any {
	results := t.fn(t.tagMatcher.matcher, contents)
	if results == nil {
		return nil
	}

	return t.tagMatcher.resultsToSliceMap(results)
}

func (t *TagMeans) Update(contents []string) map[string]any {
	results := t.fn(t.tagMatcher.matcher, contents)
	if results == nil {
		return nil
	}

	return t.tagMatcher.resultToMap(results[0])
}

func NewKey(ruleName string, opts ...RuleMatcherOption) *TagMeans {
	t := NewTagMeans(ruleName, func(m *Matcher, c []string) []*Result {
		return m.MatchKey(c)
	}, opts...)

	return t
}

func NewMostText(ruleName string, opts ...RuleMatcherOption) *TagMeans {
	t := NewTagMeans(ruleName, func(m *Matcher, c []string) []*Result {
		return m.MatchMostText(c)
	}, opts...)

	return t
}

func NewMostKey(ruleName string, opts ...RuleMatcherOption) *TagMeans {
	t := NewTagMeans(ruleName, func(m *Matcher, c []string) []*Result {
		return m.MatchMostKey(c)
	}, opts...)

	return t
}
