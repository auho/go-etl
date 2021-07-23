package means

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/auho/go-simple-db/simple"
)

// Result
// result
//
type Result struct {
	Key           string
	Num           int64
	Texts         map[string]int64
	Tags          map[string]string
	IsTextComplex bool
}

func NewResult() *Result {
	m := &Result{}
	m.Tags = make(map[string]string)
	m.Texts = make(map[string]int64)

	return m
}

// TagResult
// tag result
//
type TagResult struct {
	Identity string
	Tags     map[string]string
	Match    map[string]map[string]int
	KeyNum   int64
	TextNum  int64
}

func NewTagResult() *TagResult {
	m := &TagResult{}
	m.Tags = make(map[string]string)
	m.Match = make(map[string]map[string]int)

	return m
}

type TagMatcherOption func(*TagMatcher)

func WithTagMatcherData(d string) TagMatcherOption {
	return func(t *TagMatcher) {
		t.dataName = d
	}
}

func WittTagMatcherAlias(alias map[string]string) TagMatcherOption {
	return func(t *TagMatcher) {
		t.alias = alias
	}
}

// TagMatcher
// tag matcher
//
type TagMatcher struct {
	Matcher          *Matcher
	db               simple.Driver
	tagsName         []string
	key              string
	keyFieldName     string
	keyNumFieldName  string
	keyTableName     string
	dataName         string
	alias            map[string]string
	excludeRuleField []string
}

func NewTagMatcher(key string, db simple.Driver, Options ...TagMatcherOption) *TagMatcher {
	t := &TagMatcher{}
	t.key = key
	t.db = db
	t.excludeRuleField = []string{"id", "keyword_len"}

	for _, option := range Options {
		option(t)
	}

	return t
}

func (t *TagMatcher) getRules() []map[string]string {
	row, err := t.db.GetTableColumns(t.keyTableName)
	if err != nil {
		panic(err)
	}

	columns := make([]string, 0, len(row))
	for k := range row {
		column := row[k].(string)
		for _, ec := range t.excludeRuleField {
			if ec == column {
				goto LOOP
			}
		}

		columns = append(columns)
	LOOP:
	}

	query := fmt.Sprintf("SELECT `%s` FROM `%s` ORDER BY `keyword_len` DESC, `id` ASC", strings.Join(columns, "`, `"), t.keyTableName)
	rules, err := t.db.QueryString(query)
	if err != nil {
		panic(err)
	}

	if len(rules) <= 0 {
		panic("rules is null")
	}

	return rules
}

func (t *TagMatcher) GetResultInsertKeys() []string {
	return append([]string{t.keyFieldName, t.keyNumFieldName}, t.tagsName...)
}

func (t *TagMatcher) MatchUniqueText(contents []string) [][]interface{} {
	results := t.Matcher.MatchText(contents)
	if results == nil {
		return nil
	}

	return t.ResultsToSliceSlice(results)
}

func (t *TagMatcher) MatchUniqueKey(contents []string) [][]interface{} {
	results := t.Matcher.MatchKey(contents)
	if results == nil {
		return nil
	}

	return t.ResultsToSliceSlice(results)
}

func (t *TagMatcher) ResultsToSliceMap(results []*Result) []map[string]interface{} {
	items := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		items = append(items, t.ResultToMap(result))
	}

	return items
}

func (t *TagMatcher) ResultToMap(result *Result) map[string]interface{} {
	item := make(map[string]interface{})
	item[t.keyFieldName] = result.Key
	item[t.keyNumFieldName] = result.Num

	for _, tagName := range t.tagsName {
		item[tagName] = result.Tags[tagName]
	}

	return item
}

func (t *TagMatcher) ResultToSlice(result *Result) []interface{} {
	item := make([]interface{}, 0, len(t.tagsName)+2)
	item = append(item, result.Key)
	item = append(item, result.Num)

	for _, tagName := range t.tagsName {
		item = append(item, result.Tags[tagName])
	}

	return item
}

func (t *TagMatcher) ResultToSliceSlice(result *Result) [][]interface{} {
	item := t.ResultToSlice(result)

	return [][]interface{}{item}
}

func (t *TagMatcher) ResultsToSliceSlice(results []*Result) [][]interface{} {
	items := make([][]interface{}, 0, len(results))
	for _, result := range results {
		items = append(items, t.ResultToSlice(result))
	}

	return items
}

// MatcherOption
// tag match option
//
type MatcherOption func(mt *Matcher)

func WithTagMatcherKeyFun(f func(string) string) MatcherOption {
	return func(tm *Matcher) {
		tm.addKeyFormatFun(f)
	}
}

// Matcher
// matcher
//
type Matcher struct {
	keyFormatFunList []func(string) string
	regexpItems      map[string]map[string]string
	regexp           *regexp.Regexp
	regexpString     string
	normalRegexpName string
	badKeyMap        map[string]string
	tagsName         []string
}

func NewMatcher(keyName string, items []map[string]string, Options ...MatcherOption) *Matcher {
	m := &Matcher{}
	m.normalRegexpName = "_rEgEx_"
	m.regexpItems = make(map[string]map[string]string, len(items))
	m.badKeyMap = make(map[string]string)

	for _, option := range Options {
		option(m)
	}

	m.init(keyName, items)

	return m
}

func (m *Matcher) init(keyName string, items []map[string]string) {
	for k := range items[0] {
		if k != keyName {
			m.tagsName = append(m.tagsName, k)
		}
	}

	regexpNormalItems := make([]string, 0)
	regexpGroupItems := make([]string, 0)

	for itemK := range items {
		keyValue := items[itemK][keyName]
		delete(items[itemK], keyName)
		m.regexpItems[keyValue] = items[itemK]

		newKeyValue := regexp.QuoteMeta(keyValue)
		for _, keyFormatFun := range m.keyFormatFunList {
			newKeyValue = keyFormatFun(newKeyValue)
		}

		if newKeyValue == keyValue {
			regexpNormalItems = append(regexpNormalItems, newKeyValue)
		} else {
			keyGroupName := m.correctBadKeyOfGroupName(keyValue, itemK)
			regexpGroupItems = append(regexpGroupItems, fmt.Sprintf(`(?P<%s>%s)`, keyGroupName, newKeyValue))
		}
	}

	if len(regexpNormalItems) > 0 {
		regexpGroupItems = append([]string{fmt.Sprintf("(?P<%s>%s)", m.normalRegexpName, strings.Join(regexpNormalItems, "|"))}, regexpGroupItems...)
	}

	m.regexpString = strings.Join(regexpGroupItems, "|")
	m.regexp = regexp.MustCompile(m.regexpString)
	m.regexp.Longest()
}

func (m *Matcher) Match(contents []string) []*Result {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	return m.matchesToResults(matches)
}

func (m *Matcher) MatchText(contents []string) []*Result {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	results := make([]*Result, 0)
	resultIndex := make(map[string]int)

	for _, match := range matches {
		key := match[0]
		text := match[1]

		if index, ok := resultIndex[text]; ok {
			results[index].Texts[text] = 1
			results[index].Num += 1
		} else {
			results = append(results, m.matchToResult(match, false))
			resultIndex[key] = len(results) - 1
		}
	}

	return results
}

func (m *Matcher) MatchKey(contents []string) []*Result {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	results := make([]*Result, 0)
	resultIndex := make(map[string]int)

	for _, match := range matches {
		key := match[0]
		text := match[1]

		if index, ok := resultIndex[key]; ok {
			if _, ok := results[index].Texts[text]; ok {
				results[index].Texts[text] += 1
			} else {
				results[index].Texts[text] = 1
			}

			results[index].Num += 1
		} else {
			results = append(results, m.matchToResult(match, true))
			resultIndex[key] = len(results) - 1
		}
	}

	return results
}

func (m *Matcher) MatchFirstText(contents []string) *Result {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	return m.matchToResult(matches[0], false)
}

func (m *Matcher) MatchLastText(contents []string) *Result {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	return m.matchToResult(matches[len(matches)-1], false)
}

func (m *Matcher) MatchMostKey(contents []string) *Result {
	results := m.MatchKey(contents)
	if results == nil {
		return nil
	}

	sort.Sort(sortResults(results))

	return results[0]
}

func (m *Matcher) MatchMostText(contents []string) *Result {
	results := m.MatchText(contents)
	if results == nil {
		return nil
	}

	sort.Sort(sortResults(results))

	return results[0]
}

func (m *Matcher) MatchTag(contents []string) []*TagResult {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	results := make([]*TagResult, 0)
	resultIndex := make(map[string]int)

	for _, match := range matches {
		key := match[0]
		text := match[1]

		tags := m.regexpItems[key]

		tagsContent := ""
		for _, tag := range m.tagsName {
			tagsContent = tagsContent + "-" + tags[tag]
		}

		tagsIdentity := fmt.Sprintf("%x", md5.Sum([]byte(tagsContent)))

		if index, ok := resultIndex[tagsIdentity]; ok {
			result := results[index]
			if _, ok := result.Match[key]; ok {
				if _, ok := result.Match[key][text]; ok {
					result.Match[key][text] += 1
				} else {
					result.Match[key][text] = 1
				}
			} else {
				result.Match[key] = map[string]int{text: 1}
				result.KeyNum += 1
			}

			result.TextNum += 1
		} else {
			result := NewTagResult()
			result.Tags = tags
			result.Match[key] = map[string]int{text: 1}
			result.KeyNum += 1
			result.TextNum += 1

			results = append(results, result)
			resultIndex[tagsIdentity] = len(results) - 1
		}
	}

	return results
}

func (m *Matcher) MatchTagMostText(contents []string) *TagResult {
	results := m.MatchTag(contents)
	if results == nil {
		return nil
	}

	sort.Sort(sortTagResults(results))

	return results[0]
}

func (m *Matcher) addKeyFormatFun(f func(string) string) {
	m.keyFormatFunList = append(m.keyFormatFunList, f)
}

func (m *Matcher) correctBadKeyOfGroupName(key string, keyIndex int) string {
	newKey := fmt.Sprintf("%s%d", m.normalRegexpName, keyIndex)
	m.badKeyMap[newKey] = key

	return newKey
}

func (m *Matcher) matchesToResults(matches [][]string) []*Result {
	results := make([]*Result, 0, len(matches))
	for k := range matches {
		results = append(results, m.matchToResult(matches[k], false))
	}

	return results
}

func (m *Matcher) matchToResult(match []string, isTextComplex bool) *Result {
	mRes := NewResult()
	mRes.Key = match[0]
	mRes.Texts[match[1]] = 1
	mRes.Num = 1
	mRes.Tags = m.regexpItems[mRes.Key]
	mRes.IsTextComplex = isTextComplex

	return mRes
}

func (m *Matcher) findAllMatch(contents []string) [][]string {
	results := make([][]string, 0)

	for _, content := range contents {
		res := m.findAllSubMatch(content)
		if res != nil {
			results = append(results, res...)
		}
	}

	return results
}

func (m *Matcher) findAllSubMatch(content string) [][]string {
	allSubGroup := m.regexp.SubexpNames()
	allSubMatch := m.regexp.FindAllStringSubmatch(content, -1)

	matches := make([][]string, 0, len(allSubMatch))
	for _, subMatch := range allSubMatch {
		for k, text := range subMatch {
			if text == "" || k == 0 {
				continue
			}

			group := allSubGroup[k]

			if group == m.normalRegexpName {
				group = text
			} else {
				if key, ok := m.badKeyMap[group]; ok {
					group = key
				}
			}

			matches = append(matches, []string{group, text})

			break
		}
	}

	if len(matches) <= 0 {
		return nil
	}

	return matches
}

//
type sortResults []*Result

func (sr sortResults) Len() int {
	return len(sr)
}

func (sr sortResults) Less(i, j int) bool {
	return sr[i].Num > sr[j].Num
}

func (sr sortResults) Swap(i, j int) {
	sr[i], sr[j] = sr[j], sr[i]
}

//
type sortTagResults []*TagResult

func (str sortTagResults) Len() int {
	return len(str)
}

func (str sortTagResults) Less(i, j int) bool {
	return str[i].TextNum > str[j].TextNum
}

func (str sortTagResults) Swap(i, j int) {
	str[i], str[j] = str[j], str[i]
}
