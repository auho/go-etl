package splitword

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/job/extract"
)

var _ extract.Extractor = (*SplitWords)(nil)

type SplitWords struct {
	sep    string
	export *extract.Exporter[Results]
	format Format
}

func NewDefault(sep string) *SplitWords {
	return NewSplitWords(sep, NewExportAll())
}

func NewSplitWords(sep string, export *extract.Exporter[Results]) *SplitWords {
	return &SplitWords{sep: sep, export: export, format: DefaultFormat}
}

func (s *SplitWords) Title() string {
	return fmt.Sprintf("SplitWords[%s]", s.sep)
}

func (s *SplitWords) NewExport() extract.FieldSpec {
	return s.export
}

func (s *SplitWords) Prepare() error {
	s.format.check()

	return nil
}

func (s *SplitWords) Search(contents []string) extract.Result {
	var results Results
	for _, c := range contents {
		rets := strings.Split(c, s.sep)
		results = append(results, rets...)
	}

	return s.export.ToToken(results, len(results) > 0)
}

func (s *SplitWords) Close() error { return nil }
