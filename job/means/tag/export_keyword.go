package tag

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Exporter = (*ExportKeyword)(nil)

type GenExportKeyword func(Results, means.Ruler) *ExportKeyword

type ExportKeyword struct {
	Export
	Results           Results
	ResultsToTokenize func(results Results) []map[string]any
}

// NewExportKeyword
//
// keys []string, name of tokenize
// df: map[string]any, defaultValues
// fn: fn func(Results) []map[string]any, Results to tokenize
func NewExportKeyword(rs Results, rule means.Ruler, keys []string, df map[string]any, fn func(Results) []map[string]any) *ExportKeyword {
	e := &ExportKeyword{
		Export: Export{
			Rule:          rule,
			Keys:          keys,
			DefaultValues: df,
		},
		Results:           rs,
		ResultsToTokenize: fn,
	}

	e.init()

	return e
}

func (e *ExportKeyword) init() {
	e.Ok = true
	if e.Results == nil {
		e.Ok = false
	}
}

func (e *ExportKeyword) ToTokenize() []map[string]any {
	return e.ResultsToTokenize(e.Results)
}

func NewExportKeywordAll(rs Results, rule means.Ruler) *ExportKeyword {
	var keys []string
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		keys = append(keys, _ta)
		values[_ta] = ""
	}

	keys = append(keys, rule.KeywordNameAlias(), rule.KeywordAmountNameAlias())
	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordAmountNameAlias()] = 0

	return NewExportKeyword(rs, rule, keys, values, func(results Results) []map[string]any {
		return results.ToAll(rule)
	})
}

func NewExportKeywordLine(rs Results, rule means.Ruler) *ExportKeyword {
	var keys []string
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		keys = append(keys, _ta)
		values[_ta] = ""
	}

	keys = append(keys, rule.KeywordNameAlias(), rule.KeywordNumNameAlias())
	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordNumNameAlias()] = 0

	return NewExportKeyword(rs, rule, keys, values, func(results Results) []map[string]any {
		return results.ToLine(rule)
	})
}
