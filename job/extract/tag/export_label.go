package tag

import "github.com/auho/go-etl/v3/job/extract"

var _ extract.FieldSpec = (*extract.Exporter[LabelResults])(nil)

type ExportLabelResults = extract.Exporter[LabelResults]

// all
// line
// flag

// NewExportLabel
//
// df: map[string]any, defaultValues
func NewExportLabel(rule extract.Rule, df map[string]any, fn func(extract.ExportContext[LabelResults]) []map[string]any) *extract.Exporter[LabelResults] {
	return extract.NewExporter(df, fn, extract.WithRule[LabelResults](rule), extract.WithFormat[LabelResults](DefaultFormat))
}

func NewExportLabelAll(rule extract.Rule) *extract.Exporter[LabelResults] {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}

	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordAmountNameAlias()] = 0

	return extract.NewExporter(values, func(ctx extract.ExportContext[LabelResults]) []map[string]any {
		return ctx.Results.ToAll(rule, ctx.Format.(Format))
	}, extract.WithRule[LabelResults](rule), extract.WithFormat[LabelResults](DefaultFormat))
}

func NewExportLabelLine(rule extract.Rule) *extract.Exporter[LabelResults] {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}

	values[rule.KeywordNameAlias()] = ""
	values[rule.LabelNumNameAlias()] = 0
	values[rule.KeywordNumNameAlias()] = 0
	values[rule.KeywordAmountNameAlias()] = 0

	return extract.NewExporter(values, func(ctx extract.ExportContext[LabelResults]) []map[string]any {
		return ctx.Results.ToLine(rule, ctx.Format.(Format))
	}, extract.WithRule[LabelResults](rule), extract.WithFormat[LabelResults](DefaultFormat))
}

func NewExportLabelFlag(rule extract.Rule) *extract.Exporter[LabelResults] {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}

	values[rule.KeywordNameAlias()] = ""
	values[rule.NameAlias()] = 0

	return extract.NewExporter(values, func(ctx extract.ExportContext[LabelResults]) []map[string]any {
		return ctx.Results.ToFlag(rule, ctx.Format.(Format))
	}, extract.WithRule[LabelResults](rule), extract.WithFormat[LabelResults](DefaultFormat))
}
