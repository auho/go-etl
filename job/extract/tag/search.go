package tag

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/job/extract"
)

var _ extract.Extractor = (*search[results])(nil)
var _ extract.Extractor = (*search[labelResults])(nil)
var _ extract.Extractor = (*searchResults)(nil)
var _ extract.Extractor = (*searchLabelResults)(nil)

type searchResults = search[results]
type searchLabelResults = search[labelResults]
type searchContextResults = searchContext[results]
type searchContextLabelResults = searchContext[labelResults]

type resultsEntity interface {
	results | labelResults
}

type searchResultsFunc[T resultsEntity] func(*searchContext[T], []string) T
type searchContext[T resultsEntity] struct {
	matcher *matcher
}

type search[T resultsEntity] struct {
	matcher *matcher

	toMaps    func(T, extract.Rule, format) []map[string]any
	rule      extract.Rule
	format    format
	keys      []string
	defaults  map[string]any
	pluckKeys []string

	context          *searchContext[T]
	searchResultsFun searchResultsFunc[T]

	matcherConfig *matcherConfig
	newMatcherFun func(extract.Rule, *matcherConfig) (*matcher, error)
}

func newSearch[T resultsEntity](
	rule extract.Rule,
	toMaps func(T, extract.Rule, format) []map[string]any,
	keys []string,
	defaults map[string]any,
	fn searchResultsFunc[T],
) *search[T] {
	return &search[T]{
		rule:             rule,
		format:           defaultFormat,
		toMaps:           toMaps,
		keys:             keys,
		defaults:         defaults,
		searchResultsFun: fn,
		matcherConfig:    &matcherConfig{},
	}
}

func (s *search[T]) Title() string {
	return fmt.Sprintf("Search{%s:%s}", s.rule.Name(), strings.Join(s.Keys(), ","))
}

func (s *search[T]) Keys() []string {
	if len(s.pluckKeys) == 0 {
		return s.keys
	}

	keySet := make(map[string]bool, len(s.keys))
	for _, k := range s.keys {
		keySet[k] = true
	}

	keys := make([]string, 0, len(s.pluckKeys))
	for _, pk := range s.pluckKeys {
		if keySet[pk] {
			keys = append(keys, pk)
		}
	}

	return keys
}

func (s *search[T]) DefaultValues() map[string]any {
	if len(s.pluckKeys) == 0 {
		return s.defaults
	}

	pluckSet := make(map[string]bool, len(s.pluckKeys))
	for _, k := range s.pluckKeys {
		pluckSet[k] = true
	}

	dv := make(map[string]any)
	for k, v := range s.defaults {
		if pluckSet[k] {
			dv[k] = v
		}
	}

	return dv
}

func (s *search[T]) Extract(contents []string) extract.Result {
	results := s.searchResultsFun(s.context, contents)
	if len(results) == 0 {
		return extract.Result{}
	}

	rows := s.toMaps(results, s.rule, s.format)
	if len(s.pluckKeys) > 0 {
		rows = extract.PluckRows(rows, s.pluckKeys)
	}

	return extract.NewResult(true, rows)
}

func (s *search[T]) Prepare() error {
	if s.newMatcherFun == nil {
		s.newMatcherFun = defaultMatcher
	}

	var err error
	s.matcher, err = s.newMatcherFun(s.rule, s.matcherConfig)
	if err != nil {
		return fmt.Errorf("newMatcherFun: %w", err)
	}

	s.context = &searchContext[T]{
		matcher: s.matcher,
	}

	return nil
}

func (s *search[T]) Close() error { return nil }

func (s *search[T]) WithPluck(keys []string) *search[T] {
	s.pluckKeys = keys

	return s
}
