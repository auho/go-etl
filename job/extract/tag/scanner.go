package tag

import (
	"fmt"
	"maps"
	"regexp"
	"sort"
	"strings"

	"github.com/auho/go-etl/v3/job/extract"
)

// ScannerConfig scanner config
type ScannerConfig struct {
	Debug bool
}

// scanned text
type scannedText struct {
	keyword string
	text    string
	index   int // 多个 content， content 的序号
	start   int // 包含 [
	stop    int // 不包含 (
}

// ScannerOption
// tag scanner option
type ScannerOption func(mt *scanner)

// ScannerKeyFormatter
// 匹配前格式化 keyword 的 func list
type ScannerKeyFormatter func(string) string

func WithScannerKeyFormatter(fs ...ScannerKeyFormatter) ScannerOption {
	return func(m *scanner) {
		m.addKeyFormatters(fs...)
	}
}

func defaultScannerKeyFormatter(s string) string {
	s = strings.TrimSpace(s)
	ok, err := regexp.MatchString(`^[\w+._\s()]+$`, s)
	if err != nil {
		return s
	}

	if ok {
		return fmt.Sprintf(`\b%s\b`, s)
	} else {
		return strings.ReplaceAll(s, "_", `.{1,3}`)
	}
}

// default scanner
func defaultScanner(rule extract.Rule, sc ScannerConfig) (*scanner, error) {
	return newScannerFromRule(rule, sc, WithScannerKeyFormatter(defaultScannerKeyFormatter))
}

func newScannerFromRule(rule extract.Rule, sc ScannerConfig, opts ...ScannerOption) (*scanner, error) {
	items, err := rule.ItemsForRegexp()
	if err != nil {
		return nil, fmt.Errorf("ItemsForRegexp: %w", err)
	}

	return newScanner(rule.KeywordNameAlias(), items, opts...), nil
}

// scanner
// 从 rule 条目生成 regexp，匹配 content, 得到 keyword, scanned text
//
// key：keyword
// text：被匹配到的 text
// label：label
// tag：name +label
type scanner struct {
	keyFormatters    []ScannerKeyFormatter // 在匹配前格式化关键词（使匹配更精确、丰富）
	keyIndex         map[string]int
	regexpItems      map[string]map[string]string // 关键词规则列表 map[关键词]map[标签名][标签值]
	regexp           *regexp.Regexp               // 所有关键词的 regexp
	regexpString     string                       // regular expression "(<?P<group name of keyword>...)"
	allSubGroupNames []string

	// 普通匹配：不包含 regular expression（纯文本）
	// 非普通匹配：包含 regular expression（需要指定 group name 与 keyword 关联）
	groupNamePrefix string            // 普通匹配、非普通匹配的分组名称前缀（防止和自定义名称冲突，或 group 不支持的特殊字符）
	groupNameMap    map[string]string // 非普通匹配分组名称
	tagNames        []string          // 标签的名称
	hasItems        bool              // 是否有 items
}

func newScanner(keyName string, items []map[string]string, opts ...ScannerOption) *scanner {
	m := &scanner{}
	m.groupNamePrefix = "_rEgEx_"
	m.groupNameMap = make(map[string]string)

	for _, option := range opts {
		option(m)
	}

	m.prepare(keyName, items)

	return m
}

// prepare
// keyName keyword
// items map[keyword, tags]
func (s *scanner) prepare(keyName string, items []map[string]string) {
	if len(items) == 0 {
		return
	}

	s.hasItems = true
	s.regexpItems = make(map[string]map[string]string, len(items))
	s.keyIndex = make(map[string]int, len(items))

	for k := range items[0] {
		if k != keyName {
			s.tagNames = append(s.tagNames, k)
		}
	}

	sort.SliceStable(s.tagNames, func(i, j int) bool {
		return s.tagNames[i] < s.tagNames[j]
	})

	// 普通匹配和非普通匹配的表达式（英文、数字等需要通过前后限定符分组精确匹配）
	var normalItems []string
	var groupRegexps []string

	for item := range items {
		keyValue := items[item][keyName]
		delete(items[item], keyName)
		s.regexpItems[keyValue] = items[item]

		s.keyIndex[keyValue] = item

		newKeyValue := regexp.QuoteMeta(keyValue)
		for _, kf := range s.keyFormatters {
			newKeyValue = kf(newKeyValue)
		}

		if newKeyValue == keyValue { // 普通匹配
			normalItems = append(normalItems, newKeyValue)
		} else { // 非普通匹配
			keyGroupName := s.correctBadKeyOfGroupName(keyValue, item)
			groupRegexps = append(groupRegexps, fmt.Sprintf(`(?P<%s>%s)`, keyGroupName, newKeyValue))
		}
	}

	if len(normalItems) > 0 {
		groupRegexps = append(groupRegexps, fmt.Sprintf("(?P<%s>%s)", s.groupNamePrefix, strings.Join(normalItems, "|")))
	}

	s.regexpString = strings.Join(groupRegexps, "|")
	s.regexp = regexp.MustCompile(s.regexpString)
	s.regexp.Longest()
	s.allSubGroupNames = s.regexp.SubexpNames()
}

// Scan
// all scanned
// in regexp scanner order, text order
func (s *scanner) Scan(contents []string) results {
	sts := s.findScanAll(contents)
	if sts == nil {
		return nil
	}

	return s.scansToResults(sts)
}

// ScanInKeyOrder
// all scanned
// in scanned keyword order
func (s *scanner) ScanInKeyOrder(contents []string) results {
	sts := s.findScanAllInKeyOrder(contents)
	if sts == nil {
		return nil
	}

	return s.scansToResults(sts)
}

// ScanText
// scan text, merging identical scanned text
func (s *scanner) ScanText(contents []string) results {
	sts := s.findScanAll(contents)
	if sts == nil {
		return nil
	}

	var rets results
	resultIndex := make(map[string]int)

	for _, st := range sts {
		text := st.text

		if index, ok := resultIndex[text]; ok {
			rets[index].texts[text] += 1
			rets[index].amount += 1
		} else {
			rets = append(rets, s.scanToResult(st))
			resultIndex[text] = len(rets) - 1
		}
	}

	return rets
}

// ScanFirstText
// the leftmost scanned text
func (s *scanner) ScanFirstText(contents []string) results {
	sts := s.findScanAll(contents)
	if sts == nil {
		return nil
	}

	return s.scanToResults(sts[0])
}

// ScanLastText
// the rightmost scanned text
func (s *scanner) ScanLastText(contents []string) results {
	sts := s.findScanAll(contents)
	if sts == nil {
		return nil
	}

	return s.scanToResults(sts[len(sts)-1])
}

// ScanMostText
// the text that has been scanned the most times
func (s *scanner) ScanMostText(contents []string) results {
	rets := s.ScanText(contents)
	if rets == nil {
		return nil
	}

	sort.Slice(rets, func(i, j int) bool {
		return rets[i].amount > rets[j].amount
	})

	return rets[0:1]
}

// ScanKey
// scan key, merging identical keywords (and their scanned texts)
// in scanned key order
func (s *scanner) ScanKey(contents []string) results {
	sts := s.findScanAllInKeyOrder(contents)
	if sts == nil {
		return nil
	}

	var rets results
	resultIndex := make(map[string]int)

	for _, st := range sts {
		key := st.keyword
		text := st.text

		if index, ok := resultIndex[key]; ok {
			rets[index].texts[text] += 1
			rets[index].amount += 1
		} else {
			rets = append(rets, s.scanToResult(st))
			resultIndex[key] = len(rets) - 1
		}
	}

	return rets
}

// ScanFirstKey
// the first scanned key
func (s *scanner) ScanFirstKey(contents []string) results {
	sts := s.findScanAllInKeyOrder(contents)
	if sts == nil {
		return nil
	}

	return s.scanToResults(sts[0])
}

// ScanLastKey
// the last scanned key
func (s *scanner) ScanLastKey(contents []string) results {
	sts := s.findScanAllInKeyOrder(contents)
	if sts == nil {
		return nil
	}

	return s.scanToResults(sts[len(sts)-1])
}

// ScanMostKey
// scan most key: the keyword with the highest match count
func (s *scanner) ScanMostKey(contents []string) results {
	rets := s.ScanKey(contents)
	if rets == nil {
		return nil
	}

	sort.Slice(rets, func(i, j int) bool {
		return rets[i].amount > rets[j].amount
	})

	return rets[0:1]
}

// ScanLabel
// scan label, merging identical tag combinations
func (s *scanner) ScanLabel(contents []string) labelResults {
	sts := s.findScanAll(contents)
	if sts == nil {
		return nil
	}

	return s.scansToLabelResults(sts)
}

// ScanLabelMostText
// scan label most text: from merged tag combinations, pick the one with the most texts
func (s *scanner) ScanLabelMostText(contents []string) labelResults {
	rets := s.ScanLabel(contents)
	if rets == nil {
		return nil
	}

	sort.Slice(rets, func(i, j int) bool {
		return rets[i].amount > rets[j].amount
	})

	return rets[0:1]
}

func (s *scanner) addKeyFormatters(fs ...ScannerKeyFormatter) {
	s.keyFormatters = append(s.keyFormatters, fs...)
}

// correctBadKeyOfGroupName
// 避免不合法的分组名称
func (s *scanner) correctBadKeyOfGroupName(key string, keyIndex int) string {
	newKey := fmt.Sprintf("%s%d", s.groupNamePrefix, keyIndex)
	s.groupNameMap[newKey] = key

	return newKey
}

func (s *scanner) scansToResults(sts []scannedText) results {
	rets := make(results, 0, len(sts))
	for k := range sts {
		rets = append(rets, s.scanToResult(sts[k]))
	}

	return rets
}

func (s *scanner) scanToResults(st scannedText) results {
	return results{s.scanToResult(st)}
}

func (s *scanner) scanToResult(st scannedText) result {
	r := newResult()
	r.keyword = st.keyword
	r.texts[st.text] = 1
	r.amount = 1
	maps.Copy(r.tags, s.regexpItems[r.keyword])

	return r
}

func (s *scanner) scansToLabelResults(sts []scannedText) labelResults {
	var rets labelResults
	resultIndex := make(map[string]int)

	for _, st := range sts {
		key := st.keyword
		text := st.text

		tags := s.regexpItems[key]

		tagsIdentity := ""
		for _, tag := range s.tagNames {
			tagsIdentity += "-" + tags[tag]
		}

		if index, ok := resultIndex[tagsIdentity]; ok {
			ret := rets[index]
			if _, ok1 := ret.match[key]; ok1 {
				ret.match[key][text] += 1
			} else {
				ret.match[key] = map[string]int{text: 1}
				ret.keywords = append(ret.keywords, key)
			}

			ret.amount += 1
			rets[index] = ret
		} else {
			ret := newLabelResult()
			ret.identity = tagsIdentity
			maps.Copy(ret.tags, tags)
			ret.match[key] = map[string]int{text: 1}
			ret.keywords = append(ret.keywords, key)
			ret.amount += 1

			rets = append(rets, ret)
			resultIndex[tagsIdentity] = len(rets) - 1
		}
	}

	return rets
}

// findScanAllInKeyOrder
// all regex search results, sorted by scanned keyword order
func (s *scanner) findScanAllInKeyOrder(contents []string) []scannedText {
	sts := s.findScanAll(contents)
	if sts == nil {
		return nil
	}

	sort.SliceStable(sts, func(i, j int) bool {
		return s.keyIndex[sts[i].keyword] < s.keyIndex[sts[j].keyword]
	})

	return sts
}

// findScanAll
// all regex search results, in scan text order (leftmost first)
func (s *scanner) findScanAll(contents []string) []scannedText {
	rets := make([]scannedText, 0)

	for i, content := range contents {
		ret := s.findAllSubMatch(i, content, -1)
		rets = append(rets, ret...)
	}

	if len(rets) == 0 {
		return nil
	}

	return rets
}

// findScanFirst
// first regex match result
func (s *scanner) findScanFirst(contents []string) []scannedText {
	var rets []scannedText

	for i, content := range contents {
		rets = s.findAllSubMatch(i, content, 1)
		if len(rets) == 0 {
			break
		}
	}

	return rets
}

func (s *scanner) findAllSubMatch(index int, content string, n int) []scannedText {
	if !s.hasItems {
		return nil
	}

	// 所有分组的匹配结果
	allSubMatchIndex := s.regexp.FindAllStringSubmatchIndex(content, n)

	sts := make([]scannedText, 0, len(allSubMatchIndex))
	for _, subMatch := range allSubMatchIndex {
		_mt := scannedText{
			index: index,
			start: subMatch[0],
			stop:  subMatch[1],
		}

		smLen := len(subMatch)
		for i := 2; i < smLen; i += 2 {
			if subMatch[i] != -1 {
				_mt.text = content[subMatch[i]:subMatch[i+1]]
				_mt.keyword = s.getGroupName(i/2, _mt.text)

				break
			}
		}

		sts = append(sts, _mt)
	}

	return sts
}

func (s *scanner) getGroupName(groupIndex int, text string) string {
	group := s.allSubGroupNames[groupIndex]

	if group == s.groupNamePrefix {
		group = text
	} else {
		if key, ok := s.groupNameMap[group]; ok {
			group = key
		}
	}

	return group
}

// findAllSubMatch
// [][keyword, scanned text]
func (s *scanner) findAllSubMatchBackup(content string, n int) [][]string {
	if !s.hasItems {
		return nil
	}

	// 所有分组的匹配结果
	allSubMatch := s.regexp.FindAllStringSubmatch(content, n)

	matches := make([][]string, 0, len(allSubMatch))
	for _, subMatch := range allSubMatch {
		for k, text := range subMatch {
			if text == "" || k == 0 {
				continue
			}

			group := s.getGroupName(k, text)
			matches = append(matches, []string{group, text})

			break
		}
	}

	return matches
}
