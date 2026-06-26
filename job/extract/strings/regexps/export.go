package regexps

import (
	"github.com/auho/go-etl/v3/job/extract"
)

func NewExport(rule extract.Rule, df map[string]any, fn func(Results, extract.Rule) []map[string]any) *extract.Exporter[Results] {
	return extract.NewExporter(df, func(ctx extract.ExportContext[Results]) []map[string]any {
		return fn(ctx.Results, rule)
	}, extract.WithRule[Results](rule))
}

func NewExportDefault(rule extract.Rule, fn func(Results, extract.Rule) []map[string]any) *extract.Exporter[Results] {
	return NewExport(rule, map[string]any{rule.NameAlias(): ""}, fn)
}

func NewExportAll(rule extract.Rule) *extract.Exporter[Results] {
	df := map[string]any{
		rule.NameAlias():              "",
		rule.KeywordAmountNameAlias(): 0,
	}
	return NewExport(rule, df, func(results Results, rule extract.Rule) []map[string]any {
		if results == nil {
			return nil
		}
		return results.ToAll(rule)
	})
}

func NewExportLine(rule extract.Rule) *extract.Exporter[Results] {
	df := map[string]any{
		rule.NameAlias():              "",
		rule.KeywordNumNameAlias():    0,
		rule.KeywordAmountNameAlias(): 0,
	}
	return NewExport(rule, df, func(results Results, rule extract.Rule) []map[string]any {
		if results == nil {
			return nil
		}
		return results.ToLine(rule)
	})
}

func NewExportFlag(rule extract.Rule) *extract.Exporter[Results] {
	df := map[string]any{
		rule.NameAlias():        0,
		rule.KeywordNameAlias(): "",
	}
	return NewExport(rule, df, func(results Results, rule extract.Rule) []map[string]any {
		if results == nil {
			return nil
		}
		return results.ToFlag(rule)
	})
}
