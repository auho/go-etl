package tagor

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"sort"
	"strings"

	goEtl "github.com/auho/go-etl"
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

// LabelResult
// label result
//
type LabelResult struct {
	Identity string
	Labels   map[string]string
	Match    map[string]map[string]int
	KeyNum   int64
	TextNum  int64
}

func NewLabelResult() *LabelResult {
	l := &LabelResult{}
	l.Labels = make(map[string]string)
	l.Match = make(map[string]map[string]int)

	return l
}

type TagMatcherOption func(*TagMatcher)

func WithTagMatcherTags(s []string) TagMatcherOption {
	return func(t *TagMatcher) {
		t.tagsName = s
	}
}

func WithTagMatcherAlias(alias map[string]string) TagMatcherOption {
	return func(t *TagMatcher) {
		t.alias = alias
	}
}

func WithTagMatcherDataName(d string) TagMatcherOption {
	return func(t *TagMatcher) {
		t.dataName = d
	}
}

func WithTagMatcherShortTableName(s string) TagMatcherOption {
	return func(t *TagMatcher) {
		t.shortTableName = s
	}
}

func WithTagMatcherFixedTags(m map[string]string) TagMatcherOption {
	return func(t *TagMatcher) {
		t.fixedTags = m
	}
}

func WithTagMatcherMatcher(options []MatcherOption) TagMatcherOption {
	return func(t *TagMatcher) {
		t.Matcher = NewMatcher(options...)
	}
}

// TagMatcher
// tag matcher
//
type TagMatcher struct {
	Matcher           *Matcher
	db                simple.Driver
	tagsName          []string
	key               string
	keyFieldName      string
	keyNumFieldName   string
	tableName         string
	tableTagFields    []string
	excludeTableField []string

	dataName       string
	shortTableName string
	alias          map[string]string
	fixedTags      map[string]string
	fixedKeys      []string
	fixedValues    []interface{}
}

func newTagMatcher(key string, db simple.Driver, Options ...TagMatcherOption) *TagMatcher {
	t := &TagMatcher{}
	t.prepare(key, db, Options...)

	return t
}

func (t *TagMatcher) prepare(key string, db simple.Driver, Options ...TagMatcherOption) {
	t.key = key
	t.db = db
	t.excludeTableField = []string{"id", "keyword_len"}

	for _, option := range Options {
		option(t)
	}

	if t.Matcher == nil {
		t.Matcher = NewMatcher(WithTagMatcherKeyFun(func(s string) string {
			res, err := regexp.MatchString(`^[\w+._\s()]+$`, s)
			if err != nil {
				return s
			}

			if res {
				return fmt.Sprintf(`\b%s\b`, s)
			} else {
				return strings.ReplaceAll(s, "_", `.{1,3}`)
			}
		}))
	}

	t.init()
	t.Matcher.init(t.keyFieldName, t.getRules())
}

func (t *TagMatcher) init() {
	if t.shortTableName != "" {
		t.tableName = goEtl.RuleTableNamePrefix + "_" + t.shortTableName
	} else if t.dataName != "" {
		t.tableName = goEtl.RuleTableNamePrefix + "_" + t.dataName + "_" + t.key
	} else {
		t.tableName = goEtl.RuleTableNamePrefix + "_" + t.key
	}

	t.keyFieldName = t.key + "_keyword"
	t.keyNumFieldName = t.keyFieldName + "_num"

	row, err := t.db.GetTableColumns(t.tableName)
	if err != nil {
		panic(err)
	}

	t.tableTagFields = make([]string, 0)

	isTagsName := false
	if len(t.tagsName) > 0 {
		isTagsName = true
	} else {
		t.tagsName = make([]string, 0)
	}

	for k := range row {
		column := row[k]
		for _, ec := range t.excludeTableField {
			if ec == column {
				goto LOOP
			}
		}

		t.tableTagFields = append(t.tableTagFields, column)

		if isTagsName == false {
			if column != t.keyFieldName {
				t.tagsName = append(t.tagsName, column)
			}
		}
	LOOP:
	}

	for k, v := range t.tagsName {
		if s, ok := t.alias[v]; ok {
			t.tagsName[k] = s
		}
	}

	if s, ok := t.alias[t.keyFieldName]; ok {
		t.keyFieldName = s
	}

	if s, ok := t.alias[t.keyNumFieldName]; ok {
		t.keyNumFieldName = s
	}

	t.fixedKeys = make([]string, 0)
	t.fixedValues = make([]interface{}, 0)
	if len(t.fixedTags) > 0 {
		for k := range t.fixedTags {
			if s, ok := t.alias[k]; ok {
				k = s
			}

			t.fixedKeys = append(t.fixedKeys, k)
			t.fixedValues = append(t.fixedValues, t.fixedTags[k])
		}
	}
}

func (t *TagMatcher) getRules() []map[string]string {
	columns := make([]string, 0)
	for _, f := range t.tableTagFields {
		if s, ok := t.alias[f]; ok {
			f = fmt.Sprintf("`%s` AS '%s'", f, s)
		} else {
			f = fmt.Sprintf("`%s`", f)
		}

		columns = append(columns, f)
	}

	query := fmt.Sprintf("SELECT %s FROM `%s` ORDER BY `keyword_len` DESC, `id` ASC", strings.Join(columns, ", "), t.tableName)
	rules, err := t.db.QueryString(query)
	if err != nil {
		panic(err)
	}

	if len(rules) <= 0 {
		panic("rules is null")
	}

	return rules
}

func (t *TagMatcher) getResultInsertKeys() []string {
	return append([]string{t.keyFieldName, t.keyNumFieldName}, append(t.tagsName, t.fixedKeys...)...)
}

func (t *TagMatcher) resultsToSliceMap(results []*Result) []map[string]interface{} {
	items := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		items = append(items, t.resultToMap(result))
	}

	return items
}

func (t *TagMatcher) resultToMap(result *Result) map[string]interface{} {
	item := make(map[string]interface{})
	item[t.keyFieldName] = result.Key
	item[t.keyNumFieldName] = result.Num

	for _, tagName := range t.tagsName {
		item[tagName] = result.Tags[tagName]
	}

	if len(t.fixedTags) > 0 {
		for k := range t.fixedTags {
			item[k] = t.fixedTags[k]
		}
	}

	return item
}

func (t *TagMatcher) resultToSlice(result *Result) []interface{} {
	item := make([]interface{}, 0, len(t.tagsName)+2)
	item = append(item, result.Key)
	item = append(item, result.Num)

	for _, tagName := range t.tagsName {
		item = append(item, result.Tags[tagName])
	}

	if len(t.fixedValues) > 0 {
		item = append(item, t.fixedValues...)
	}

	return item
}

func (t *TagMatcher) resultToSliceSlice(result *Result) [][]interface{} {
	item := t.resultToSlice(result)

	return [][]interface{}{item}
}

func (t *TagMatcher) resultsToSliceSlice(results []*Result) [][]interface{} {
	items := make([][]interface{}, 0, len(results))
	for _, result := range results {
		items = append(items, t.resultToSlice(result))
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

func NewMatcher(Options ...MatcherOption) *Matcher {
	m := &Matcher{}
	m.normalRegexpName = "_rEgEx_"
	m.badKeyMap = make(map[string]string)

	for _, option := range Options {
		option(m)
	}

	return m
}

func (m *Matcher) init(keyName string, items []map[string]string) {
	m.regexpItems = make(map[string]map[string]string, len(items))

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

func (m *Matcher) MatchLabel(contents []string) []*LabelResult {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	results := make([]*LabelResult, 0)
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
			result := NewLabelResult()
			result.Labels = tags
			result.Match[key] = map[string]int{text: 1}
			result.KeyNum += 1
			result.TextNum += 1

			results = append(results, result)
			resultIndex[tagsIdentity] = len(results) - 1
		}
	}

	return results
}

func (m *Matcher) MatchLabelMostText(contents []string) *LabelResult {
	results := m.MatchLabel(contents)
	if results == nil {
		return nil
	}

	sort.Sort(sortLabelResults(results))

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

	if len(results) <= 0 {
		return nil
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
type sortLabelResults []*LabelResult

func (str sortLabelResults) Len() int {
	return len(str)
}

func (str sortLabelResults) Less(i, j int) bool {
	return str[i].TextNum > str[j].TextNum
}

func (str sortLabelResults) Swap(i, j int) {
	str[i], str[j] = str[j], str[i]
}
