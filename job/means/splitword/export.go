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

	keys          []string
	defaultValues map[string]any
}

func NewExport(keys []string, df map[string]any, fn func(ExportContext) []map[string]any) *Export {
	return &Export{
		keys:           keys,
		defaultValues:  df,
		resultsToToken: fn,
		format:         DefaultFormat,
	}
}

func (e *Export) GetKeys() []string {
	return e.keys
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
	return NewExport([]string{WordName}, map[string]any{WordName: ""}, func(ctx ExportContext) []map[string]any {
		return ctx.Results.ToAll()
	})
}

func NewExportLine() *Export {
	return NewExport([]string{WordName}, map[string]any{WordName: ""}, func(ctx ExportContext) []map[string]any {
		return ctx.Results.ToLine(ctx.Format)
	})
}
