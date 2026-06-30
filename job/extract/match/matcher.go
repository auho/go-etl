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

type MatcherResults = Matcher[results]
type MatcherLabelResults = Matcher[labelResults]
type matcherContextResults = matcherContext[results]
type matcherContextLabelResults = matcherContext[labelResults]

type resultsEntity interface {
	results | labelResults
}

type matcherResultsFunc[T resultsEntity] func(*matcherContext[T], []string) T
type matcherContext[T resultsEntity] struct {
	scanner *scanner
}

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

func (m *Matcher[T]) Title() string {
	return fmt.Sprintf("Matcher{%s:%s}", m.rule.Name(), strings.Join(m.keys, ","))
}

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

func (m *Matcher[T]) Close() error { return nil }

func (m *Matcher[T]) Keys() []string {
	return m.keys
}

func (m *Matcher[T]) DefaultValues() map[string]any {
	return m.defaults
}

func (m *Matcher[T]) WithPluckKeys(keys []string) *Matcher[T] {
	m.pluckKeys = keys

	return m
}

func (m *Matcher[T]) WithFormat(f Format) *Matcher[T] {
	m.format = f

	return m
}

func (m *Matcher[T]) WithScannerConfig(sc ScannerConfig) *Matcher[T] {
	m.scannerConfig = sc

	return m
}

func (m *Matcher[T]) WithIgnoreCase() *Matcher[T] {
	m.scannerConfig.IgnoreCase = true

	return m
}

func (m *Matcher[T]) WithModePriorityAccurate() *Matcher[T] {
	m.scannerConfig.Mode = modePriorityAccurate

	return m
}

func (m *Matcher[T]) WithModePriorityFuzzy() *Matcher[T] {
	m.scannerConfig.Mode = modePriorityFuzzy

	return m
}

func (m *Matcher[T]) WithFuzzy(fc FuzzyConfig) *Matcher[T] {
	fc.enabled = true
	m.scannerConfig.Fuzzy = fc

	return m
}

func (m *Matcher[T]) WithDebug() *Matcher[T] {
	m.scannerConfig.Debug = true

	return m
}

func (m *Matcher[T]) WithScanner(keyName string, items []map[string]string) *Matcher[T] {
	m.newScannerFun = func(rule extract.Rule, sc ScannerConfig) (*scanner, error) {
		return newScanner(keyName, items, sc), nil
	}

	return m
}
