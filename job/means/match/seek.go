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
	key    string            // origin keyword
	labels map[string]string // 匹配项 labels
	amount int               // keyword matched amount
}

type seekResults []seekResult

// seeker
// content，keyword 大小写在 match 已经处理过
// 这里区分大小写
type seeker interface {
	// seekResult: seek result
	// string: 去除匹配项后的 content，如果未匹配项则是之前的 content
	// bool
	seeking(string) (seekResult, string, bool)
}

type seek struct{}

func (s *seek) replaceKeyPoint(content, key string) string {
	return strings.ReplaceAll(content, key, _placeholder)
}

func newSeeker(index int, originKey, key string, labels map[string]string, config *matcherConfig) (seeker, int) {
	if config.enableFuzzy && strings.Index(key, config.fuzzyConfig.Sep) > -1 {
		return newFuzzy(index, originKey, key, labels, config.fuzzyConfig), seekFuzzy
	} else {
		return newAccurate(index, originKey, key, labels), seekAccurate
	}
}
