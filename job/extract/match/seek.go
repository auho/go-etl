package match

import (
	"fmt"
	"maps"
	"strings"
)

var placeholder = fmt.Sprintf("%c", 0x00)

type seekType uint8

const (
	seekAccurate seekType = iota
	seekFuzzy
)

// seek content
type seekContent struct {
	maxSeekNum int

	// 多个 content， content 的序号
	index int

	// 去除匹配项后（匹配项被替换为 placeholder）的 origin，如果无匹配项则是匹配前 origin
	origin string

	// 如果 ignore case，content 为 lower
	// 去除匹配项后（匹配项被替换为 placeholder）的 content，如果无匹配项则是匹配前 content
	content string
}

// seek result
type seekResult struct {
	index   int               // 多个 content， content 的序号
	start   int               // start 包含
	width   int               // width unit byte
	keyword string            // origin keyword
	text    string            // matched text
	tags    map[string]string // matched tags
}

type seekResults []seekResult

// seek config
type seekConfig struct {
	debug bool
}

// seeker
// content，keyword 大小写在 scanner 已经处理过
// 这里区分大小写
type seeker interface {
	// origin string:
	// toLower string:

	// seekResult
	// seekContent
	// bool: true has matched；false has not matched
	seeking(seekContent) (seekResults, seekContent, bool)
}

// seek
type seek struct {
	config seekConfig
}

func (s *seek) replaceMatchedToPlaceholder(content, matched string) string {
	return strings.ReplaceAll(content, matched, placeholder)
}

func (s *seek) matchedToPlaceholder(matched string) string {
	return strings.Repeat(placeholder, len(matched))
}

func newSeeker(keyIndex int, originKey, key string, tags map[string]string, fc FuzzyConfig, sc seekConfig) (seeker, seekType) {
	newTags := maps.Clone(tags)

	if fc.shouldFuzzy(key) {
		return newFuzzy(keyIndex, originKey, key, newTags, fc, sc), seekFuzzy
	} else {
		return newAccurate(keyIndex, originKey, key, newTags, sc), seekAccurate
	}
}
