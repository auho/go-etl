package match

import (
	"strings"
)

var _ seeker = (*accurate)(nil)

type accurate struct {
	seek
	index     int
	originKey string            // origin keyword
	key       string            // if ignore case, all to lower
	tags      map[string]string // tags name and value
}

func newAccurate(index int, originKey, key string, labels map[string]string) *accurate {
	return &accurate{
		index:     index,
		originKey: originKey,
		key:       key,
		tags:      labels,
	}
}

func (a *accurate) seeking(content string) (seekResult, string, bool) {
	_count := strings.Count(content, a.key)
	if _count > 0 {
		// replace 防止重复 count
		content = a.replaceKeyPoint(content, a.key)

		result := newSeekResult()
		result.keyword = a.originKey
		result.texts = append(result.texts, a.key)
		result.textsAmount[a.key] = _count
		result.amount = _count

		return result, content, true
	} else {
		return seekResult{}, content, false
	}
}
