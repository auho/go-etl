package splitword

import (
	"github.com/auho/go-etl/v2/job/explore/search"
)

type ExportContext struct {
	Results Results
	Format  Format
}

var _ search.Exporter = (*Export)(nil)

type Export struct {
	format         Format
	resultsToToken func(ExportContext) []map[string]any

	defaultValues map[string]any
}

func NewExport(df map[string]any, fn func(ExportContext) []map[string]any) *Export {
	return &Export{
		defaultValues:  df,
		resultsToToken: fn,
		format:         DefaultFormat,
	}
}

func (e *Export) GetKeys() []string {
	var keys []string
	for k := range e.defaultValues {
		keys = append(keys, k)
	}

	return keys
}

func (e *Export) GetDefaultValues() map[string]any {
	return e.defaultValues
}

func (e *Export) WithFormat(format Format) *Export {
	e.format = format

	return e
}

func (e *Export) ToToken(results Results) search.Token {
	token := search.Token{}

	if len(results) > 0 {
		token.SetOk()
		token.SetTokenizerFunc(func() []map[string]any {
			return e.resultsToToken(ExportContext{
				Results: results,
				Format:  e.format,
			})
		})
	}

	return token
}

func NewExportAll() *Export {
	format := DefaultFormat

	return NewExport(map[string]any{format.WordName: ""}, func(ctx ExportContext) []map[string]any {
		return ctx.Results.ToAll(ctx.Format)
	}).WithFormat(format)
}

func NewExportLine() *Export {
	format := DefaultFormat

	return NewExport(map[string]any{format.WordName: ""}, func(ctx ExportContext) []map[string]any {
		return ctx.Results.ToLine(ctx.Format)
	}).WithFormat(format)
}
