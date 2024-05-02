package tag

import (
	"github.com/auho/go-etl/v2/job/means"
)

// all
// line
// flag

// NewExportKeyword
//
// df: map[string]any, defaultValues
func NewExportKeyword(df map[string]any, fn func(ExportContextResults) []map[string]any) *ExportResults {
	return NewExport[Results](df, fn)
}

func NewExportKeywordAll(rule means.Ruler) *ExportResults {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}

	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordAmountNameAlias()] = 0

	return NewExportKeyword(values, func(ctx ExportContextResults) []map[string]any {
		return ctx.Results.ToAll(rule)
	})
}

func NewExportKeywordLine(rule means.Ruler) *ExportResults {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}

	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordNumNameAlias()] = 0

	return NewExportKeyword(values, func(ctx ExportContextResults) []map[string]any {
		return ctx.Results.ToLine(rule, ctx.Format)
	})
}

func NewExportKeywordFlag(rule means.Ruler) *ExportResults {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}

	values[rule.KeywordNameAlias()] = ""
	values[rule.NameAlias()] = 0

	return NewExportKeyword(values, func(ctx ExportContextResults) []map[string]any {
		return ctx.Results.ToFlag(rule, ctx.Format)
	})
}
