package match

import (
	"strings"
)

const (
	modeSequence = iota
	modePriorityAccurate
	modePriorityFuzzy
)

type FuzzyConfig struct {
	Enable bool
	Window int
	Sep    string
}

type matcherConfig struct {
	ignoreCase  bool
	mode        int
	fuzzyConfig FuzzyConfig
}

type matcher struct {
	keyName    string
	labelsName []string
	hasItems   bool

	allSeek      []seeker
	fuzzySeek    []seeker
	accurateSeek []seeker

	config *matcherConfig
}

func newMatcher(keyName string, items []map[string]string, config *matcherConfig) *matcher {
	if config == nil {
		config = &matcherConfig{}
	}

	m := &matcher{
		keyName: keyName,
		config:  config,
	}

	if len(items) > 0 {
		m.hasItems = true

		// labels name
		for k, _ := range items[0] {
			if k != keyName {
				m.labelsName = append(m.labelsName, k)
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

			_labels := make(map[string]string)
			for _, _ln := range m.labelsName {
				_labels[_ln] = item[_ln]
			}

			_seeker, _sm := newSeeker(_i, _originKeyValue, _keyValue, _labels, config.fuzzyConfig)
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

	return m.toResults(items)
}

func (m *matcher) MatchFirstKey(contents []string) Results {
	items := m.findFirst(contents)
	if items == nil {
		return nil
	}

	return m.toResults(items)
}

func (m *matcher) MatchLabel(contents []string) LabelResults {
	items := m.findAll(contents)
	if items == nil {
		return nil
	}

	return m.toLabelResults(items)
}

func (m *matcher) MatchFirstLabel(contents []string) LabelResults {
	items := m.findFirst(contents)
	if items == nil {
		return nil
	}

	return m.toLabelResults(items)
}

func (m *matcher) findAll(contents []string) seekResults {
	if !m.hasItems {
		return nil
	}

	var results seekResults
	for _, content := range contents {
		if m.config.ignoreCase {
			content = strings.ToLower(content)
		}

		var ok bool
		var _rts seekResults
		switch m.config.mode {
		case modeSequence:
			_rts, ok = m.seekingAll(m.allSeek, content)
			if ok {
				results = append(results, _rts...)
			}
		case modePriorityAccurate:
			_rts, ok = m.seekingAll(m.accurateSeek, content)
			if ok {
				results = append(results, _rts...)
			}

			_rts, ok = m.seekingAll(m.fuzzySeek, content)
			if ok {
				results = append(results, _rts...)
			}
		case modePriorityFuzzy:
			_rts, ok = m.seekingAll(m.fuzzySeek, content)
			if ok {
				results = append(results, _rts...)
			}

			_rts, ok = m.seekingAll(m.accurateSeek, content)
			if ok {
				results = append(results, _rts...)
			}
		default:
			panic("unknown mode")
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
		if m.config.ignoreCase {
			content = strings.ToLower(content)
		}

		var ok bool
		var _rts seekResults
		switch m.config.mode {
		case modeSequence:
			_rts, ok = m.seekingFirst(m.allSeek, content)
			if ok {
				results = append(results, _rts...)
			}
		case modePriorityAccurate:
			_rts, ok = m.seekingFirst(m.accurateSeek, content)
			if ok {
				results = append(results, _rts...)
				break
			}

			_rts, ok = m.seekingFirst(m.fuzzySeek, content)
			if ok {
				results = append(results, _rts...)
			}
		case modePriorityFuzzy:
			_rts, ok = m.seekingFirst(m.fuzzySeek, content)
			if ok {
				results = append(results, _rts...)
				break
			}

			_rts, ok = m.seekingFirst(m.accurateSeek, content)
			if ok {
				results = append(results, _rts...)
			}
		default:
			panic("unknown mode")
		}
	}

	return results
}

func (m *matcher) seekingAll(seekers []seeker, content string) (seekResults, bool) {
	return m.seeking(seekers, content, false)
}

func (m *matcher) seekingFirst(seekers []seeker, content string) (seekResults, bool) {
	return m.seeking(seekers, content, true)
}

func (m *matcher) seeking(seekers []seeker, content string, onlyFirst bool) (seekResults, bool) {
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
				break
			}
		}
	}

	return results, has
}

func (m *matcher) toResults(items seekResults) Results {
	var results Results
	keysIndex := make(map[string]int)

	for _, item := range items {
		if _index, ok := keysIndex[m.keyName]; ok {
			results[_index].Num += item.amount
		} else {
			result := NewResult()
			result.Key = item.key
			result.Num = item.amount
			result.Labels = item.labels

			results = append(results, result)
		}
	}

	return results
}

func (m *matcher) toLabelResults(items seekResults) LabelResults {
	var results LabelResults
	labelsIndex := make(map[string]int)

	for _, item := range items {
		_labelsIdentity := ""
		for _, _labelName := range m.labelsName {
			_labelsIdentity += "-" + item.labels[_labelName]
		}

		if _index, ok := labelsIndex[_labelsIdentity]; ok {
			result := results[_index]
			if _, ok1 := result.Match[item.key]; !ok1 {
				result.Keys = append(result.Keys, item.key)
			}

			result.Match[item.key] += item.amount
			result.MatchAmount += item.amount
		} else {
			result := NewLabelResult()
			result.Identity = _labelsIdentity
			result.Labels = item.labels
			result.Keys = append(result.Keys, item.key)
			result.Match[item.key] = item.amount
			result.MatchAmount = item.amount

			results = append(results, result)
			labelsIndex[_labelsIdentity] = len(results) - 1
		}
	}

	return results
}
