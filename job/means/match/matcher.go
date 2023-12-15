package match

import (
	"fmt"
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
	keyName  string
	tagsName []string
	hasItems bool

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
		for k, _ := range items[0] {
			if k != keyName {
				m.tagsName = append(m.tagsName, k)
			}
		}

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

func (m *matcher) MatchKey(contents []string) Results {
	items := m.findAll(contents)
	if items == nil {
		return nil
	}

	var results Results
	resultsIndex := make(map[string]int)

	for _, item := range items {
		if index, ok := resultsIndex[item.keyword]; ok {
			for text, _n := range item.textsAmount {
				results[index].TextsAmount[text] += _n
			}

			results[index].Amount += item.amount
		} else {
			results = append(results, m.toResult(item))
		}
	}

	return results
}

func (m *matcher) MatchFirstKey(contents []string) Results {
	results := m.MatchKey(contents)
	if results == nil {
		return nil
	}

	return results[0:1]
}

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

func (m *matcher) MatchText(contents []string) Results {
	items := m.findAll(contents)
	if items == nil {
		return nil
	}

	var results Results
	resultsIndex := make(map[string]int)

	for _, item := range items {
		for text, _n := range item.textsAmount {
			if index, ok := resultsIndex[text]; ok {
				results[index].TextsAmount[text] += _n
				results[index].Amount += _n
			} else {
				result := NewResult()
				result.Keyword = item.keyword
				result.Tags = item.tags
				result.TextsAmount[text] = item.textsAmount[text]
				result.Amount = item.textsAmount[text]

				results = append(results, result)
				resultsIndex[text] = len(results) - 1
			}
		}
	}

	return results
}

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

func (m *matcher) MatchFirstText(contents []string) Results {
	items := m.findAll(contents)
	if items == nil {
		return nil
	}

	rets := m.toResults(items)
	sort.Slice(rets, func(i, j int) bool {
		return rets[i].FirstIndex < rets[j].FirstIndex
	})

	return rets[0:1]
}

func (m *matcher) MatchLastText(contents []string) Results {
	items := m.findAll(contents)
	if items == nil {
		return nil
	}

	return m.toResults(items[len(items)-2:])
}

func (m *matcher) MatchLabel(contents []string) LabelResults {
	items := m.findAll(contents)
	if items == nil {
		return nil
	}

	return m.toLabelResults(items)
}

func (m *matcher) MatchFirstLabel(contents []string) LabelResults {
	results := m.MatchLabel(contents)
	if results == nil {
		return nil
	}

	return results[0:1]
}

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

func (m *matcher) findAll(contents []string) seekResults {
	if !m.hasItems {
		return nil
	}

	var results seekResults
	for _, content := range contents {
		originContent := content

		if m.config.ignoreCase {
			content = strings.ToLower(content)
		}

		rets, sc, ok := m.seekingAll(m.allSeek, originContent, content)
		if ok {
			results = append(results, rets...)
		}

		if m.config.debug {
			m.debugInfo(originContent, sc, results)
		}
	}

	return results
}

func (m *matcher) findFirst(contents []string) seekResults {
	if !m.hasItems {
		return nil
	}

	var results seekResults
	for _, content := range contents {
		originContent := content

		if m.config.ignoreCase {
			content = strings.ToLower(content)
		}

		rets, sc, ok := m.seekingFirst(m.allSeek, originContent, content)
		if ok {
			results = append(results, rets...)
		}

		if m.config.debug {
			m.debugInfo(originContent, sc, results)
		}

		if results != nil {
			break
		}
	}

	return results
}

func (m *matcher) seekingAll(seekers []seeker, originContent, content string) (seekResults, seekContent, bool) {
	return m.seeking(seekers, originContent, content, false)
}

func (m *matcher) seekingFirst(seekers []seeker, originContent, content string) (seekResults, seekContent, bool) {
	return m.seeking(seekers, originContent, content, true)
}

func (m *matcher) seeking(seekers []seeker, originContent, content string, onlyFirst bool) (seekResults, seekContent, bool) {
	var results seekResults

	has := false
	var ok bool
	var sr seekResult
	var sc seekContent
	for _, _seeker := range seekers {
		sr, sc, ok = _seeker.seeking(originContent, content)
		if ok {
			results = append(results, sr)
			has = true

			originContent = sc.origin
			content = sc.content

			if onlyFirst {
				results = results[0:1]

				break
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
	result.Tags = item.tags
	result.TextsAmount = item.textsAmount
	for _, text := range item.texts {
		result.Texts = append(result.Texts, Text{
			Text:  text.text,
			Start: text.start,
			Width: text.width,
		})
	}
	result.FirstIndex = result.Texts[0].Start
	result.LastIndex = result.Texts[len(result.Texts)-1].Start
	result.Amount = item.amount

	return result
}

func (m *matcher) toLabelResults(items seekResults) LabelResults {
	var results LabelResults
	tagsIndex := make(map[string]int)

	for _, item := range items {
		tagsIdentity := ""
		for _, _tn := range m.tagsName {
			tagsIdentity += "-" + item.tags[_tn]
		}

		if _index, ok := tagsIndex[tagsIdentity]; ok {
			result := results[_index]
			if _, ok1 := result.Match[item.keyword]; !ok1 {
				result.Keywords = append(result.Keywords, item.keyword)
				result.Match[item.keyword] = make(map[string]int)
			}

			for _text, _n := range item.textsAmount {
				result.Match[item.keyword][_text] += _n
			}

			result.Amount += item.amount

			results[_index] = result
		} else {
			result := NewLabelResult()
			result.Identity = tagsIdentity
			result.Tags = item.tags
			result.Keywords = append(result.Keywords, item.keyword)
			result.Match[item.keyword] = item.textsAmount
			result.Amount = item.amount

			results = append(results, result)
			tagsIndex[tagsIdentity] = len(results) - 1
		}
	}

	return results
}

func (m *matcher) debugInfo(origin string, sc seekContent, rets seekResults) {
	fmt.Println(fmt.Sprintf("%-16s", "origin:"), origin)
	fmt.Println(fmt.Sprintf("%s", "matched origin: "), sc.origin)
	fmt.Println(fmt.Sprintf("%s", "matched content:"), sc.content)
	fmt.Println("results:")
	for i, rt := range rets {
		fmt.Println(fmt.Sprintf("  %-3d%+v", i, rt))
	}

	fmt.Println()
}
