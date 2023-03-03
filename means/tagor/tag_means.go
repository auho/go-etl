package tagor

import (
	"fmt"
	"strings"

	go_simple_db "github.com/auho/go-simple-db/v2"
)

// TagMeans
// tag means
//
type TagMeans struct {
	TagMatcher
}

func (tm *TagMeans) GetTitle() string {
	return fmt.Sprintf("Tag:%s{%s}", tm.key, strings.Join(tm.tagsName, ", "))
}

func (tm *TagMeans) GetKeys() []string {
	return tm.getResultInsertKeys()
}

func (tm *TagMeans) Close() {

}

func (tm *TagMeans) insertResult(f func() *Result) []map[string]interface{} {
	result := f()
	if result == nil {
		return nil
	}

	return tm.resultToSliceMap(result)
}

func (tm *TagMeans) insertResults(f func() []*Result) []map[string]interface{} {
	results := f()
	if results == nil {
		return nil
	}

	return tm.resultsToSliceMap(results)
}

func (tm *TagMeans) updateResult(f func() *Result) map[string]interface{} {
	result := f()
	if result == nil {
		return nil
	}

	return tm.resultToMap(result)
}

// Key
// Text
type Key struct {
	TagMeans
}

func NewKey(ruleName string, db *go_simple_db.SimpleDB, Options ...TagMatcherOption) *Key {
	t := &Key{}
	t.prepare(ruleName, db, Options...)

	return t
}

func (t *Key) Insert(contents []string) []map[string]interface{} {
	return t.insertResults(func() []*Result {
		return t.Matcher.MatchKey(contents)
	})
}

// MostText
// Most Text
//
type MostText struct {
	TagMeans
}

func NewMostText(ruleName string, db *go_simple_db.SimpleDB, Options ...TagMatcherOption) *MostText {
	t := &MostText{}
	t.prepare(ruleName, db, Options...)

	return t
}

func (t *MostText) Insert(contents []string) []map[string]interface{} {
	return t.insertResult(func() *Result {
		return t.Matcher.MatchMostText(contents)
	})
}

func (t *MostText) Update(contents []string) map[string]interface{} {
	return t.updateResult(func() *Result {
		return t.Matcher.MatchMostText(contents)
	})
}

// MostKey
// Most Key
//
type MostKey struct {
	TagMeans
}

func NewMostKey(ruleName string, db *go_simple_db.SimpleDB, Options ...TagMatcherOption) *MostKey {
	t := &MostKey{}
	t.prepare(ruleName, db, Options...)

	return t
}

func (t *MostKey) Insert(contents []string) []map[string]interface{} {
	return t.insertResult(func() *Result {
		return t.Matcher.MatchMostKey(contents)
	})
}

func (t *MostKey) Update(contents []string) map[string]interface{} {
	return t.updateResult(func() *Result {
		return t.Matcher.MatchMostKey(contents)
	})
}
