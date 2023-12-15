package match

import (
	"strings"
	"unicode/utf8"
)

var _ seeker = (*fuzzy)(nil)

type fuzzyKey struct {
	key    string
	keyLen int
}

// fuzzy
// keys: 关键词列表
// windows 窗口
// keys: 	a b	| a b c	| a b c d
// windows: 3	| 2 5	| 1 4 5
type fuzzy struct {
	seek
	index     int
	originKey string            // origin keyword
	key       string            // if ignore case, all to lower
	tags      map[string]string // tags name and value
	keys      []fuzzyKey        // split by sep
	windows   []int             // 前后两个关键词的距离
	keysWidth int               // 所有词的总宽度 byte
}

func newFuzzy(index int, originKey, key string, tags map[string]string, config FuzzyConfig) *fuzzy {
	keys := strings.Split(key, config.Sep)

	f := &fuzzy{}
	f.index = index
	f.originKey = originKey
	f.key = key
	f.tags = tags

	for _, _k := range keys {
		_kLen := len(_k)
		f.keys = append(f.keys, fuzzyKey{
			key:    _k,
			keyLen: _kLen,
		})

		f.windows = append(f.windows, config.Window)
		f.keysWidth += _kLen
	}

	return f
}

// seeking
func (f *fuzzy) seeking(origin, content string) (seekResult, seekContent, bool) {
	result := newSeekResult()
	result.keyword = f.originKey
	result.tags = f.tags

	var matchedIndex, beforeLen, textLen int // 每次 matched 的结束 index
	var matchedOrigin, matchedContent, text, matchedText, before string
	var hasMatch, ok bool

	for {
		before, text, content, ok = f.match(content)
		beforeLen = len(before)

		if ok {
			hasMatch = true

			textLen = len(text)
			matchedContent += before + _placeholder
			matchedOrigin += origin[matchedIndex:matchedIndex+beforeLen] + _placeholder
			matchedIndex += len(before)
			matchedText = origin[matchedIndex : matchedIndex+textLen]
			result.texts = append(result.texts, textResult{
				text:  matchedText,
				start: matchedIndex,
				width: textLen,
			})

			matchedIndex += textLen

			result.textsAmount[matchedText] += 1
			result.amount += 1

			if len(content) < f.keysWidth {
				break
			}
		} else {
			// 匹配失败进行回溯，只回溯第一个 byte // TODO
			// 剩余 content 长度足够进行匹配
			if len(before) >= f.keysWidth {
				matchedContent += before[0:1]
				matchedOrigin += origin[matchedIndex : matchedIndex+1]

				matchedIndex += 1
				content = before[1:]
			} else { // 待匹配长度不够所有词的总宽度

				matchedContent += before
				matchedOrigin += origin[matchedIndex:]

				break
			}
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

// seekingExpression
//
// 搜索第一个 key，在搜索第二个 key，计算两者直接的距离
// 如果不符合条件，退出
// 如果符合条件，继续搜索第三个key，计算距离，依此类推
//
// string: 匹配项前面的部分；如未匹配项，则是全部原始 content
// string: 匹配项
// string: 匹配项后面的部分；
// bool: true 有匹配项，false 没有匹配项
func (f *fuzzy) match(content string) (string, string, string, bool) {
	var before, gap, text string
	var contentIndex int

	hasMatch := true
	originContent := content

	for _i, _key := range f.keys {
		_index := 0
		if _key.key == "" {
			_index = 0
		} else {
			_index = strings.Index(content, _key.key)
		}

		if _index > -1 {
			contentIndex += _index + _key.keyLen

			if _i == 0 { // 第一个 key 取匹配项前面的部分
				before = originContent[0:_index]
				text = _key.key
			} else { // 第二个 key 开始计算与前面 key 的 gap
				gap = content[0:_index]
				text += gap + _key.key

				// 取和上一个 key 的 window，使用 rune 长度
				// 匹配到的关键词的距离是否在窗口内（包含窗口）
				if utf8.RuneCountInString(gap) > f.windows[_i-1] {
					hasMatch = false
					break
				}
			}

			// 截取匹配到的关键词的后面部分
			content = content[_index+_key.keyLen:]

		} else { // 未匹配到
			hasMatch = false
			break
		}
	}

	if hasMatch {
		var after string
		if contentIndex+1 < len(originContent) {
			after = originContent[contentIndex:]
		}

		return before, text, after, true
	} else {
		return originContent, "", "", false
	}
}
