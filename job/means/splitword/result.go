package splitword

import (
	"strings"
)

const WordName = "word"

var DefaultFormat = Format{
	Sep: " ",
}

type Format struct {
	Sep string
}

type Results []string

func (rs Results) ToAll() []map[string]any {
	var rets []map[string]any
	for _, r := range rs {
		rets = append(rets, map[string]any{
			WordName: r,
		})
	}

	return rets
}

func (rs Results) ToLine(format Format) []map[string]any {
	return []map[string]any{{
		WordName: strings.Join(rs, format.Sep),
	}}
}
