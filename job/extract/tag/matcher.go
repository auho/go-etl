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

// resultsEntity is a type constraint that accepts both results and labelResults.
type resultsEntity interface {
	results | labelResults
}

// resultsFunc is the scanning strategy: given a scanner and content,
// it returns the matched results T.
type resultsFunc[T resultsEntity] func(*scanner, []string) T

// Matcher is a generic extractor that drives a scanner to find keywords/labels
// in content via regex, then transforms the results into structured rows.
//
// T is either results (keyword-oriented) or labelResults (label-oriented).
type Matcher[T resultsEntity] struct {
	scanner *scanner

	rule      extract.Rule
	format    Format
	keys      []string
	defaults  map[string]any
	pluckKeys []string

	rowsFunc      func(T, extract.Rule, Format) []map[string]any
	resultsFunc   resultsFunc[T]
	scannerFunc   func(ScannerConfig) (*scanner, error)
	scannerConfig ScannerConfig
}

func newMatcher[T resultsEntity](
	rule extract.Rule,
	rowsFunc func(T, extract.Rule, Format) []map[string]any,
	keys []string,
	defaults map[string]any,
	fn resultsFunc[T],
) *Matcher[T] {
	return &Matcher[T]{
		rule:          rule,
		format:        defaultFormat,
		rowsFunc:      rowsFunc,
		keys:          keys,
		defaults:      defaults,
		resultsFunc:   fn,
		scannerConfig: ScannerConfig{},
	}
}

func (m *Matcher[T]) Prepare() error {
	if m.scanner == nil {
		if m.scannerFunc == nil {
			m.scannerFunc = func(sc ScannerConfig) (*scanner, error) {
				return defaultScannerFromRule(m.rule, WithScannerConfig(sc))
			}
		}

		var err error
		m.scanner, err = m.scannerFunc(m.scannerConfig)
		if err != nil {
			return fmt.Errorf("scannerFunc: %w", err)
		}
	}

	if len(m.pluckKeys) > 0 {
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

// Title returns a human-readable identifier for the matcher.
func (m *Matcher[T]) Title() string {
	return fmt.Sprintf("Matcher{%s:%s}", m.rule.Name(), strings.Join(m.Keys(), ","))
}

// Keys returns the column names for extraction results.
func (m *Matcher[T]) Keys() []string {
	return m.keys
}

// DefaultValues returns default values for each key.
func (m *Matcher[T]) DefaultValues() map[string]any {
	return m.defaults
}

// Extract runs the scanner on the given contents and converts results to rows.
func (m *Matcher[T]) Extract(contents []string) extract.Result {
	rets := m.resultsFunc(m.scanner, contents)
	if len(rets) == 0 {
		return extract.Result{}
	}

	rows := m.rowsFunc(rets, m.rule, m.format)
	if len(m.pluckKeys) > 0 {
		rows = extract.PluckRows(rows, m.pluckKeys)
	}

	return extract.NewResult(true, rows)
}

// Close releases resources. Currently a no-op.
func (m *Matcher[T]) Close() error { return nil }

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

// WithDebug enables debug output during scanning.
func (m *Matcher[T]) WithDebug() *Matcher[T] {
	m.scannerConfig.Debug = true

	return m
}

// WithRuleScanner replaces the default scanner factory with one that builds
// the scanner from the rule using the given ScannerOptions.
func (m *Matcher[T]) WithRuleScanner(opts ...ScannerOption) *Matcher[T] {
	m.scannerFunc = func(sc ScannerConfig) (*scanner, error) {
		opts = append([]ScannerOption{WithScannerConfig(sc)}, opts...)
		return newScannerFromRule(m.rule, opts...)
	}

	return m
}
