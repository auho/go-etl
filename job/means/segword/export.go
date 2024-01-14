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

	keys          []string
	defaultValues map[string]any
}

func NewExport(keys []string, df map[string]any, fn func(ExportContext) []map[string]any) *Export {
	return &Export{
		keys:           keys,
		defaultValues:  df,
		resultsToToken: fn,
		filterFunc:     DefaultFilterFunc,
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
		[]string{},
		map[string]any{format.TokenName: "", format.FlagName: ""},
		func(ctx ExportContext) []map[string]any {
			return ctx.Results.ToAll(ctx.Format)
		},
	).WithFormat(format)
}

func NewExportLine() *Export {
	format := DefaultFormat

	return NewExport(
		[]string{format.TokenName},
		map[string]any{format.TokenName: ""},
		func(ctx ExportContext) []map[string]any {
			return ctx.Results.ToLine(ctx.Format)
		},
	).WithFormat(format)
}
