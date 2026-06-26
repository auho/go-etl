package match

import "github.com/auho/go-etl/v3/job/extract"

var _ extract.FieldSpec = (*extract.Exporter[Results])(nil)

type ExportResults = extract.Exporter[Results]

// all
// line
// flag

// NewExportKeyword
//
// df: map[string]any, defaultValues
func NewExportKeyword(rule extract.Rule, df map[string]any, fn func(extract.ExportContext[Results]) []map[string]any) *extract.Exporter[Results] {
	return extract.NewExporter(df, fn, extract.WithRule[Results](rule), extract.WithFormat[Results](DefaultFormat))
}

func NewExportKeywordAll(rule extract.Rule) *extract.Exporter[Results] {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}

	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordAmountNameAlias()] = 0

	return extract.NewExporter(values, func(ctx extract.ExportContext[Results]) []map[string]any {
		return ctx.Results.ToAll(rule)
	}, extract.WithRule[Results](rule))
}

func NewExportKeywordLine(rule extract.Rule) *extract.Exporter[Results] {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}

	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordNumNameAlias()] = 0

	return extract.NewExporter(values, func(ctx extract.ExportContext[Results]) []map[string]any {
		return ctx.Results.ToLine(rule, ctx.Format.(Format))
	}, extract.WithRule[Results](rule), extract.WithFormat[Results](DefaultFormat))
}

func NewExportKeywordFlag(rule extract.Rule) *extract.Exporter[Results] {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}

	values[rule.KeywordNameAlias()] = ""
	values[rule.NameAlias()] = 0

	return extract.NewExporter(values, func(ctx extract.ExportContext[Results]) []map[string]any {
		return ctx.Results.ToFlag(rule, ctx.Format.(Format))
	}, extract.WithRule[Results](rule), extract.WithFormat[Results](DefaultFormat))
}
