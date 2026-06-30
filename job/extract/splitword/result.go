package splitword

import (
	"strings"
)

const nameWord = "word"

var defaultFormat = format{
	wordName: nameWord,
	sep:      " ",
}

type format struct {
	wordName string
	sep      string
}

func (f *format) check() {
	if f.wordName == "" {
		f.wordName = nameWord
	}

	if f.sep == "" {
		f.sep = " "
	}
}

type results []string

func (rs results) toAll(format format) []map[string]any {
	var rets []map[string]any
	for _, r := range rs {
		rets = append(rets, map[string]any{
			format.wordName: r,
		})
	}

	return rets
}

func (rs results) toLine(format format) []map[string]any {
	return []map[string]any{{
		format.wordName: strings.Join(rs, format.sep),
	}}
}
