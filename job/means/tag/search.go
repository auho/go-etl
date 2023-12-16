package tag

import (
	"fmt"

	"github.com/auho/go-etl/v2/job/means"
)

type tagSearch struct {
	rule    means.Ruler
	Matcher *Matcher
}

func (s *tagSearch) Prepare() error {
	s.Matcher = DefaultMatcher()

	items, err := s.rule.ItemsForRegexp()
	if err != nil {
		return fmt.Errorf("ItemsForRegexp error; %w", err)
	}

	s.Matcher.prepare(s.rule.KeywordNameAlias(), items, s.rule.FixedAlias())

	return nil
}

func (s *tagSearch) Close() error { return nil }
