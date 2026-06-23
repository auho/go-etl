package segword

import (
	"unicode/utf8"

	"github.com/auho/go-etl/v2/job/explore/search"
)

var _ search.Exporter = (*Export)(nil)

var DefaultFilterFunc = func(result Result) bool {
	return utf8.RuneCountInString(result.Token) < 2 || result.Flag == "eng" || result.Flag == "m"
}

type ExportContext struct {
	Results Results
	Format  Format
}

type Export struct {
	format         Format
	resultsToToken func(ExportContext) []map[string]any
	filterFunc     func(Result) bool

	defaultValues map[string]any
}

func NewExport(df map[string]any, fn func(ExportContext) []map[string]any) *Export {
	return &Export{
		defaultValues:  df,
		resultsToToken: fn,
		filterFunc:     DefaultFilterFunc,
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
func (e *Export) WithFilterFunc(fn func(Result) bool) *Export {
	e.filterFunc = fn

	return e
}

func (e *Export) ToToken(results Results) search.Token {
	token := search.Token{}

	var newResults []Result
	for _, result := range results {
		if !e.filterFunc(result) {
			newResults = append(newResults, result)
		}
	}

	if len(newResults) > 0 {
		token.SetOk()
		token.SetTokenizerFunc(func() []map[string]any {
			return e.resultsToToken(ExportContext{
				Results: newResults,
				Format:  e.format,
			})
		})
	}

	return token
}

func NewExportAll() *Export {
	format := DefaultFormat

	return NewExport(
		map[string]any{format.TokenName: "", format.FlagName: ""},
		func(ctx ExportContext) []map[string]any {
			return ctx.Results.ToAll(ctx.Format)
		},
	).WithFormat(format)
}

func NewExportLine() *Export {
	format := DefaultFormat

	return NewExport(
		map[string]any{format.TokenName: ""},
		func(ctx ExportContext) []map[string]any {
			return ctx.Results.ToLine(ctx.Format)
		},
	).WithFormat(format)
}
