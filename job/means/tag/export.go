package tag

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
	rule           means.Ruler
	resultsToToken func(ctx ExportContext[T]) []map[string]any
	format         Format
	keys           []string
	defaultValues  map[string]any
}

func NewExport[T ResultsEntity](rule means.Ruler, df map[string]any, fn func(ctx ExportContext[T]) []map[string]any) *Export[T] {
	var keys []string
	for k := range df {
		keys = append(keys, k)
	}

	return &Export[T]{
		rule:           rule,
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

func (e *Export[T]) GetRule() means.Ruler {
	return e.rule
}

func (e *Export[T]) Pluck(keys []string) *Export[T] {
	df := maps.Clone(e.defaultValues)

	e.keys = make([]string, 0)
	e.defaultValues = make(map[string]any)

	for _, key := range keys {
		if v, ok := df[key]; ok {
			e.keys = append(e.keys, key)
			e.defaultValues[key] = v
		}
	}

	df = nil

	return e
}

func (e *Export[T]) WithFormat(format Format) *Export[T] {
	e.format = format

	return e
}

func (e *Export[T]) ToToken(results T) search.Token {
	token := search.Token{}

	if len(results) > 0 {
		token.SetOk()
		token.SetTokenizerFunc(func() []map[string]any {
			ret := e.resultsToToken(ExportContext[T]{
				Rule:    e.rule,
				Results: results,
				Format:  e.format,
			})

			// for pluck
			return maps2.PluckSliceMap(ret, e.keys)
		})
	}

	return token
}
