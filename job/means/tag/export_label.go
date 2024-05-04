package tag

import (
	"github.com/auho/go-etl/v2/job/means"
)

// all
// line
// flag

// NewExportLabel
//
// df: map[string]any, defaultValues
func NewExportLabel(rule means.Ruler, df map[string]any, fn func(ExportContextLabelResults) []map[string]any) *ExportLabelResults {
	return NewExport[LabelResults](rule, df, fn)
}

func NewExportLabelAll(rule means.Ruler) *ExportLabelResults {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}

	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordAmountNameAlias()] = 0

	return NewExportLabel(rule, values, func(ctx ExportContextLabelResults) []map[string]any {
		return ctx.Results.ToAll(rule, ctx.Format)
	})
}

func NewExportLabelLine(rule means.Ruler) *ExportLabelResults {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}

	values[rule.KeywordNameAlias()] = ""
	values[rule.LabelNumNameAlias()] = 0
	values[rule.KeywordNumNameAlias()] = 0
	values[rule.KeywordAmountNameAlias()] = 0

	return NewExportLabel(rule, values, func(ctx ExportContextLabelResults) []map[string]any {
		return ctx.Results.ToLine(rule, ctx.Format)
	})
}

func NewExportLabelFlag(rule means.Ruler) *ExportLabelResults {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}

	values[rule.KeywordNameAlias()] = ""
	values[rule.NameAlias()] = 0

	return NewExportLabel(rule, values, func(ctx ExportContextLabelResults) []map[string]any {
		return ctx.Results.ToFlag(rule, ctx.Format)
	})
}
