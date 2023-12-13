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
	keysWidth int               // 所有词的宽度
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
func (f *fuzzy) seeking(origin, content string) (seekResult, string, bool) {
	result := newSeekResult()
	result.keyword = f.originKey
	result.tags = f.tags

	var matchedIndex, textLen int // 每次 matched 的结束 index
	var matchedContent, text, originText, before string
	var hasMatch, ok bool

	for {
		before, text, content, ok = f.match(content)
		if ok {
			hasMatch = true

			textLen = len(text)
			matchedIndex += len(before)
			originText = origin[matchedIndex : matchedIndex+textLen]

			if _, ok1 := result.textsAmount[originText]; !ok1 {
				result.texts = append(result.texts, originText)
			}

			result.textsAmount[originText] += 1
			result.amount += 1

			// 防止重复 count
			matchedContent += before + _placeholder
			matchedIndex += textLen

			// TODO optimize 最后一次如果长度不够，不再进行 match
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
	hasMatch := true
	var before, gap, text string

	originContent := content

	for _i, _key := range f.keys {
		_index := 0
		if _key.key == "" {
			_index = 0
		} else {
			_index = strings.Index(content, _key.key)
		}

		if _index > -1 {
			if _i == 0 { // 第一个 key 取匹配项前面的部分
				before = content[0:_index]
				text = _key.key
			} else { // 第二个 key 开始计算与前面 key 的 gap
				gap = content[0:_index]
				text += gap + _key.key
			}

			content = content[_index+_key.keyLen:] // 截取匹配到的关键词的后面部分

			if _i > 0 { // 取和上一个 key 的 window
				// 使用 rune 长度
				if utf8.RuneCountInString(gap) > f.windows[_i-1] { // 匹配到的关键词的距离是否在窗口内（包含窗口）
					hasMatch = false
					break
				}
			}
		} else { // 未匹配到
			hasMatch = false
			break
		}
	}

	if hasMatch {
		return before, text, content, true
	} else {
		return originContent, "", "", false
	}
}
