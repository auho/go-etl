package match

import (
	"fmt"
	"strings"
)

var _placeholder = fmt.Sprintf("%c", 0x00)

const (
	seekAccurate = iota
	seekFuzzy
)

type seekResult struct {
	key    string            // origin key
	labels map[string]string // 匹配项 labels
	amount int
}

type seekResults []seekResult

type seeker interface {
	// seekResult: seek result
	// string: 去除匹配项后的 content，如果未匹配项则是之前的 content
	// bool
	seeking(string) (seekResult, string, bool)
}

type seek struct {
}

func (s *seek) replaceKeyPoint(content, key string) string {
	return strings.ReplaceAll(content, key, _placeholder)
}

func newSeeker(index int, originKey, key string, labels map[string]string, config FuzzyConfig) (seeker, int) {
	if config.Enable && strings.Index(key, config.Sep) > -1 {
		return newFuzzyFromItem(index, originKey, key, labels, config), seekFuzzy
	} else {
		return newAccurate(index, originKey, key, labels), seekAccurate
	}
}
