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

func newAccurate(index int, originKey, key string, tags map[string]string) *accurate {
	return &accurate{
		index:     index,
		originKey: originKey,
		key:       key,
		tags:      tags,
	}
}

func (a *accurate) seeking(origin, content string) (seekResult, string, bool) {
	result := newSeekResult()
	result.keyword = a.originKey
	result.tags = a.tags

	var matchedIndex int // 每次 matched 的结束 index
	var matchedContent, originText, before string
	var hasMatch, ok bool

	keyLen := len(a.key)
	for {
		before, content, ok = strings.Cut(content, a.key)
		if ok {
			hasMatch = true

			matchedIndex += len(before)
			originText = origin[matchedIndex : matchedIndex+keyLen]

			if _, ok1 := result.textsAmount[originText]; !ok1 {
				result.texts = append(result.texts, originText)
			}

			result.textsAmount[originText] += 1
			result.amount += 1

			// 防止重复 count
			matchedContent += before + _placeholder
			matchedIndex += keyLen
		} else {
			matchedContent += before

			break
		}
	}

	if hasMatch {
		return result, matchedContent, true
	} else {
		return seekResult{}, matchedContent, false
	}
}
