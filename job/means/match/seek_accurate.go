package match

import (
	"strings"
)

var _ seeker = (*accurate)(nil)

type accurate struct {
	seek
	index     int
	originKey string
	key       string
	labels    map[string]string
}

func newAccurate(index int, originKey, key string, labels map[string]string) *accurate {
	return &accurate{
		index:     index,
		originKey: originKey,
		key:       key,
		labels:    labels,
	}
}

func (a *accurate) seeking(content string) (seekResult, string, bool) {
	_count := strings.Count(content, a.key)
	if _count > 0 {
		// replace 防止重复 count
		content = a.replaceKeyPoint(content, a.key)

		return seekResult{
			key:    a.originKey,
			labels: a.labels,
			amount: _count,
		}, content, true
	} else {
		return seekResult{}, content, false
	}
}
