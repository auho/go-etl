package splitword

import (
	"strings"
)

const WordName = "word"

var DefaultFormat = Format{
	Name: WordName,
	Sep:  " ",
}

type Format struct {
	Name string
	Sep  string
}

func (f *Format) check() {
	if f.Name == "" {
		f.Name = WordName
	}
}

type Results []string

func (rs Results) ToAll(format Format) []map[string]any {
	var rets []map[string]any
	for _, r := range rs {
		rets = append(rets, map[string]any{
			format.Name: r,
		})
	}

	return rets
}

func (rs Results) ToLine(format Format) []map[string]any {
	return []map[string]any{{
		format.Name: strings.Join(rs, format.Sep),
	}}
}
