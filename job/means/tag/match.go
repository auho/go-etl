package tag

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

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
	keyFormatFunc []MatcherKeyFormatFunc       // 在匹配前格式化关键词（使匹配更精确、丰富）
	regexpItems   map[string]map[string]string // 关键词规则列表 map[关键词]map[标签名称][标签]
	regexp        *regexp.Regexp               // 所有关键词的 regexp
	regexpString  string                       // regular expression "(<?P<group name of keyword>...)"

	// 普通匹配：不包含 regular expression（纯文本）
	// 非普通匹配：包含 regular expression（需要指定 group name 与 keyword 关联）
	normalRegexpName string            // 普通匹配、非普通匹配的分组名称前缀（防止和自定义名称冲突，或 group 不支持的特殊字符）
	badKeyMap        map[string]string // 非普通匹配分组名称
	tagsName         []string          // 标签名称列表
	hasItems         bool
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
func (m *Matcher) prepare(keyName string, items []map[string]string) {
	if len(items) <= 0 {
		return
	}

	m.hasItems = true

	m.regexpItems = make(map[string]map[string]string, len(items))

	for k := range items[0] {
		if k != keyName {
			m.tagsName = append(m.tagsName, k)
		}
	}

	// 普通匹配和非普通匹配的表达式（英文、数字等需要通过前后限定符分组精确匹配）
	var normalItems []string
	var groupRegexps []string

	for itemK := range items {
		keyValue := items[itemK][keyName]
		delete(items[itemK], keyName)
		m.regexpItems[keyValue] = items[itemK]

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
}

// Match
// 匹配所有结果，重复结果不合并
func (m *Matcher) Match(contents []string) Results {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	return m.matchesToResults(matches)
}

// MatchText
// match text 合并相同的 matched text
func (m *Matcher) MatchText(contents []string) Results {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	var results Results
	resultIndex := make(map[string]int)

	for _, match := range matches {
		text := match[1]

		if index, ok := resultIndex[text]; ok {
			results[index].Texts[text] += 1
			results[index].Amount += 1
		} else {
			results = append(results, m.matchToResult(match, false))
			resultIndex[text] = len(results) - 1
		}
	}

	return results
}

// MatchKey
// match key 合并相同的 keyword（同时也合并 matched text）
func (m *Matcher) MatchKey(contents []string) Results {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	var results Results
	resultIndex := make(map[string]int)

	for _, match := range matches {
		key := match[0]
		text := match[1]

		if index, ok := resultIndex[key]; ok {
			results[index].Texts[text] += 1
			results[index].Amount += 1
		} else {
			results = append(results, m.matchToResult(match, true))
			resultIndex[key] = len(results) - 1
		}
	}

	return results
}

// MatchFirstText
// match first text 匹配第一个被匹配的 text
func (m *Matcher) MatchFirstText(contents []string) Results {
	matches := m.findFirstMatch(contents)
	if matches == nil {
		return nil
	}

	return m.matchToResults(matches[0], false)
}

// MatchLastText
// match first text 最后一个被匹配的 text
func (m *Matcher) MatchLastText(contents []string) Results {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	return m.matchToResults(matches[len(matches)-1], false)
}

// MatchMostKey
// match most key 被匹配次数最多的 keyword
func (m *Matcher) MatchMostKey(contents []string) Results {
	results := m.MatchKey(contents)
	if results == nil {
		return nil
	}

	sort.Sort(results)

	return results[0:1]
}

// MatchMostText
// match most text 被匹配次数最多的 matched text
func (m *Matcher) MatchMostText(contents []string) Results {
	results := m.MatchText(contents)
	if results == nil {
		return nil
	}

	sort.Sort(results)

	return results[0:1]
}

// MatchLabel
// match label 合并重复的 tags 组合
func (m *Matcher) MatchLabel(contents []string) LabelResults {
	matches := m.findAllMatch(contents)
	if matches == nil {
		return nil
	}

	var results LabelResults
	resultIndex := make(map[string]int)

	for _, match := range matches {
		key := match[0]
		text := match[1]

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
				result.Keys = append(result.Keys, key)
			}

			result.Amount += 1
			results[index] = result
		} else {
			result := NewLabelResult()
			result.Identity = tagsIdentity
			result.Tags = tags
			result.Match[key] = map[string]int{text: 1}
			result.Keys = append(result.Keys, key)
			result.Amount += 1

			results = append(results, result)
			resultIndex[tagsIdentity] = len(results) - 1
		}
	}

	return results
}

// MatchLabelMostText
// match label most text 合并重复的 tags 组合中，text 最多次数
func (m *Matcher) MatchLabelMostText(contents []string) LabelResults {
	results := m.MatchLabel(contents)
	if results == nil {
		return nil
	}

	sort.Sort(results)

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

func (m *Matcher) matchesToResults(matches [][]string) Results {
	results := make(Results, 0, len(matches))
	for k := range matches {
		results = append(results, m.matchToResult(matches[k], false))
	}

	return results
}

func (m *Matcher) matchToResults(match []string, isKeyMerge bool) Results {
	return Results{m.matchToResult(match, isKeyMerge)}
}

func (m *Matcher) matchToResult(match []string, isKeyMerge bool) Result {
	mRes := NewResult()
	mRes.Keyword = match[0]
	mRes.Texts[match[1]] = 1
	mRes.Amount = 1
	mRes.Tags = m.regexpItems[mRes.Keyword]
	mRes.IsKeyMerge = isKeyMerge

	return mRes
}

// findAllMatch
// [][keyword, matched text]
func (m *Matcher) findAllMatch(contents []string) [][]string {
	results := make([][]string, 0)

	for _, content := range contents {
		res := m.findAllSubMatch(content, -1)
		if res != nil {
			results = append(results, res...)
		}
	}

	if len(results) <= 0 {
		return nil
	}

	return results
}

// findFirstMatch
// [][keyword, matched text]
func (m *Matcher) findFirstMatch(contents []string) [][]string {
	var res [][]string

	for _, content := range contents {
		res = m.findAllSubMatch(content, 1)
		if res != nil {
			break
		}
	}

	return res
}

// findAllSubMatch
// [][keyword, matched text]
func (m *Matcher) findAllSubMatch(content string, n int) [][]string {
	if !m.hasItems {
		return nil
	}

	// 所有分组的匹配结果
	allSubGroup := m.regexp.SubexpNames()
	allSubMatch := m.regexp.FindAllStringSubmatch(content, n)

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
