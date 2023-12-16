package match

import (
	"fmt"
	"maps"
	"strings"
	"unicode/utf8"
)

var _placeholder = fmt.Sprintf("%c", 0x00)

type seekType uint8

const (
	seekAccurate seekType = iota
	seekFuzzy
)

type seekContent struct {
	// 多个 content， content 的序号
	index int

	// 去除匹配项后（匹配项被替换为 placeholder）的 origin，如果无匹配项则是匹配前 origin
	origin string

	// 如果 ignore case，content 为 lower
	// 去除匹配项后（匹配项被替换为 placeholder）的 content，如果无匹配项则是匹配前 content
	content string
}

type seekResult struct {
	index   int               // 多个 content， content 的序号
	start   int               // start 包含
	width   int               // width unit byte
	keyword string            // origin keyword
	text    string            // matched text
	tags    map[string]string // matched tags
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
	seeking(seekContent) (seekResults, seekContent, bool)
}

type seek struct {
	config seekConfig
}

func (s *seek) replaceMatchedToPlaceholder(content, matched string) string {
	return strings.ReplaceAll(content, matched, _placeholder)
}

func (s *seek) matchedToPlaceholder(matched string) string {
	return strings.Repeat(_placeholder, len(matched))
}

func (s *seek) debugToPlaceholder(matched string) string {
	_len := len(matched)
	_runeLen := utf8.RuneCountInString(matched)
	_zhLen := (_len - _runeLen) / 2
	return strings.Repeat(_placeholder, _runeLen+_zhLen)
}

func newSeeker(keyIndex int, originKey, key string, tags map[string]string, matcherConfig *matcherConfig) (seeker, seekType) {
	config := seekConfig{debug: matcherConfig.debug}
	newTags := maps.Clone(tags)

	if matcherConfig.enableFuzzy && strings.Index(key, matcherConfig.fuzzyConfig.Sep) > -1 {
		return newFuzzy(keyIndex, originKey, key, newTags, matcherConfig.fuzzyConfig, config), seekFuzzy
	} else {
		return newAccurate(keyIndex, originKey, key, newTags, config), seekAccurate
	}
}
