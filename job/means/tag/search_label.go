package tag

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

type SearchLabel struct {
	rule      means.Ruler
	newExport NewExportLabel
}

func (s *SearchLabel) GetTitle() string {
	//TODO implement me
	panic("implement me")
}

func (s *SearchLabel) GetExport() search.Exporter {
	//TODO implement me
	panic("implement me")
}

func (s *SearchLabel) Do(contents []string) search.Exporter {
	//TODO implement me
	panic("implement me")
}
