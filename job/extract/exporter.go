package extract

import (
	"maps"

	"github.com/auho/go-etl/v3/tool/mapx"
)

// ExportContext is passed to the convert function.
type ExportContext[T any] struct {
	Rule    Rule // optional, used by match/tag/contains/regexps
	Results T
	Format  any // optional, sub-package custom type passed through
}

// Exporter is the unified export executor.
// It holds defaultValues/keys, runs convert, applies pluck, and builds Result.
// Implements FieldSpec.
type Exporter[T any] struct {
	rule          Rule
	keys          []string
	defaultValues map[string]any
	format        any
	convert       func(ExportContext[T]) []map[string]any
	pluckKeys     []string // nil = no pluck
}

type ExporterOption[T any] func(*Exporter[T])

func WithRule[T any](rule Rule) ExporterOption[T] {
	return func(e *Exporter[T]) { e.rule = rule }
}

func WithFormat[T any](format any) ExporterOption[T] {
	return func(e *Exporter[T]) { e.format = format }
}

func NewExporter[T any](
	df map[string]any,
	convert func(ExportContext[T]) []map[string]any,
	opts ...ExporterOption[T],
) *Exporter[T] {
	keys := make([]string, 0, len(df))
	for k := range df {
		keys = append(keys, k)
	}
	e := &Exporter[T]{
		keys:          keys,
		defaultValues: df,
		convert:       convert,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// FieldSpec implementation

func (e *Exporter[T]) Keys() []string             { return e.keys }
func (e *Exporter[T]) DefaultValues() map[string]any { return e.defaultValues }
func (e *Exporter[T]) GetRule() Rule              { return e.rule }

// Pluck selects a subset of fields.
func (e *Exporter[T]) Pluck(keys []string) *Exporter[T] {
	df := maps.Clone(e.defaultValues)

	e.keys = make([]string, 0)
	e.defaultValues = make(map[string]any)
	e.pluckKeys = make([]string, 0)

	for _, key := range keys {
		if v, ok := df[key]; ok {
			e.keys = append(e.keys, key)
			e.defaultValues[key] = v
			e.pluckKeys = append(e.pluckKeys, key)
		}
	}

	return e
}

// ToToken converts results to extract.Result.
// ok: whether results are valid (caller decides emptiness check).
func (e *Exporter[T]) ToToken(results T, ok bool) Result {
	token := Result{}
	if !ok {
		return token
	}

	token.SetOK()
	token.SetResultsFunc(func() []map[string]any {
		ret := e.convert(ExportContext[T]{
			Rule:    e.rule,
			Results: results,
			Format:  e.format,
		})

		if e.pluckKeys != nil {
			ret = mapx.PluckSliceMap(ret, e.pluckKeys)
		}

		return ret
	})

	return token
}
