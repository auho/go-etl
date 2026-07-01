package match

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/auho/go-etl/v3/job/extract"
)

// seekMode controls the order in which seekers (accurate vs fuzzy) are tried.
type seekMode int

const (
	// modeSequence: try all seekers (accurate + fuzzy) in definition order.
	modeSequence seekMode = iota
	// modePriorityAccurate: try accurate seekers first, fall back to fuzzy.
	modePriorityAccurate
	// modePriorityFuzzy: try fuzzy seekers first, fall back to accurate.
	modePriorityFuzzy
)

// ScannerConfig holds all knobs that control scanner behavior.
type ScannerConfig struct {
	IgnoreCase bool        // if true, keywords are lowercased before matching
	Debug      bool        // if true, prints debug info during scanning
	Mode       seekMode    // ordering strategy for accurate/fuzzy seekers
	Fuzzy      FuzzyConfig // fuzzy matching configuration
}

func (sc *ScannerConfig) check() {
	sc.Fuzzy.check()
}

func defaultScanner(rule extract.Rule, sc ScannerConfig) (*scanner, error) {
	return newScannerFromRule(rule, sc)
}

func newScannerFromRule(rule extract.Rule, sc ScannerConfig) (*scanner, error) {
	items, err := rule.ItemsAlias()
	if err != nil {
		return nil, fmt.Errorf("ItemsAlias: %w", err)
	}

	return newScanner(rule.KeywordNameAlias(), items, sc), nil
}

// scanner drives keyword scanning over content using a list of seekers.
type scanner struct {
	hasItems bool

	keyName  string
	tagNames []string

	allSeekers      []seeker // all seekers in evaluation order
	fuzzySeekers    []seeker // fuzzy (approximate) seekers
	accurateSeekers []seeker // accurate (exact) seekers

	config ScannerConfig
}

// newScanner creates a scanner from items. Each item's keyName value becomes a keyword;
// remaining columns become tags. The ScannerConfig controls case sensitivity, mode, and fuzzy settings.
func newScanner(keyName string, items []map[string]string, sc ScannerConfig) *scanner {
	sc.check()

	m := &scanner{
		keyName: keyName,
		config:  sc,
	}

	if len(items) == 0 {
		return m
	}

	m.hasItems = true

	// tags name
	for k := range items[0] {
		if k != keyName {
			m.tagNames = append(m.tagNames, k)
		}
	}

	sort.SliceStable(m.tagNames, func(i, j int) bool {
		return m.tagNames[i] < m.tagNames[j]
	})

	for _i, item := range items {
		var _keyValue string
		_originKeyValue := item[keyName]
		if sc.IgnoreCase {
			_keyValue = strings.ToLower(_originKeyValue)
		} else {
			_keyValue = _originKeyValue
		}

		_tags := make(map[string]string)
		for _, _ln := range m.tagNames {
			_tags[_ln] = item[_ln]
		}

		_seeker, _st := newSeeker(_i, _originKeyValue, _keyValue, _tags, sc.Fuzzy, seekConfig{debug: sc.Debug})
		if sc.Mode == modeSequence {
			m.allSeekers = append(m.allSeekers, _seeker)
		} else {
			if _st == seekAccurate {
				m.accurateSeekers = append(m.accurateSeekers, _seeker)
			} else {
				m.fuzzySeekers = append(m.fuzzySeekers, _seeker)
			}
		}
	}

	switch sc.Mode {
	case modeSequence:
	case modePriorityAccurate:
		m.allSeekers = append(m.accurateSeekers, m.fuzzySeekers...)
	case modePriorityFuzzy:
		m.allSeekers = append(m.fuzzySeekers, m.accurateSeekers...)
	default:
		panic(fmt.Sprintf("unknown mode[%d]", sc.Mode))
	}

	return m
}

// Scan
// all scanned
// in scanned keyword order
func (s *scanner) Scan(contents []string) results {
	items := s.findAll(contents)
	if items == nil {
		return nil
	}

	return s.toResults(items)
}

// ScanInTextOrder
// all scanned
// in scanned text order
func (s *scanner) ScanInTextOrder(contents []string) results {
	items := s.findAllInTextOrder(contents)
	if items == nil {
		return nil
	}

	return s.toResults(items)
}

// ScanText
// scan text 合并相同的 scanned text
func (s *scanner) ScanText(contents []string) results {
	items := s.findAllInTextOrder(contents)
	if items == nil {
		return nil
	}

	var rets results
	resultIndex := make(map[string]int)

	for _, item := range items {
		text := item.text

		if index, ok := resultIndex[text]; ok {
			rets[index].texts[text] += 1
			rets[index].amount += 1
		} else {
			rets = append(rets, s.toResult(item))
			resultIndex[text] = len(rets) - 1
		}
	}

	return rets
}

// ScanFirstText
// the leftmost scanned text
func (s *scanner) ScanFirstText(contents []string) results {
	items := s.findAllInTextOrder(contents)
	if items == nil {
		return nil
	}

	rets := s.toResults(items)
	return rets[0:1]
}

// ScanLastText
// the rightmost scanned text
func (s *scanner) ScanLastText(contents []string) results {
	items := s.findAllInTextOrder(contents)
	if items == nil {
		return nil
	}

	return s.toResults(items[len(items)-1:])
}

// ScanMostText
// the text that has been scanned the most times
func (s *scanner) ScanMostText(contents []string) results {
	rets := s.ScanText(contents)
	if rets == nil {
		return nil
	}

	sort.SliceStable(rets, func(i, j int) bool {
		return rets[i].amount > rets[j].amount
	})

	return rets[0:1]
}

// ScanKey
// scan key 合并相同的 keyword（同时也合并 scanned text）
// in scanned key order
func (s *scanner) ScanKey(contents []string) results {
	items := s.findAll(contents)
	if items == nil {
		return nil
	}

	var rets results
	resultIndex := make(map[string]int)

	for _, item := range items {
		key := item.keyword
		text := item.text

		if index, ok := resultIndex[key]; ok {
			rets[index].texts[text] += 1
			rets[index].amount += 1
		} else {
			rets = append(rets, s.toResult(item))
			resultIndex[key] = len(rets) - 1
		}
	}

	return rets
}

// ScanFirstKey
// the first scanned key
func (s *scanner) ScanFirstKey(contents []string) results {
	rets := s.findFirst(contents)
	if rets == nil {
		return nil
	}

	return s.toResults(rets)
}

// ScanLastKey
// the last scanned key
func (s *scanner) ScanLastKey(contents []string) results {
	items := s.findAll(contents)
	if items == nil {
		return nil
	}

	return s.toResults(items[len(items)-1:])
}

// ScanMostKey
// scan most key 被匹配次数最多的 keyword
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
// scan label 合并重复的 tags 组合
func (s *scanner) ScanLabel(contents []string) labelResults {
	items := s.findAll(contents)
	if items == nil {
		return nil
	}

	return s.toLabelResults(items)
}

// ScanLabelMostText
// scan label most text 合并重复的 tags 组合中，text 最多次数
func (s *scanner) ScanLabelMostText(contents []string) labelResults {
	labels := s.ScanLabel(contents)
	if labels == nil {
		return nil
	}

	sort.Slice(labels, func(i, j int) bool {
		return labels[i].amount > labels[j].amount
	})

	return labels[0:1]
}

// findAllInTextOrder
// all scan, in scanned text order
// the leftmost text is at the front
func (s *scanner) findAllInTextOrder(contents []string) seekResults {
	items := s.seekContents(contents, false)
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
// all scanned results, in scanned keyword order
func (s *scanner) findAll(contents []string) seekResults {
	return s.seekContents(contents, false)
}

// findFirst
// the first scanned keyword found
func (s *scanner) findFirst(contents []string) seekResults {
	return s.seekContents(contents, true)
}

// seekContents runs all seekers against every content string.
// If onlyFirst is true, stops after the first match.
func (s *scanner) seekContents(contents []string, onlyFirst bool) seekResults {
	if !s.hasItems {
		return nil
	}

	var allRets seekResults
	var isBreak bool
	for i, content := range contents {
		var contentResults seekResults

		originContent := content
		if s.config.IgnoreCase {
			content = strings.ToLower(content)
		}

		sc := seekContent{
			maxSeekNum: -1,
			index:      i,
			origin:     originContent,
			content:    content,
		}

		if onlyFirst {
			sc.maxSeekNum = 1
		}

		var ok bool
		var rets seekResults
		rets, sc, ok = s.seeking(s.allSeekers, sc, onlyFirst)
		if ok {
			if onlyFirst {
				contentResults = rets[0:1]

				isBreak = true
				goto LOOP
			} else {
				contentResults = append(contentResults, rets...)
			}
		}
	LOOP:

		if s.config.Debug {
			s.debugInfo(i, originContent, sc, contentResults)
		}

		allRets = append(allRets, contentResults...)

		if isBreak {
			break
		}
	}

	return allRets
}

// seeking iterates seekers and collects their results.
// If onlyFirst, returns the first result set and stops.
func (s *scanner) seeking(seekers []seeker, sc seekContent, onlyFirst bool) (seekResults, seekContent, bool) {
	var allRets seekResults

	has := false
	var ok bool
	var rets seekResults
	for _, _seeker := range seekers {
		rets, sc, ok = _seeker.seeking(sc)
		if ok {
			has = true
			if onlyFirst {
				allRets = rets[0:1]

				break
			} else {
				allRets = append(allRets, rets...)
			}
		}
	}

	return allRets, sc, has
}

// toResults converts seekResults to a results slice.
func (s *scanner) toResults(items seekResults) results {
	var rets results

	for _, item := range items {
		rets = append(rets, s.toResult(item))
	}

	return rets
}

// toResult converts a single seekResult to a result.
func (s *scanner) toResult(item seekResult) result {
	rets := newResult()
	rets.keyword = item.keyword
	rets.tags = maps.Clone(item.tags)
	rets.texts = map[string]int{item.text: 1}
	rets.amount = 1

	return rets
}

// toLabelResults converts seekResults to labelResults, merging by tags identity.
func (s *scanner) toLabelResults(items seekResults) labelResults {
	var rets labelResults
	resultIndex := make(map[string]int)

	for _, item := range items {
		key := item.keyword
		text := item.text

		tagsIdentity := ""
		for _, _tn := range s.tagNames {
			tagsIdentity += "-" + item.tags[_tn]
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
			maps.Copy(ret.tags, item.tags)
			ret.match[key] = map[string]int{text: 1}
			ret.keywords = append(ret.keywords, key)
			ret.amount += 1

			rets = append(rets, ret)
			resultIndex[tagsIdentity] = len(rets) - 1
		}
	}

	return rets
}

func (s *scanner) debugInfo(index int, originContent string, sc seekContent, rets seekResults) {
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
	preStart := 0
	preStop := 0
	for _, ret := range newRest {
		debugContent += originContent[preStart:ret.start]

		_len := len(ret.text)
		_runeLen := utf8.RuneCountInString(ret.text)
		_zhLen := (_len - _runeLen) / 2

		debugContent += strings.Repeat(placeholder, _runeLen+_zhLen)
		preStart = ret.start + ret.width
		preStop = ret.start + ret.width
	}

	debugContent += originContent[preStop:]
	fmt.Println(fmt.Sprintf("%-16s", "index:"), index)
	fmt.Println(fmt.Sprintf("%-16s", "origin:"), originContent)
	fmt.Println(fmt.Sprintf("%-16s", "debug origin:"), debugContent)
	fmt.Println(fmt.Sprintf("%-16s", "scanned origin:"), sc.origin)
	fmt.Println(fmt.Sprintf("%-16s", "scanned content:"), sc.content)
	fmt.Println("results:")
	for i, rt := range rets {
		fmt.Println(fmt.Sprintf("  %-3d%+v", i, rt))
	}

	fmt.Println()
}
