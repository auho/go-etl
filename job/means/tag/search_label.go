package tag

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Searcher = (*SearchLabel)(nil)

type SearchLabel struct {
	rule      means.Ruler
	newExport NewExportLabel
}

func (s *SearchLabel) Prepare() error {
	//TODO implement me
	panic("implement me")
}

func (s *SearchLabel) Close() error {
	//TODO implement me
	panic("implement me")
}

func (s *SearchLabel) GetTitle() string {
	//TODO implement me
	panic("implement me")
}

func (s *SearchLabel) GenExport() search.Exporter {
	//TODO implement me
	panic("implement me")
}

func (s *SearchLabel) Do(contents []string) search.Exporter {
	//TODO implement me
	panic("implement me")
}
