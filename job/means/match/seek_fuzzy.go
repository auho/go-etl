package match

import (
	"strings"
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
	originKey string
	key       string
	labels    map[string]string
	keys      []fuzzyKey
	windows   []int // 前后两个关键词的距离
	keysWidth int   // 所有词的宽度
}

func newFuzzyFromItem(index int, originKey, key string, labels map[string]string, config FuzzyConfig) *fuzzy {
	keys := strings.Split(key, config.Sep)

	f := &fuzzy{}
	f.index = index
	f.originKey = originKey
	f.key = key
	f.labels = labels

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

func (f *fuzzy) seeking(content string) (seekResult, string, bool) {
	var ok, hasMatch bool
	var prefix string
	var amount int
	var newContent string

	for {
		prefix, content, ok = f.match(content)
		if ok {
			amount += 1
			hasMatch = true

			newContent += prefix + _placeholder

			// TODO optimize 最后一次如果长度不够，不再进行 match
		} else {
			newContent += content

			break
		}
	}

	if hasMatch {
		return seekResult{
			key:    f.originKey,
			labels: f.labels,
			amount: amount,
		}, newContent, true
	} else {
		return seekResult{}, newContent, false
	}
}

// seekingExpression
// string: 匹配项前面的部分；如未匹配项，则是全部原始 content
// string: 匹配项后面的部分；
// bool: true 有匹配项，false 没有匹配项
func (f *fuzzy) match(content string) (string, string, bool) {
	hasMatch := true
	prefix := ""

	originContent := content

	for _i, _key := range f.keys {
		_index := strings.Index(content, _key.key)
		if _index > -1 {
			if _i == 0 { // 第一个 key 取匹配项前面的部分
				prefix = content[0:_index]
			}

			content = content[_index+_key.keyLen:] // 截取匹配到的关键词的后面部分
		} else { // 未匹配到
			hasMatch = false
			break
		}

		if _i > 0 { // 取和上一个 key 的 window
			if _index > f.windows[_i-1] { // 匹配到的关键词的距离是否在窗口内（包含窗口）
				hasMatch = false
				break
			}
		}
	}

	if hasMatch {
		return prefix, content, true
	} else {
		return originContent, "", false
	}
}
