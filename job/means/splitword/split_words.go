package splitword

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/explore/search"
)

var _ search.Searcher = (*SplitWords)(nil)

type SplitWords struct {
	sep    string
	export *Export
}

func NewDefault(sep string) *SplitWords {
	return NewSplitWords(sep, NewExportAll())
}

func NewSplitWords(sep string, export *Export) *SplitWords {
	return &SplitWords{sep: sep, export: export}
}

func (s *SplitWords) GetTitle() string {
	return fmt.Sprintf("SplitWords[%s]", s.sep)
}

func (s *SplitWords) GenExport() search.Exporter {
	return s.export
}

func (s *SplitWords) Prepare() error {
	s.export.format.check()

	return nil
}

func (s *SplitWords) Do(contents []string) search.Token {
	var results Results
	for _, c := range contents {
		rets := strings.Split(c, s.sep)
		results = append(results, rets...)
	}

	return s.export.ToToken(results)
}

func (s *SplitWords) Close() error { return nil }
