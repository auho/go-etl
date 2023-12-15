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

func (a *accurate) seeking(origin, content string) (seekResult, seekContent, bool) {
	result := newSeekResult()
	result.keyword = a.originKey
	result.tags = a.tags

	var matchedIndex, beforeLen int // 每次 matched 的结束 index
	var matchedOrigin, matchedContent, matchedText, before string
	var hasMatch, ok bool

	keyLen := len(a.key)
	for {
		before, content, ok = strings.Cut(content, a.key)
		beforeLen = len(before)

		matchedContent += before
		matchedOrigin += origin[matchedIndex : matchedIndex+beforeLen]

		if ok {
			hasMatch = true

			matchedContent += _placeholder
			matchedOrigin += _placeholder

			matchedIndex += beforeLen
			matchedText = origin[matchedIndex : matchedIndex+keyLen]

			result.texts = append(result.texts, textResult{
				text:  matchedText,
				start: matchedIndex,
				width: keyLen,
			})

			// + key
			matchedIndex += keyLen

			result.textsAmount[matchedText] += 1
			result.amount += 1

		} else {
			break
		}
	}

	if hasMatch {
		return result, seekContent{
			origin:  matchedOrigin,
			content: matchedContent,
		}, true
	} else {
		return seekResult{}, seekContent{
			origin:  matchedOrigin,
			content: matchedContent,
		}, false
	}
}
