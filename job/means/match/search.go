package match

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Searcher = (*Search[Results])(nil)
var _ search.Searcher = (*Search[LabelResults])(nil)
var _ search.Searcher = (*SearchResults)(nil)
var _ search.Searcher = (*SearchLabelResults)(nil)

type SearchResults = Search[Results]
type SearchLabelResults = Search[LabelResults]
type SearchContextResults = SearchContext[Results]
type SearchContextLabelResults = SearchContext[LabelResults]

type ResultsEntity interface {
	Results | LabelResults
}

type SearchResultsFun[T ResultsEntity] func(*SearchContext[T], []string) T
type SearchContext[T ResultsEntity] struct {
	Matcher *matcher
}

type Search[T ResultsEntity] struct {
	matcher *matcher
	export  *Export[T]

	context          *SearchContext[T]
	searchResultsFun SearchResultsFun[T]

	matcherConfig *matcherConfig
	newMatcherFun func(means.Ruler, *matcherConfig) (*matcher, error)
}

func NewSearch[T ResultsEntity](export *Export[T], fn SearchResultsFun[T]) *Search[T] {
	return &Search[T]{
		export:           export,
		searchResultsFun: fn,
		matcherConfig:    &matcherConfig{},
	}
}

func (s *Search[T]) GetTitle() string {
	return fmt.Sprintf("Search{%s:%s}", s.export.GetRule().Name(), strings.Join(s.export.GetKeys(), ","))
}

func (s *Search[T]) GenExport() search.Exporter {
	return s.export
}

func (s *Search[T]) Do(contents []string) search.Token {
	rets := s.searchResultsFun(s.context, contents)

	return s.export.ToToken(rets)
}

func (s *Search[T]) Prepare() error {
	if s.newMatcherFun == nil {
		s.newMatcherFun = defaultMatcher
	}

	var err error
	s.matcher, err = s.newMatcherFun(s.export.GetRule(), s.matcherConfig)
	if err != nil {
		return fmt.Errorf("prepare error; %w", err)
	}

	s.context = &SearchContext[T]{
		Matcher: s.matcher,
	}

	return nil
}

func (s *Search[T]) Close() error { return nil }

func (s *Search[T]) ToMeans() *means.Means {
	return means.NewMeans(s)
}

func (s *Search[T]) WithPluck(keys []string) *Search[T] {
	s.export.Pluck(keys)

	return s
}

func (s *Search[T]) WithIgnoreCase() *Search[T] {
	s.matcherConfig.ignoreCase = true

	return s
}

func (s *Search[T]) WithFuzzy(config FuzzyConfig) *Search[T] {
	s.matcherConfig.enableFuzzy = true
	s.matcherConfig.fuzzyConfig = config

	return s
}

func (s *Search[T]) WithPriorityFuzzy() *Search[T] {
	s.matcherConfig.mode = modePriorityFuzzy

	return s
}

func (m *Search[T]) WithMatcher(keyName string, items []map[string]string) *Search[T] {
	m.newMatcherFun = func(rule means.Ruler, config *matcherConfig) (*matcher, error) {
		return newMatcher(keyName, items, config), nil
	}

	return m
}

func (s *Search[T]) WithDebug() *Search[T] {
	s.matcherConfig.debug = true

	return s
}
