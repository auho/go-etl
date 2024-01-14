package splitword

import (
	"strings"
)

const NameWord = "word"

var DefaultFormat = Format{
	WorkdName: NameWord,
	Sep:       " ",
}

type Format struct {
	WorkdName string
	Sep       string
}

func (f *Format) check() {
	if f.WorkdName == "" {
		f.WorkdName = NameWord
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
			format.WorkdName: r,
		})
	}

	return rets
}

func (rs Results) ToLine(format Format) []map[string]any {
	return []map[string]any{{
		format.WorkdName: strings.Join(rs, format.Sep),
	}}
}
