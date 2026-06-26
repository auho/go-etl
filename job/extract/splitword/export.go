package splitword

import (
	"github.com/auho/go-etl/v3/job/extract"
)

type ExportContext struct {
	Results Results
	Format  Format
}

var _ extract.FieldSpec = (*Export)(nil)

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

func (e *Export) Keys() []string {
	var keys []string
	for k := range e.defaultValues {
		keys = append(keys, k)
	}

	return keys
}

func (e *Export) DefaultValues() map[string]any {
	return e.defaultValues
}

func (e *Export) WithFormat(format Format) *Export {
	e.format = format

	return e
}

func (e *Export) ToToken(results Results) extract.Result {
	token := extract.Result{}

	if len(results) > 0 {
		token.SetOK()
		token.SetResultsFunc(func() []map[string]any {
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
