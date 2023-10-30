package table

import (
	"unicode/utf8"
)

type Tool struct {
}

func (t Tool) stringUtf8Len(s string) int {
	return utf8.RuneCountInString(s)
}
