package tag

import (
	"fmt"

	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Searcher = (*SearchKeyword)(nil)

type SearchKeyword struct {
	rule      means.Ruler
	matcher   *Matcher
	newExport NewExportKeyword
	fn        func(means.Ruler, *Matcher, []string) []map[string]any
}

func NewSearchKeyword(rule means.Ruler, ne NewExportKeyword) *SearchKeyword {
	return &SearchKeyword{rule: rule, newExport: ne}
}

func (s *SearchKeyword) Prepare() error {
	s.matcher = DefaultMatcher()

	items, err := s.rule.ItemsForRegexp()
	if err != nil {
		return fmt.Errorf("ItemsForRegexp error; %w", err)
	}

	s.matcher.prepare(s.rule.KeywordNameAlias(), items, s.rule.FixedAlias())

	return nil
}

func (s *SearchKeyword) Close() error {
	//TODO implement me
	panic("implement me")
}

func (s *SearchKeyword) GetTitle() string {
	return "SearchKeyword"
}

func (s *SearchKeyword) GenExport() search.Exporter {
	return s.newExport(nil, s.rule)
}

func (s *SearchKeyword) Do(contents []string) search.Exporter {
	for _, content := range contents {
		_ = content //TODO
	}

	return nil // TODO
}
