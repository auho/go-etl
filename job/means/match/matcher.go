package match

import (
	"fmt"
	"sort"
	"strings"
)

const (
	modeSequence = iota
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
	mode        int
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
				results[index].Texts[text] += _n
			}

			results[index].Amount += item.amount
		} else {
			results = append(results, m.toResult(item))
		}
	}

	return results
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
				results[index].Texts[text] += _n
				results[index].Amount += _n
			} else {
				result := NewResult()
				result.Keyword = item.keyword
				result.Tags = item.tags
				result.Texts[text] = item.textsAmount[text]
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
	items := m.findFirst(contents)
	if items == nil {
		return nil
	}

	return m.toResults(items)
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

		var ok bool
		var _rts seekResults
		switch m.config.mode {
		case modeSequence:
			_rts, content, ok = m.seekingAll(m.allSeek, content)
			if ok {
				results = append(results, _rts...)
			}
		case modePriorityAccurate:
			_rts, content, ok = m.seekingAll(m.accurateSeek, content)
			if ok {
				results = append(results, _rts...)
			}

			_rts, _, ok = m.seekingAll(m.fuzzySeek, content)
			if ok {
				results = append(results, _rts...)
			}
		case modePriorityFuzzy:
			_rts, content, ok = m.seekingAll(m.fuzzySeek, content)
			if ok {
				results = append(results, _rts...)
			}

			_rts, content, ok = m.seekingAll(m.accurateSeek, content)
			if ok {
				results = append(results, _rts...)
			}
		default:
			panic("unknown mode")
		}

		if m.config.debug {
			fmt.Println(originContent, content, results)
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

		var ok bool
		var _rts seekResults
		switch m.config.mode {
		case modeSequence:
			_rts, content, ok = m.seekingFirst(m.allSeek, content)
			if ok {
				results = append(results, _rts...)
			}
		case modePriorityAccurate:
			_rts, content, ok = m.seekingFirst(m.accurateSeek, content)
			if ok {
				results = append(results, _rts...)
				break
			}

			_rts, content, ok = m.seekingFirst(m.fuzzySeek, content)
			if ok {
				results = append(results, _rts...)
			}
		case modePriorityFuzzy:
			_rts, content, ok = m.seekingFirst(m.fuzzySeek, content)
			if ok {
				results = append(results, _rts...)
				break
			}

			_rts, content, ok = m.seekingFirst(m.accurateSeek, content)
			if ok {
				results = append(results, _rts...)
			}
		default:
			panic("unknown mode")
		}

		if m.config.debug {
			m.debugInfo(originContent, content, results)
		}

		if results != nil {
			break
		}
	}

	return results
}

func (m *matcher) seekingAll(seekers []seeker, content string) (seekResults, string, bool) {
	return m.seeking(seekers, content, false)
}

func (m *matcher) seekingFirst(seekers []seeker, content string) (seekResults, string, bool) {
	return m.seeking(seekers, content, true)
}

func (m *matcher) seeking(seekers []seeker, content string, onlyFirst bool) (seekResults, string, bool) {
	var results seekResults

	has := false
	var ok bool
	var _srt seekResult
	for _, _seeker := range seekers {
		_srt, content, ok = _seeker.seeking(content)
		if ok {
			results = append(results, _srt)
			has = true

			if onlyFirst {
				results = results[0:1]

				break
			}
		}
	}

	return results, content, has
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
	result.Texts = item.textsAmount
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

func (m *matcher) debugInfo(os, s string, rts seekResults) {
	fmt.Println(os)
	fmt.Println(s)
	fmt.Println(rts)
}
