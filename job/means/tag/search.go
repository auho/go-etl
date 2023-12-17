package tag

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Searcher = (*Search[Results])(nil)
var _ search.Searcher = (*Search[LabelResults])(nil)

type SearchEntity interface {
	Results | LabelResults
}

type GenExport[T SearchEntity] func(T, means.Ruler) search.Exporter
type SearchToExport[T SearchEntity] func(*SearchContext[T], []string) T
type SearchContext[T SearchEntity] struct {
	Matcher *matcher
}

type Search[T SearchEntity] struct {
	rule    means.Ruler
	matcher *matcher

	context          *SearchContext[T]
	genExportFn      GenExport[T]
	searchToExportFn SearchToExport[T]

	matcherConfig *matcherConfig
	newMatcherFun func(means.Ruler, *matcherConfig) (*matcher, error)
}

func NewSearch[T SearchEntity](rule means.Ruler, gek GenExport[T], fn SearchToExport[T]) *Search[T] {
	return &Search[T]{
		rule:             rule,
		genExportFn:      gek,
		searchToExportFn: fn,
	}
}

func (s *Search[T]) GetTitle() string {
	return fmt.Sprintf("Search{%s}", strings.Join(s.GenExport().GetKeys(), ","))
}

func (s *Search[T]) GenExport() search.Exporter {
	return s.genExportFn(nil, s.rule)
}

func (s *Search[T]) Do(contents []string) search.Exporter {
	rets := s.searchToExportFn(s.context, contents)

	return s.genExportFn(rets, s.rule)
}

func (s *Search[T]) Prepare() error {
	if s.newMatcherFun == nil {
		s.newMatcherFun = defaultMatcher
	}

	var err error
	s.matcher, err = s.newMatcherFun(s.rule, s.matcherConfig)
	if err != nil {
		return fmt.Errorf("prepare error; %w", err)
	}

	s.context = &SearchContext[T]{
		Matcher: s.matcher,
	}

	return nil
}

func (s *Search[T]) Close() error { return nil }
