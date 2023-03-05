package tagor

import (
	"fmt"
	"strings"

	go_simple_db "github.com/auho/go-simple-db/v2"
)

// TagMeans
// tag means
type TagMeans struct {
	tagMatcher *TagMatcher
	fn         func(*Matcher, []string) []*Result
}

func NewTagMeans(ruleName string, db *go_simple_db.SimpleDB, fn func(*Matcher, []string) []*Result, opts ...TagMatcherOption) *TagMeans {
	t := &TagMeans{
		tagMatcher: newTagMatcher(ruleName, db, opts...),
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

func (t *TagMeans) Insert(contents []string) []map[string]interface{} {
	results := t.fn(t.tagMatcher.Matcher, contents)
	if results == nil {
		return nil
	}

	return t.tagMatcher.resultsToSliceMap(results)
}

func (t *TagMeans) Update(contents []string) map[string]interface{} {
	results := t.fn(t.tagMatcher.Matcher, contents)
	if results == nil {
		return nil
	}

	return t.tagMatcher.resultToMap(results[0])
}

func NewKey(ruleName string, db *go_simple_db.SimpleDB, opts ...TagMatcherOption) *TagMeans {
	t := NewTagMeans(ruleName, db, func(m *Matcher, c []string) []*Result {
		return m.MatchKey(c)
	}, opts...)

	return t
}

func NewMostText(ruleName string, db *go_simple_db.SimpleDB, opts ...TagMatcherOption) *TagMeans {
	t := NewTagMeans(ruleName, db, func(m *Matcher, c []string) []*Result {
		return m.MatchMostText(c)
	}, opts...)

	return t
}

func NewMostKey(ruleName string, db *go_simple_db.SimpleDB, opts ...TagMatcherOption) *TagMeans {
	t := NewTagMeans(ruleName, db, func(m *Matcher, c []string) []*Result {
		return m.MatchMostKey(c)
	}, opts...)

	return t
}
