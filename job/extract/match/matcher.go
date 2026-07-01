package match

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
// in content, then transforms the results into structured rows.
//
// T is either results (keyword-oriented) or labelResults (label-oriented).
type Matcher[T resultsEntity] struct {
	scanner *scanner

	rule    extract.Rule
	format  Format
	toMaps  func(T, extract.Rule, Format) []map[string]any
	keysFun func(extract.Rule) ([]string, map[string]any)

	keys      []string
	defaults  map[string]any
	pluckKeys []string

	context           *matcherContext[T]
	matcherResultsFun matcherResultsFunc[T]

	scannerConfig ScannerConfig
	newScannerFun func(extract.Rule, ScannerConfig) (*scanner, error)
}

func newMatcher[T resultsEntity](rule extract.Rule, toMaps func(T, extract.Rule, Format) []map[string]any, keysFun func(extract.Rule) ([]string, map[string]any), fn matcherResultsFunc[T]) *Matcher[T] {
	return &Matcher[T]{
		rule:              rule,
		format:            defaultFormat,
		toMaps:            toMaps,
		keysFun:           keysFun,
		matcherResultsFun: fn,
		scannerConfig:     ScannerConfig{},
	}
}

// Title returns a human-readable identifier for the matcher.
func (m *Matcher[T]) Title() string {
	return fmt.Sprintf("Matcher{%s:%s}", m.rule.Name(), strings.Join(m.keys, ","))
}

// Extract runs the scanner on the given contents and converts results to rows.
func (m *Matcher[T]) Extract(contents []string) extract.Result {
	rets := m.matcherResultsFun(m.context, contents)
	if len(rets) == 0 {
		return extract.Result{}
	}

	rows := m.toMaps(rets, m.rule, m.format)
	if m.pluckKeys != nil {
		rows = extract.PluckRows(rows, m.pluckKeys)
	}

	return extract.NewResult(true, rows)
}

// Prepare initializes the scanner via newScannerFun and computes keys/defaults.
// Must be called before Extract.
func (m *Matcher[T]) Prepare() error {
	if m.newScannerFun == nil {
		m.newScannerFun = defaultScanner
	}

	var err error
	m.scanner, err = m.newScannerFun(m.rule, m.scannerConfig)
	if err != nil {
		return fmt.Errorf("newScannerFun: %w", err)
	}

	m.context = &matcherContext[T]{
		scanner: m.scanner,
	}

	m.keys, m.defaults = m.keysFun(m.rule)
	if m.pluckKeys != nil {
		pluckedKeys := make([]string, 0, len(m.pluckKeys))
		pluckedDefaults := make(map[string]any)
		for _, k := range m.pluckKeys {
			if v, ok := m.defaults[k]; ok {
				pluckedKeys = append(pluckedKeys, k)
				pluckedDefaults[k] = v
			}
		}
		m.keys = pluckedKeys
		m.defaults = pluckedDefaults
	}

	return nil
}

// Close releases resources. Currently a no-op.
func (m *Matcher[T]) Close() error { return nil }

// Keys returns the column names for extraction results.
func (m *Matcher[T]) Keys() []string {
	return m.keys
}

// DefaultValues returns default values for each key.
func (m *Matcher[T]) DefaultValues() map[string]any {
	return m.defaults
}

// WithPluckKeys restricts output rows to the specified keys only.
func (m *Matcher[T]) WithPluckKeys(keys []string) *Matcher[T] {
	m.pluckKeys = keys

	return m
}

// WithFormat sets a custom output format (keyword amount suffix and separator).
func (m *Matcher[T]) WithFormat(f Format) *Matcher[T] {
	m.format = f

	return m
}

// WithScannerConfig replaces the default ScannerConfig entirely.
func (m *Matcher[T]) WithScannerConfig(sc ScannerConfig) *Matcher[T] {
	m.scannerConfig = sc

	return m
}

// WithIgnoreCase enables case-insensitive keyword matching.
func (m *Matcher[T]) WithIgnoreCase() *Matcher[T] {
	m.scannerConfig.IgnoreCase = true

	return m
}

// WithModePriorityAccurate sets the scanner to try accurate (exact) matches first,
// falling back to fuzzy matches only when no accurate match is found.
func (m *Matcher[T]) WithModePriorityAccurate() *Matcher[T] {
	m.scannerConfig.Mode = modePriorityAccurate

	return m
}

// WithModePriorityFuzzy sets the scanner to try fuzzy matches first,
// falling back to accurate matches only when no fuzzy match is found.
func (m *Matcher[T]) WithModePriorityFuzzy() *Matcher[T] {
	m.scannerConfig.Mode = modePriorityFuzzy

	return m
}

// WithFuzzy enables fuzzy matching with the given FuzzyConfig.
// It marks enabled=true internally so callers don't need to set it.
func (m *Matcher[T]) WithFuzzy(fc FuzzyConfig) *Matcher[T] {
	fc.enabled = true
	m.scannerConfig.Fuzzy = fc

	return m
}

// WithDebug enables debug output during scanning.
func (m *Matcher[T]) WithDebug() *Matcher[T] {
	m.scannerConfig.Debug = true

	return m
}

// WithScanner configures the matcher to use a custom key-name and items list
// instead of the default rule-based scanner.
func (m *Matcher[T]) WithScanner(keyName string, items []map[string]string) *Matcher[T] {
	m.newScannerFun = func(rule extract.Rule, sc ScannerConfig) (*scanner, error) {
		return newScanner(keyName, items, sc), nil
	}

	return m
}
