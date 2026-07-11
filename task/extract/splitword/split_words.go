package splitword

import (
	"fmt"
	"maps"
	"strings"

	"github.com/auho/go-etl/v3/task/extract"
)

var _ extract.Extractor = (*SplitWords)(nil)

type SplitWords struct {
	sep      string
	format   format
	toRows   func(results, format) []map[string]any
	keys     []string
	defaults map[string]any
}

func NewSplitWordsAll(sep string) *SplitWords {
	return &SplitWords{
		sep:    sep,
		format: defaultFormat,
		toRows: func(r results, f format) []map[string]any { return r.toAll(f) },
	}
}

func NewSplitWordsLine(sep string) *SplitWords {
	return &SplitWords{
		sep:    sep,
		format: defaultFormat,
		toRows: func(r results, f format) []map[string]any { return r.toLine(f) },
	}
}

func (s *SplitWords) Title() string {
	return fmt.Sprintf("SplitWords[%s]", s.sep)
}

func (s *SplitWords) Prepare() error {
	s.format.check()
	s.defaults = map[string]any{s.format.wordName: ""}
	s.keys = []string{s.format.wordName}
	return nil
}

func (s *SplitWords) Keys() []string {
	return s.keys
}

func (s *SplitWords) DefaultValues() map[string]any {
	return maps.Clone(s.defaults)
}

func (s *SplitWords) Extract(contents []string) extract.Result {
	var rets results
	for _, c := range contents {
		rets = append(rets, strings.Split(c, s.sep)...)
	}
	if len(rets) == 0 {
		return extract.Result{}
	}
	return extract.NewResult(true, s.toRows(rets, s.format))
}

func (s *SplitWords) Close() error { return nil }
