package match

import (
	"fmt"
	"strings"
)

type matcherConfig struct {
	ignoreCase bool
}

type matcher struct {
	keyName     string
	labelsName  []string
	icKeyValues []string
	items       []map[string]string
	matchItems  map[string]map[string]string // map[key name]map[label name]label value
	hasItems    bool

	config *matcherConfig
}

func newMatcher(keyName string, items []map[string]string, config *matcherConfig) *matcher {
	if config == nil {
		config = &matcherConfig{}
	}

	m := &matcher{
		keyName: keyName,
		items:   items,
		config:  config,
	}

	if len(items) > 0 {
		m.hasItems = true

		for k, _ := range items[0] {
			if k != keyName {
				m.labelsName = append(m.labelsName, k)
			}
		}

		m.matchItems = make(map[string]map[string]string)
		for _, item := range m.items {
			m.matchItems[item[m.keyName]] = make(map[string]string)
			for _, labelName := range m.labelsName {
				m.matchItems[item[m.keyName]][labelName] = item[labelName]
			}

			if m.config.ignoreCase {
				m.icKeyValues = append(m.icKeyValues, strings.ToLower(item[m.keyName]))
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

func (m *matcher) findAll(contents []string) []map[string]string {
	if !m.hasItems {
		return nil
	}

	var results []map[string]string
	for _, content := range contents {
		if m.config.ignoreCase {
			content = strings.ToLower(content)
		}

		for k, item := range m.items {
			_subStr := ""
			if m.config.ignoreCase {
				_subStr = m.icKeyValues[k]
			} else {
				_subStr = item[m.keyName]
			}

			_count := strings.Count(content, _subStr)
			if _count > 0 {
				for i := 0; i < _count; i++ {
					results = append(results, item)
				}

				// replace 防止重复 count
				content = strings.ReplaceAll(content, _subStr, fmt.Sprintf("%c", 0x00))
			}
		}
	}

	return results
}

func (m *matcher) findFirst(contents []string) []map[string]string {
	if !m.hasItems {
		return nil
	}

	var results []map[string]string
	for _, content := range contents {
		if m.config.ignoreCase {
			content = strings.ToLower(content)
		}

		for k, item := range m.items {
			_subStr := ""
			if m.config.ignoreCase {
				_subStr = m.icKeyValues[k]
			} else {
				_subStr = item[m.keyName]
			}

			if strings.Contains(content, _subStr) {
				results = append(results, item)
				break
			}
		}
	}

	return results
}

func (m *matcher) toResults(items []map[string]string) Results {
	var results Results
	keysIndex := make(map[string]int)

	for _, item := range items {
		if _index, ok := keysIndex[m.keyName]; ok {
			results[_index].Num += 1
		} else {
			result := NewResult()
			result.Key = item[m.keyName]
			result.Num = 1
			result.Tags = m.matchItems[item[m.keyName]]

			results = append(results, result)
		}
	}

	return results
}

func (m *matcher) toLabelResults(items []map[string]string) LabelResults {
	var results LabelResults
	labelsIndex := make(map[string]int)

	for _, item := range items {
		_labelsIdentity := ""
		for _, _labelName := range m.labelsName {
			_labelsIdentity += "-" + item[_labelName]
		}

		if _index, ok := labelsIndex[_labelsIdentity]; ok {
			result := results[_index]
			if _, ok1 := result.Match[item[m.keyName]]; !ok1 {
				result.Keys = append(result.Keys, item[m.keyName])
			}

			result.Match[item[m.keyName]] += 1
			result.MatchAmount += 1
		} else {
			result := NewLabelResult()
			result.Identity = _labelsIdentity
			result.Labels = m.matchItems[item[m.keyName]]
			result.Keys = append(result.Keys, item[m.keyName])
			result.Match[item[m.keyName]] = 1
			result.MatchAmount = 1

			results = append(results, result)
			labelsIndex[_labelsIdentity] = len(results) - 1
		}
	}

	return results
}
