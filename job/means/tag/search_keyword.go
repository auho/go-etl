package tag

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Searcher = (*SearchKeyword)(nil)
var _ search.Searcher = (*SearchLabel)(nil)

type tagSearch struct {
}

type SearchKeyword struct {
	rule      means.Ruler
	newExport NewExportKeyword
}

func NewSearchKeyword(rule means.Ruler, ne NewExportKeyword) *SearchKeyword {
	return &SearchKeyword{rule: rule, newExport: ne}
}

func (s *SearchKeyword) GetTitle() string {
	return "SearchKeyword"
}

func (s *SearchKeyword) GetExport() search.Exporter {
	return s.newExport(nil, s.rule)
}

func (s *SearchKeyword) Do(contents []string) search.Exporter {
	for _, content := range contents {
		_ = content //TODO
	}

	return nil // TODO
}
