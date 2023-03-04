package tagor

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"sort"
	"strings"

	goetl "github.com/auho/go-etl"
	go_simple_db "github.com/auho/go-simple-db/v2"
)

// #TODO 重命名（key keyword  label=>tag）

// Result
// result 匹配结果
type Result struct {
	Key           string            // keyword
	Num           int64             // matched num
	Texts         map[string]int64  // matched text map[matched text]num
	Tags          map[string]string // tags map[tag name]tag
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
type LabelResult struct {
	Identity string
	Labels   map[string]string         // tags map[tag name]tag
	Match    map[string]map[string]int // keyword and match text map[keyword]map[matched text]num
	KeyNum   int64                     // keyword num
	TextNum  int64                     // all text word num
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
// table name
// - rule prefix + data name + rule key name
// - rule prefix + short table name
// - rule prefix + rule key name
//
type TagMatcher struct {
	Matcher           *Matcher
	db                *go_simple_db.SimpleDB
	tagsName          []string            // 关键词匹配的标签列表名称： [tagA1, tagA2]
	key               string              // 关键词名称： tagA
	keyFieldName      string              // 关键词字段名称：tagA_keyword
	keyNumFieldName   string              // 关键词出现数量：tagA_keyword_num
	tableName         string              // rule 表：rule_tagA
	tableTagFields    []string            // 数据表标签字段：[tagA1, tagA2, tagA_keyword]
	excludeTableField map[string]struct{} // 排除的数据表字段：[id, keyword_len]

	dataName       string            // data name
	shortTableName string            // short name of tag data table
	alias          map[string]string // 别名 [data name => output name]
	fixedTags      map[string]string // fixed tags data
	fixedKeys      []string          // keys of fixed data
	fixedValues    []interface{}     // values of fixed data
}

func newTagMatcher(key string, db *go_simple_db.SimpleDB, Options ...TagMatcherOption) *TagMatcher {
	t := &TagMatcher{}
	t.prepare(key, db, Options...)

	return t
}

func (t *TagMatcher) prepare(key string, db *go_simple_db.SimpleDB, Options ...TagMatcherOption) {
	t.key = key
	t.db = db
	t.excludeTableField = map[string]struct{}{"id": {}, "keyword_len": {}}

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
		t.tableName = goetl.RuleTableNamePrefix + "_" + t.shortTableName
	} else if t.dataName != "" {
		t.tableName = goetl.RuleTableNamePrefix + "_" + t.dataName + "_" + t.key
	} else {
		t.tableName = goetl.RuleTableNamePrefix + "_" + t.key
	}

	t.keyFieldName = t.key + "_keyword"
	t.keyNumFieldName = t.keyFieldName + "_num"

	row, err := t.db.GetTableColumns(t.tableName)
	if err != nil {
		panic(err)
	}

	t.tableTagFields = make([]string, 0)

	hasAssignTagsName := false
	if len(t.tagsName) > 0 {
		hasAssignTagsName = true
	} else {
		t.tagsName = make([]string, 0)
	}

	// 获取被匹配的标签字段列表
	for k := range row {
		column := row[k]
		if _, ok := t.excludeTableField[column]; ok {
			continue
		}

		t.tableTagFields = append(t.tableTagFields, column)

		if hasAssignTagsName == false {
			if column != t.keyFieldName {
				t.tagsName = append(t.tagsName, column)
			}
		}
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

	rules := make([]map[string]interface{}, 0)
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY keyword_len DESC, id ASC", strings.Join(columns, ", "), t.tableName)
	err := t.db.Raw(query).Scan(&rules).Error
	if err != nil {
		panic(err)
	}

	if len(rules) <= 0 {
		panic("rules is null")
	}

	_rules := make([]map[string]string, 0)
	for _, rule := range rules {
		_rule := make(map[string]string)
		for k, v := range rule {
			if _v, ok := v.(string); ok {
				_rule[k] = _v
			} else {
				panic(fmt.Sprintf("TagMatcher type of v is not string %s => %v", k, v))
			}
		}

		_rules = append(_rules, _rule)
	}

	return _rules
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

func (t *TagMatcher) resultToSliceMap(result *Result) []map[string]interface{} {
	return []map[string]interface{}{t.resultToMap(result)}
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
	keyFormatFunList []func(string) string        // 在匹配前处理关键词（使匹配更精确、丰富）
	regexpItems      map[string]map[string]string // 关键词规则列表 map[关键词]map[标签名称][标签]
	regexp           *regexp.Regexp               // 所有关键词的 regexp
	regexpString     string                       // regular expression
	normalRegexpName string                       // 普通匹配的分组名称
	badKeyMap        map[string]string            // 非普通匹配分组名称
	tagsName         []string                     // 标签名称列表
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

	// 普通匹配和非普通匹配的表达式（英文、数字等需要通过前后限定符分组精确匹配）
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

// Match
// 重复结果不合并
func (m *Matcher) Match(contents []string) []*Result {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	return m.matchesToResults(matches)
}

// MatchText
// match text 合并相同的 matched text
func (m *Matcher) MatchText(contents []string) []*Result {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	results := make([]*Result, 0)
	resultIndex := make(map[string]int)

	for _, match := range matches {
		text := match[1]

		if index, ok := resultIndex[text]; ok {
			results[index].Texts[text] = 1
			results[index].Num += 1
		} else {
			results = append(results, m.matchToResult(match, false))
			resultIndex[text] = len(results) - 1
		}
	}

	return results
}

// MatchKey
// match key 合并相同的 keyword（同时也合并 matched text）
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
			if _, ok1 := results[index].Texts[text]; ok1 {
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

// MatchFirstText
// match first text 匹配第一个被匹配的 text
func (m *Matcher) MatchFirstText(contents []string) *Result {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	return m.matchToResult(matches[0], false)
}

// MatchLastText
// match first text 最后一个被匹配的 text
func (m *Matcher) MatchLastText(contents []string) *Result {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	return m.matchToResult(matches[len(matches)-1], false)
}

// MatchMostKey
// match most key 被匹配次数最多的 keyword
func (m *Matcher) MatchMostKey(contents []string) *Result {
	results := m.MatchKey(contents)
	if results == nil {
		return nil
	}

	sort.Sort(sortResults(results))

	return results[0]
}

// MatchMostText
// match most text 被匹配次数最多的 matched text
func (m *Matcher) MatchMostText(contents []string) *Result {
	results := m.MatchText(contents)
	if results == nil {
		return nil
	}

	sort.Sort(sortResults(results))

	return results[0]
}

// MatchLabel
// match label 合并重复的 tags 组合
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

		// # TODO delete
		tagsIdentity := fmt.Sprintf("%x", md5.Sum([]byte(tagsContent)))

		if index, ok := resultIndex[tagsIdentity]; ok {
			result := results[index]
			if _, ok1 := result.Match[key]; ok1 {
				if _, ok2 := result.Match[key][text]; ok2 {
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

// MatchLabelMostText
// match label most text 合并重复的 tags 组合中，text 最多次数
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

// correctBadKeyOfGroupName
// 避免不合法的分组名称
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

// findAllMatch
// [][keyword, matched text]
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

// findAllSubMatch
// [][keyword, matched text]
func (m *Matcher) findAllSubMatch(content string) [][]string {
	// 所有分组的匹配结果
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
