package match

import (
	"strings"
)

var _ seeker = (*accurate)(nil)

type accurate struct {
	seek
	keyIndex  int
	originKey string            // origin keyword
	key       string            // if ignore case, all to lower
	tags      map[string]string // tags name and value
}

func newAccurate(keyIndex int, originKey, key string, tags map[string]string, config seekConfig) *accurate {
	a := &accurate{
		keyIndex:  keyIndex,
		originKey: originKey,
		key:       key,
		tags:      tags,
	}

	a.config = config

	return a
}

func (a *accurate) seeking(sc seekContent) (seekResults, seekContent, bool) {
	var results seekResults

	// matchedIndex: 每次 matched 的结束 index
	// beforeLen: 匹配项前面的内容
	// soughtNum: seek num
	var matchedIndex, beforeLen, soughtNum int
	var matchedOrigin, matchedContent, matchedText, before string
	var hasMatch, ok bool

	soughtNum = sc.maxSeekNum
	keyLen := len(a.key)
	content := sc.content
	for {
		before, content, ok = strings.Cut(content, a.key)
		beforeLen = len(before)

		matchedContent += before
		matchedOrigin += sc.origin[matchedIndex : matchedIndex+beforeLen]

		if ok {
			hasMatch = true

			matchedIndex += beforeLen
			matchedText = sc.origin[matchedIndex : matchedIndex+keyLen]
			_ph := a.matchedToPlaceholder(matchedText)
			matchedContent += _ph
			matchedOrigin += _ph

			results = append(results, seekResult{
				index:   sc.index,
				start:   matchedIndex,
				width:   keyLen,
				keyword: a.originKey,
				text:    matchedText,
				tags:    a.tags,
			})

			// + key
			matchedIndex += keyLen
			soughtNum -= 1

			if soughtNum == 0 {
				break
			}
		} else {
			break
		}
	}

	sc.origin = matchedOrigin
	sc.content = matchedContent

	if hasMatch {
		return results, sc, true
	} else {
		return nil, sc, false
	}
}
