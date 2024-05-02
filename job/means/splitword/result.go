package splitword

import (
	"strings"
)

const NameWord = "word"

var DefaultFormat = Format{
	WordName: NameWord,
	Sep:      " ",
}

type Format struct {
	WordName string
	Sep      string
}

func (f *Format) check() {
	if f.WordName == "" {
		f.WordName = NameWord
	}

	if f.Sep == "" {
		f.Sep = " "
	}
}

type Results []string

func (rs Results) ToAll(format Format) []map[string]any {
	var rets []map[string]any
	for _, r := range rs {
		rets = append(rets, map[string]any{
			format.WordName: r,
		})
	}

	return rets
}

func (rs Results) ToLine(format Format) []map[string]any {
	return []map[string]any{{
		format.WordName: strings.Join(rs, format.Sep),
	}}
}
