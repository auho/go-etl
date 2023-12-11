package tag

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Searcher = (*SearchLabel)(nil)

type SearchLabel struct {
	tagSearch
	genExport GenExportLabel
	fn        func(*SearchLabel, []string) search.Exporter
}

func NewSearchLabel(rule means.Ruler, gel GenExportLabel, fn func(*SearchLabel, []string) search.Exporter) *SearchLabel {
	return &SearchLabel{
		tagSearch: tagSearch{rule: rule},
		genExport: gel,
		fn:        fn,
	}
}

func (s *SearchLabel) GetTitle() string {
	return fmt.Sprintf("SearchLabel{%s}", strings.Join(s.GenExport().GetKeys(), ","))
}

func (s *SearchLabel) GenExport() search.Exporter {
	return s.genExport(nil, s.rule)
}

func (s *SearchLabel) Do(contents []string) search.Exporter {
	return s.fn(s, contents)
}

func newSearchLabel(rule means.Ruler, gel GenExportLabel, genDoFn func(*SearchLabel) func([]string) LabelResults) *SearchLabel {
	return NewSearchLabel(rule, gel, func(s *SearchLabel, contents []string) search.Exporter {
		rets := genDoFn(s)(contents)
		if rets == nil {
			return nil
		}

		return s.genExport(rets, rule)
	})
}

func NewSearchWholeLabels(rule means.Ruler) *SearchLabel {
	return newSearchLabel(rule, NewExportLabelLine, func(s *SearchLabel) func([]string) LabelResults {
		return s.matcher.MatchLabel
	})
}

func NewSearchLabels(rule means.Ruler, gel GenExportLabel) *SearchLabel {
	return newSearchLabel(rule, gel, func(s *SearchLabel) func([]string) LabelResults {
		return s.matcher.MatchLabel
	})
}
