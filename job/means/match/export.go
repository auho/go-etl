package match

import (
	"maps"

	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
	maps2 "github.com/auho/go-etl/v2/tool/maps"
)

var _ search.Exporter = (*Export[Results])(nil)
var _ search.Exporter = (*Export[LabelResults])(nil)
var _ search.Exporter = (*ExportResults)(nil)
var _ search.Exporter = (*ExportLabelResults)(nil)

type ExportContextResults = ExportContext[Results]
type ExportContextLabelResults = ExportContext[LabelResults]

type ExportContext[T ResultsEntity] struct {
	Rule    means.Ruler
	Results T
	Format  Format
}

type ExportResults = Export[Results]
type ExportLabelResults = Export[LabelResults]

type Export[T ResultsEntity] struct {
	resultsToToken func(ctx ExportContext[T]) []map[string]any
	format         Format
	keys           []string
	defaultValues  map[string]any
}

func NewExport[T ResultsEntity](keys []string, df map[string]any, fn func(ctx ExportContext[T]) []map[string]any) *Export[T] {
	return &Export[T]{
		keys:           keys,
		defaultValues:  df,
		resultsToToken: fn,
		format:         DefaultFormat,
	}
}

func (e *Export[T]) GetKeys() []string {
	return e.keys
}

func (e *Export[T]) GetDefaultValues() map[string]any {
	return e.defaultValues
}

func (e *Export[T]) Pluck(keys []string) *Export[T] {
	df := maps.Clone(e.defaultValues)

	e.keys = keys
	e.defaultValues = make(map[string]any, len(e.keys))

	for _, key := range e.keys {
		e.defaultValues[key] = df[key]
	}

	df = nil

	return e
}

func (e *Export[T]) WithFormat(format Format) *Export[T] {
	e.format = format

	return e
}

func (e *Export[T]) ToToken(results T, rule means.Ruler) search.Token {
	token := search.Token{}

	if len(results) > 0 {
		token.SetOk()
		token.SetTokenizerFunc(func() []map[string]any {
			ret := e.resultsToToken(ExportContext[T]{
				Rule:    rule,
				Results: results,
				Format:  e.format,
			})

			// for pluck
			return maps2.PluckSliceMap(ret, e.keys)
		})
	}

	return token
}
