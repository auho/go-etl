package match

import (
	"fmt"

	"github.com/auho/go-etl/v2/job/means"
)

type matchSearch struct {
	rule    means.Ruler
	matcher *matcher

	fn func(means.Ruler, *matcher, []string) []map[string]any

	matcherConfig *matcherConfig
	newMatcherFun func(*matcherConfig) (*matcher, error)
}

func (s *matchSearch) Prepare() error {
	if s.newMatcherFun == nil {
		s.newMatcherFun = func(config *matcherConfig) (*matcher, error) {
			items, err := s.rule.ItemsAlias()
			if err != nil {
				return nil, fmt.Errorf("ItemsAlias error; %w", err)
			}

			return newMatcher(s.rule.KeywordNameAlias(), items, config), nil
		}
	}

	var err error
	s.matcher, err = s.newMatcherFun(s.matcherConfig)
	if err != nil {
		return fmt.Errorf("prepare error; %w", err)
	}

	return nil
}

func (s *matchSearch) Close() error { return nil }
