package match

import (
	"fmt"
	"maps"
	"strings"
)

var _placeholder = fmt.Sprintf("%c", 0x00)

type seekType uint8

const (
	seekAccurate seekType = iota
	seekFuzzy
)

type seekContent struct {
	// 去除匹配项后（匹配项被替换为 placeholder）的 origin，如果无匹配项则是匹配前 origin
	origin string

	// 如果 ignore case，content 为 lower
	// 去除匹配项后（匹配项被替换为 placeholder）的 content，如果无匹配项则是匹配前 content
	content string
}

type textResult struct {
	text  string
	start int // start 包含
	width int // width unit byte
}

type seekResult struct {
	keyword     string            // origin keyword
	texts       []textResult      // matched texts
	textsAmount map[string]int    // map[text]text amount
	tags        map[string]string // matched tags
	amount      int               // keyword matched amount
}

func newSeekResult() seekResult {
	return seekResult{
		textsAmount: make(map[string]int),
	}
}

type seekResults []seekResult

type seekConfig struct {
	debug bool
}

// seeker
// content，keyword 大小写在 match 已经处理过
// 这里区分大小写
type seeker interface {
	// origin string:
	// toLower string:

	// seekResult
	// seekContent
	// bool: true has matched；false has not matched
	seeking(origin string, toLower string) (seekResult, seekContent, bool)
}

type seek struct {
	config seekConfig
}

func (s *seek) replaceMatchedToPlaceholder(content, matched string) string {
	c := s.matchedToPlaceholder(matched)
	return strings.ReplaceAll(content, matched, c)
}

func (s *seek) matchedToPlaceholder(matched string) string {
	if s.config.debug {
		return strings.Repeat(_placeholder, len(matched))
	} else {
		return _placeholder
	}
}

func newSeeker(index int, originKey, key string, tags map[string]string, matcherConfig *matcherConfig) (seeker, seekType) {
	config := seekConfig{debug: matcherConfig.debug}
	newTags := maps.Clone(tags)

	if matcherConfig.enableFuzzy && strings.Index(key, matcherConfig.fuzzyConfig.Sep) > -1 {
		return newFuzzy(index, originKey, key, newTags, matcherConfig.fuzzyConfig, config), seekFuzzy
	} else {
		return newAccurate(index, originKey, key, newTags, config), seekAccurate
	}
}
