package tag

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/job/extract"
)

var _ extract.Extractor = (*Matcher[results])(nil)
var _ extract.Extractor = (*Matcher[labelResults])(nil)
var _ extract.Extractor = (*MatcherResults)(nil)
var _ extract.Extractor = (*MatcherLabelResults)(nil)

// MatcherResults is a pre-defined alias for Matcher[results].
type MatcherResults = Matcher[results]

// MatcherLabelResults is a pre-defined alias for Matcher[labelResults].
type MatcherLabelResults = Matcher[labelResults]

// matcherContextResults and matcherContextLabelResults are type aliases
// used internally to reduce generic boilerplate.
type matcherContextResults = matcherContext[results]
type matcherContextLabelResults = matcherContext[labelResults]

// resultsEntity is a type constraint that accepts both results and labelResults.
type resultsEntity interface {
	results | labelResults
}

// matcherResultsFunc is the scanning strategy: given a context and content,
// it returns the matched results T.
type matcherResultsFunc[T resultsEntity] func(*matcherContext[T], []string) T

// matcherContext holds the scanner instance used during extraction.
type matcherContext[T resultsEntity] struct {
	scanner *scanner
}

// Matcher is a generic extractor that drives a scanner to find keywords/labels
// in content via regex, then transforms the results into structured rows.
//
// T is either results (keyword-oriented) or labelResults (label-oriented).
type Matcher[T resultsEntity] struct {
	scanner *scanner

	toMaps    func(T, extract.Rule, Format) []map[string]any
	rule      extract.Rule
	format    Format
	keys      []string
	defaults  map[string]any
	pluckKeys []string

	context           *matcherContext[T]
	matcherResultsFun matcherResultsFunc[T]

	scannerConfig scannerConfig
	newScannerFun func(extract.Rule, scannerConfig) (*scanner, error)
}

func newMatcher[T resultsEntity](
	rule extract.Rule,
	toMaps func(T, extract.Rule, Format) []map[string]any,
	keys []string,
	defaults map[string]any,
	fn matcherResultsFunc[T],
) *Matcher[T] {
	return &Matcher[T]{
		rule:              rule,
		format:            defaultFormat,
		toMaps:            toMaps,
		keys:              keys,
		defaults:          defaults,
		matcherResultsFun: fn,
		scannerConfig:     scannerConfig{},
	}
}

// Title returns a human-readable identifier for the matcher.
func (s *Matcher[T]) Title() string {
	return fmt.Sprintf("Matcher{%s:%s}", s.rule.Name(), strings.Join(s.Keys(), ","))
}

// Keys returns the column names for extraction results.
// If WithPluckKeys was called, only plucked keys are returned.
func (s *Matcher[T]) Keys() []string {
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

// DefaultValues returns default values for each key.
// If WithPluckKeys was called, only plucked key defaults are returned.
func (s *Matcher[T]) DefaultValues() map[string]any {
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

// Extract runs the scanner on the given contents and converts results to rows.
func (s *Matcher[T]) Extract(contents []string) extract.Result {
	rets := s.matcherResultsFun(s.context, contents)
	if len(rets) == 0 {
		return extract.Result{}
	}

	rows := s.toMaps(rets, s.rule, s.format)
	if len(s.pluckKeys) > 0 {
		rows = extract.PluckRows(rows, s.pluckKeys)
	}

	return extract.NewResult(true, rows)
}

func (s *Matcher[T]) Prepare() error {
	if s.newScannerFun == nil {
		s.newScannerFun = defaultScanner
	}

	var err error
	s.scanner, err = s.newScannerFun(s.rule, s.scannerConfig)
	if err != nil {
		return fmt.Errorf("newScannerFun: %w", err)
	}

	s.context = &matcherContext[T]{
		scanner: s.scanner,
	}

	return nil
}

// Close releases resources. Currently a no-op.
func (s *Matcher[T]) Close() error { return nil }

// WithPluckKeys restricts output rows to the specified keys only.
func (s *Matcher[T]) WithPluckKeys(keys []string) *Matcher[T] {
	s.pluckKeys = keys

	return s
}
