package tag

import (
	"fmt"
	"maps"
	"regexp"
	"sort"
	"strings"
)

type matchedText struct {
	group string
	text  string
	index int // 多个 content， content 的序号
	start int // 包含 [
	stop  int // 不包含 (
}

// MatcherOption
// tag match option
type MatcherOption func(mt *Matcher)

// MatcherKeyFormatFunc
// 匹配前格式化 keyword 的 func list
type MatcherKeyFormatFunc func(string) string

func WithMatcherKeyFormatFunc(f ...MatcherKeyFormatFunc) MatcherOption {
	return func(m *Matcher) {
		m.addKeyFormatFunc(f...)
	}
}

func defaultMatcherKeyFormatFunc(s string) string {
	s = strings.TrimSpace(s)
	res, err := regexp.MatchString(`^[\w+._\s()]+$`, s)
	if err != nil {
		return s
	}

	if res {
		return fmt.Sprintf(`\b%s\b`, s)
	} else {
		return strings.ReplaceAll(s, "_", `.{1,3}`)
	}
}

func DefaultMatcher() *Matcher {
	return NewMatcher(WithMatcherKeyFormatFunc(defaultMatcherKeyFormatFunc))
}

// Matcher
// 从 rule 条目生成 regexp，匹配 content, 得到 keyword, matched text
//
// key：keyword
// text：被匹配到的 text
// label：label
// tag：name +label
type Matcher struct {
	fixed            map[string]string
	keyFormatFunc    []MatcherKeyFormatFunc // 在匹配前格式化关键词（使匹配更精确、丰富）
	keysIndex        map[string]int
	regexpItems      map[string]map[string]string // 关键词规则列表 map[关键词]map[标签名][标签值]
	regexp           *regexp.Regexp               // 所有关键词的 regexp
	regexpString     string                       // regular expression "(<?P<group name of keyword>...)"
	allSubGroupsName []string

	// 普通匹配：不包含 regular expression（纯文本）
	// 非普通匹配：包含 regular expression（需要指定 group name 与 keyword 关联）
	normalRegexpName string            // 普通匹配、非普通匹配的分组名称前缀（防止和自定义名称冲突，或 group 不支持的特殊字符）
	badKeyMap        map[string]string // 非普通匹配分组名称
	tagsName         []string          // 标签的名称
	hasItems         bool              // 是否有 items
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

// prepare
// keyName keyword
// items map[keyword, tags]
func (m *Matcher) prepare(keyName string, items []map[string]string, fixed map[string]string) {
	if len(items) <= 0 {
		return
	}

	m.hasItems = true
	m.fixed = fixed
	m.regexpItems = make(map[string]map[string]string, len(items))
	m.keysIndex = make(map[string]int, len(items))

	for k := range items[0] {
		if k != keyName {
			m.tagsName = append(m.tagsName, k)
		}
	}

	sort.SliceStable(m.tagsName, func(i, j int) bool {
		return m.tagsName[i] < m.tagsName[j]
	})

	// 普通匹配和非普通匹配的表达式（英文、数字等需要通过前后限定符分组精确匹配）
	var normalItems []string
	var groupRegexps []string

	for itemK := range items {
		keyValue := items[itemK][keyName]
		delete(items[itemK], keyName)
		m.regexpItems[keyValue] = items[itemK]

		m.keysIndex[keyValue] = itemK

		newKeyValue := regexp.QuoteMeta(keyValue)
		for _, keyFormatFun := range m.keyFormatFunc {
			newKeyValue = keyFormatFun(newKeyValue)
		}

		if newKeyValue == keyValue { // 普通匹配
			normalItems = append(normalItems, newKeyValue)
		} else { // 非普通匹配
			keyGroupName := m.correctBadKeyOfGroupName(keyValue, itemK)
			groupRegexps = append(groupRegexps, fmt.Sprintf(`(?P<%s>%s)`, keyGroupName, newKeyValue))
		}
	}

	if len(normalItems) > 0 {
		groupRegexps = append([]string{fmt.Sprintf("(?P<%s>%s)", m.normalRegexpName, strings.Join(normalItems, "|"))}, groupRegexps...)
	}

	m.regexpString = strings.Join(groupRegexps, "|")
	m.regexp = regexp.MustCompile(m.regexpString)
	m.regexp.Longest()
	m.allSubGroupsName = m.regexp.SubexpNames()
}

// Match
// all matched, in regexp match order, key order
func (m *Matcher) Match(contents []string) Results {
	matches := m.findMatchAll(contents)
	if matches == nil {
		return nil
	}

	return m.matchesToResults(matches)
}

// MatchInKeyOrder
// all matched
// the order of matched keyword
func (m *Matcher) MatchInKeyOrder(contents []string) Results {
	matches := m.findMatchAllInKeyOrder(contents)
	if matches == nil {
		return nil
	}

	return m.matchesToResults(matches)
}

// MatchText
// match text 合并相同的 matched text
func (m *Matcher) MatchText(contents []string) Results {
	matches := m.findMatchAll(contents)
	if matches == nil {
		return nil
	}

	var results Results
	resultIndex := make(map[string]int)

	for _, matchText := range matches {
		text := matchText.text

		if index, ok := resultIndex[text]; ok {
			results[index].Texts[text] += 1
			results[index].Amount += 1
		} else {
			results = append(results, m.matchToResult(matchText))
			resultIndex[text] = len(results) - 1
		}
	}

	return results
}

// MatchFirstText
// the leftmost matched text
func (m *Matcher) MatchFirstText(contents []string) Results {
	matches := m.findMatchAll(contents)
	if matches == nil {
		return nil
	}

	return m.matchToResults(matches[0])
}

// MatchLastText
// the rightmost matched text
func (m *Matcher) MatchLastText(contents []string) Results {
	matches := m.findMatchAll(contents)
	if matches == nil {
		return nil
	}

	return m.matchToResults(matches[len(matches)-1])
}

// MatchMostText
// the text that has been matched the most times
func (m *Matcher) MatchMostText(contents []string) Results {
	results := m.MatchText(contents)
	if results == nil {
		return nil
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Amount > results[j].Amount
	})

	return results[0:1]
}

// MatchKey
// match key 合并相同的 keyword（同时也合并 matched text）
// in matched key order
func (m *Matcher) MatchKey(contents []string) Results {
	matches := m.findMatchAllInKeyOrder(contents)
	if matches == nil {
		return nil
	}

	var results Results
	resultIndex := make(map[string]int)

	for _, matchText := range matches {
		key := matchText.group
		text := matchText.text

		if index, ok := resultIndex[key]; ok {
			results[index].Texts[text] += 1
			results[index].Amount += 1
		} else {
			results = append(results, m.matchToResult(matchText))
			resultIndex[key] = len(results) - 1
		}
	}

	return results
}

// MatchFirstKey
// the first matched key
func (m *Matcher) MatchFirstKey(contents []string) Results {
	matches := m.findMatchAllInKeyOrder(contents)
	if matches == nil {
		return nil
	}

	return m.matchToResults(matches[0])
}

// MatchLastKey
// the last matched key
func (m *Matcher) MatchLastKey(contents []string) Results {
	matches := m.findMatchAllInKeyOrder(contents)
	if matches == nil {
		return nil
	}

	return m.matchToResults(matches[len(matches)-1])
}

// MatchMostKey
// match most key 被匹配次数最多的 keyword
func (m *Matcher) MatchMostKey(contents []string) Results {
	results := m.MatchKey(contents)
	if results == nil {
		return nil
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Amount > results[j].Amount
	})

	return results[0:1]
}

// MatchLabel
// match label 合并重复的 tags 组合
func (m *Matcher) MatchLabel(contents []string) LabelResults {
	matches := m.findMatchAll(contents)
	if matches == nil {
		return nil
	}

	var results LabelResults
	resultIndex := make(map[string]int)

	for _, matchText := range matches {
		key := matchText.group
		text := matchText.text

		tags := m.regexpItems[key]

		tagsIdentity := ""
		for _, tag := range m.tagsName {
			tagsIdentity = tagsIdentity + "-" + tags[tag]
		}

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
				result.Keywords = append(result.Keywords, key)
			}

			result.Amount += 1
			results[index] = result
		} else {
			result := NewLabelResult()
			result.Identity = tagsIdentity
			maps.Copy(result.Tags, tags)
			maps.Copy(result.Tags, m.fixed)
			result.Match[key] = map[string]int{text: 1}
			result.Keywords = append(result.Keywords, key)
			result.Amount += 1

			results = append(results, result)
			resultIndex[tagsIdentity] = len(results) - 1
		}
	}

	return results
}

// MatchFirstLabel
// first label
func (m *Matcher) MatchFirstLabel(contents []string) LabelResults {
	results := m.MatchLabel(contents)
	if results == nil {
		return nil
	}

	return results[0:1]
}

// MatchLabelMostText
// match label most text 合并重复的 tags 组合中，text 最多次数
func (m *Matcher) MatchLabelMostText(contents []string) LabelResults {
	results := m.MatchLabel(contents)
	if results == nil {
		return nil
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Amount > results[j].Amount
	})

	return results[0:1]
}

func (m *Matcher) addKeyFormatFunc(f ...MatcherKeyFormatFunc) {
	m.keyFormatFunc = append(m.keyFormatFunc, f...)
}

// correctBadKeyOfGroupName
// 避免不合法的分组名称
func (m *Matcher) correctBadKeyOfGroupName(key string, keyIndex int) string {
	newKey := fmt.Sprintf("%s%d", m.normalRegexpName, keyIndex)
	m.badKeyMap[newKey] = key

	return newKey
}

func (m *Matcher) matchesToResults(matches []matchedText) Results {
	results := make(Results, 0, len(matches))
	for k := range matches {
		results = append(results, m.matchToResult(matches[k]))
	}

	return results
}

func (m *Matcher) matchToResults(match matchedText) Results {
	return Results{m.matchToResult(match)}
}

func (m *Matcher) matchToResult(match matchedText) Result {
	r := NewResult()
	r.Keyword = match.group
	r.Texts[match.text] = 1
	r.Amount = 1
	maps.Copy(r.Tags, m.regexpItems[r.Keyword])
	maps.Copy(r.Tags, m.fixed)

	return r
}

// findMatchAllInKeyOrder
// all regexp match, in matched keyword order
func (m *Matcher) findMatchAllInKeyOrder(contents []string) []matchedText {
	matches := m.findMatchAll(contents)
	if matches == nil {
		return nil
	}

	sort.SliceStable(matches, func(i, j int) bool {
		return m.keysIndex[matches[i].group] < m.keysIndex[matches[j].group]
	})

	return matches
}

// findMatchAll
// all regexp match, is regexp match order, in matched text order
// the leftmost text is at the front
func (m *Matcher) findMatchAll(contents []string) []matchedText {
	rets := make([]matchedText, 0)

	for i, content := range contents {
		ret := m.findAllSubMatch(i, content, -1)
		if ret != nil {
			rets = append(rets, ret...)
		}
	}

	if len(rets) <= 0 {
		return nil
	}

	return rets
}

// findMatchFirst
// first regexp match
func (m *Matcher) findMatchFirst(contents []string) []matchedText {
	var rets []matchedText

	for i, content := range contents {
		rets = m.findAllSubMatch(i, content, 1)
		if rets != nil {
			break
		}
	}

	return rets
}

func (m *Matcher) findAllSubMatch(index int, content string, n int) []matchedText {
	if !m.hasItems {
		return nil
	}

	// 所有分组的匹配结果
	allSubMatchIndex := m.regexp.FindAllStringSubmatchIndex(content, n)

	matches := make([]matchedText, 0, len(allSubMatchIndex))
	for _, subMatch := range allSubMatchIndex {
		_mt := matchedText{
			index: index,
			start: subMatch[0],
			stop:  subMatch[1],
		}

		smLen := len(subMatch)
		for i := 2; i < smLen; i += 2 {
			if subMatch[i] != -1 {
				_mt.text = content[subMatch[i]:subMatch[i+1]]
				_mt.group = m.getGroupName(i/2, _mt.text)

				break
			}
		}

		matches = append(matches, _mt)
	}

	if len(matches) <= 0 {
		return nil
	}

	return matches
}

func (m *Matcher) getGroupName(groupIndex int, text string) string {
	group := m.allSubGroupsName[groupIndex]

	if group == m.normalRegexpName {
		group = text
	} else {
		if key, ok := m.badKeyMap[group]; ok {
			group = key
		}
	}

	return group
}

// findAllSubMatch
// [][keyword, matched text]
func (m *Matcher) findAllSubMatchBackup(content string, n int) [][]string {
	if !m.hasItems {
		return nil
	}

	// 所有分组的匹配结果
	allSubMatch := m.regexp.FindAllStringSubmatch(content, n)

	matches := make([][]string, 0, len(allSubMatch))
	for _, subMatch := range allSubMatch {
		for k, text := range subMatch {
			if text == "" || k == 0 {
				continue
			}

			group := m.getGroupName(k, text)
			matches = append(matches, []string{group, text})

			break
		}
	}

	if len(matches) <= 0 {
		return nil
	}

	return matches
}
