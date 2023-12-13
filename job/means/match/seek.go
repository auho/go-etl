package match

import (
	"fmt"
	"maps"
	"strings"
)

var _placeholder = fmt.Sprintf("%c", 0x00)

const (
	seekAccurate = iota
	seekFuzzy
)

type seekResult struct {
	keyword     string            // origin keyword
	texts       []string          // matched texts
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

// seeker
// content，keyword 大小写在 match 已经处理过
// 这里区分大小写
type seeker interface {
	// origin string:
	// toLower string:

	// seekResult: seek result
	// string: 去除匹配项后的 content，如果未匹配项则是之前的 content
	// bool
	seeking(origin string, toLower string) (seekResult, string, bool)
}

type seek struct{}

func (s *seek) replaceKeyPoint(content, key string) string {
	return strings.ReplaceAll(content, key, _placeholder)
}

func newSeeker(index int, originKey, key string, tags map[string]string, config *matcherConfig) (seeker, int) {
	newTags := maps.Clone(tags)
	if config.enableFuzzy && strings.Index(key, config.fuzzyConfig.Sep) > -1 {
		return newFuzzy(index, originKey, key, newTags, config.fuzzyConfig), seekFuzzy
	} else {
		return newAccurate(index, originKey, key, newTags), seekAccurate
	}
}
