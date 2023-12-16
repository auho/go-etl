package match

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"
)

type seekMode int

const (
	modeSequence seekMode = iota
	modePriorityAccurate
	modePriorityFuzzy
)

type FuzzyConfig struct {
	Window int
	Sep    string
}

func (fc *FuzzyConfig) check() {
	if fc.Sep == "" {
		fc.Sep = "_"
	}

	if fc.Window <= 0 {
		fc.Window = 3
	}
}

type matcherConfig struct {
	ignoreCase  bool
	mode        seekMode
	enableFuzzy bool
	fuzzyConfig FuzzyConfig
	debug       bool
}

func (mc *matcherConfig) check() {
	mc.fuzzyConfig.check()
}

type matcher struct {
	hasItems bool

	keyName  string
	tagsName []string

	allSeek      []seeker
	fuzzySeek    []seeker
	accurateSeek []seeker

	config *matcherConfig
}

func newMatcher(keyName string, items []map[string]string, config *matcherConfig) *matcher {
	if config == nil {
		config = &matcherConfig{}
	}

	config.check()

	m := &matcher{
		keyName: keyName,
		config:  config,
	}

	if len(items) > 0 {
		m.hasItems = true

		// tags name
		for k := range items[0] {
			if k != keyName {
				m.tagsName = append(m.tagsName, k)
			}
		}

		sort.SliceStable(m.tagsName, func(i, j int) bool {
			return m.tagsName[i] < m.tagsName[j]
		})

		for _i, item := range items {
			var _keyValue string
			_originKeyValue := item[keyName]
			if config.ignoreCase {
				_keyValue = strings.ToLower(_originKeyValue)
			} else {
				_keyValue = _originKeyValue
			}

			_tags := make(map[string]string)
			for _, _ln := range m.tagsName {
				_tags[_ln] = item[_ln]
			}

			_seeker, _sm := newSeeker(_i, _originKeyValue, _keyValue, _tags, config)
			if config.mode == modeSequence {
				m.allSeek = append(m.allSeek, _seeker)
			} else {
				if _sm == seekAccurate {
					m.accurateSeek = append(m.accurateSeek, _seeker)
				} else {
					m.fuzzySeek = append(m.fuzzySeek, _seeker)
				}
			}
		}

		switch config.mode {
		case modeSequence:
		case modePriorityAccurate:
			m.allSeek = append(m.accurateSeek, m.fuzzySeek...)
		case modePriorityFuzzy:
			m.allSeek = append(m.fuzzySeek, m.accurateSeek...)
		default:
			panic(fmt.Sprintf("mode [%d] error", config.mode))
		}
	}

	return m
}

// Match
// all matched
// in key order
func (m *matcher) Match(contents []string) Results {
	items := m.findAll(contents)
	if items == nil {
		return nil
	}

	return m.toResults(items)
}

// MatchInTextOrder
// all matched
// in matched text order
func (m *matcher) MatchInTextOrder(contents []string) Results {
	items := m.findAllInTextOrder(contents)
	if items == nil {
		return nil
	}

	return m.toResults(items)
}

// MatchText
// match text 合并相同的 matched text
func (m *matcher) MatchText(contents []string) Results {
	items := m.findAllInTextOrder(contents)
	if items == nil {
		return nil
	}

	var results Results
	resultIndex := make(map[string]int)

	for _, item := range items {
		text := item.text

		if index, ok := resultIndex[text]; ok {
			results[index].Texts[text] += 1
			results[index].Amount += 1
		} else {
			results = append(results, m.toResult(item))
			resultIndex[text] = len(results) - 1
		}
	}

	return results
}

// MatchFirstText
// the leftmost matched text
func (m *matcher) MatchFirstText(contents []string) Results {
	items := m.findAllInTextOrder(contents)
	if items == nil {
		return nil
	}

	rets := m.toResults(items)
	return rets[0:1]
}

// MatchLastText
// the rightmost matched text
func (m *matcher) MatchLastText(contents []string) Results {
	items := m.findAllInTextOrder(contents)
	if items == nil {
		return nil
	}

	return m.toResults(items[len(items)-2:])
}

// MatchMostText
// the text that has been matched the most times
func (m *matcher) MatchMostText(contents []string) Results {
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
func (m *matcher) MatchKey(contents []string) Results {
	items := m.findAll(contents)
	if items == nil {
		return nil
	}

	var results Results
	resultIndex := make(map[string]int)

	for _, item := range items {
		key := item.keyword
		text := item.text

		if index, ok := resultIndex[key]; ok {
			results[index].Texts[text] += 1
			results[index].Amount += 1
		} else {
			results = append(results, m.toResult(item))
			resultIndex[key] = len(results) - 1
		}
	}

	return results
}

// MatchFirstKey
// the first matched key
func (m *matcher) MatchFirstKey(contents []string) Results {
	results := m.findFirst(contents)
	if results == nil {
		return nil
	}

	return m.toResults(results)
}

// MatchLastKey
// the last matched key
func (m *matcher) MatchLastKey(contents []string) Results {
	items := m.findFirst(contents)
	if items == nil {
		return nil
	}

	return m.toResults(items[len(items)-1:])
}

// MatchMostKey
// match most key 被匹配次数最多的 keyword
func (m *matcher) MatchMostKey(contents []string) Results {
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
func (m *matcher) MatchLabel(contents []string) LabelResults {
	items := m.findAll(contents)
	if items == nil {
		return nil
	}

	return m.toLabelResults(items)
}

// MatchLabelMostText
// match label most text 合并重复的 tags 组合中，text 最多次数
func (m *matcher) MatchLabelMostText(contents []string) LabelResults {
	results := m.MatchLabel(contents)
	if results == nil {
		return nil
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Amount > results[j].Amount
	})

	return results[0:1]
}

// findAllInTextOrder
// all match, in matched text order
// the leftmost text is at the front
func (m *matcher) findAllInTextOrder(contents []string) seekResults {
	items := m.seekContents(contents, false)
	if items == nil {
		return nil
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].index < items[j].index {
			return true
		} else if items[i].index > items[j].index {
			return false
		} else {
			return items[i].start < items[j].start
		}
	})

	return items
}

// findAll
// all match, in matched keyword order
func (m *matcher) findAll(contents []string) seekResults {
	return m.seekContents(contents, false)
}

// findMatchFirst
// first matched keyword
func (m *matcher) findFirst(contents []string) seekResults {
	return m.seekContents(contents, true)
}

func (m *matcher) seekContents(contents []string, onlyFirst bool) seekResults {
	if !m.hasItems {
		return nil
	}

	var results seekResults
	for i, content := range contents {
		originContent := content
		if m.config.ignoreCase {
			content = strings.ToLower(content)
		}

		sc := seekContent{
			index:   i,
			origin:  originContent,
			content: content,
		}

		var ok bool
		var rets seekResults
		rets, sc, ok = m.seeking(m.allSeek, sc, onlyFirst)
		if ok {
			if onlyFirst {
				results = rets[0:1]

				break
			} else {
				results = append(results, rets...)
			}
		}

		if m.config.debug {
			m.debugInfo(sc, results)
		}
	}

	return results
}

func (m *matcher) seeking(seekers []seeker, sc seekContent, onlyFirst bool) (seekResults, seekContent, bool) {
	var results seekResults

	has := false
	var ok bool
	var rets seekResults
	for _, _seeker := range seekers {
		rets, sc, ok = _seeker.seeking(sc)
		if ok {
			has = true
			if onlyFirst {
				results = rets[0:1]

				break
			} else {
				results = append(results, rets...)
			}
		}
	}

	return results, sc, has
}

func (m *matcher) toResults(items seekResults) Results {
	var results Results

	for _, item := range items {
		results = append(results, m.toResult(item))
	}

	return results
}

func (m *matcher) toResult(item seekResult) Result {
	result := NewResult()
	result.Keyword = item.keyword
	result.Tags = maps.Clone(item.tags)
	result.Texts = map[string]int{item.text: 1}
	result.Amount = 1

	return result
}

func (m *matcher) toLabelResults(items seekResults) LabelResults {
	var results LabelResults
	resultIndex := make(map[string]int)

	for _, item := range items {
		key := item.keyword
		text := item.text

		tagsIdentity := ""
		for _, _tn := range m.tagsName {
			tagsIdentity += "-" + item.tags[_tn]
		}

		if index, ok := resultIndex[tagsIdentity]; ok {
			result := results[index]
			if _, ok1 := result.Match[key]; ok1 {
				result.Match[key][text] += 1
			} else {
				result.Match[key] = map[string]int{text: 1}
				result.Keywords = append(result.Keywords, key)
			}

			result.Amount += 1
			results[index] = result
		} else {
			result := NewLabelResult()
			result.Identity = tagsIdentity
			maps.Copy(result.Tags, item.tags)
			result.Match[key] = map[string]int{text: 1}
			result.Keywords = append(result.Keywords, key)
			result.Amount += 1

			results = append(results, result)
			resultIndex[tagsIdentity] = len(results) - 1
		}
	}

	return results
}

func (m *matcher) debugInfo(sc seekContent, rets seekResults) {
	newRest := slices.Clone(rets)
	sort.SliceStable(newRest, func(i, j int) bool {
		if newRest[i].index < newRest[j].index {
			return true
		} else if newRest[i].index > newRest[j].index {
			return false
		} else {
			return newRest[i].start < newRest[j].start
		}
	})

	debugContent := ""
	//preStart := 0
	//for _, ret := range newRest {
	//	//debugContent += sc.originContent[preStart:ret.start]
	//	//debugContent += strings.Repeat(_placeholder, ret.width)
	//	preStart = ret.start + ret.width - 1
	//}

	fmt.Println(fmt.Sprintf("%-16s", "debug origin:"), debugContent)
	fmt.Println(fmt.Sprintf("%-16s", "matched origin:"), sc.origin)
	fmt.Println(fmt.Sprintf("%-16s", "matched content:"), sc.content)
	fmt.Println("results:")
	for i, rt := range rets {
		fmt.Println(fmt.Sprintf("  %-3d%+v", i, rt))
	}

	fmt.Println()
}
