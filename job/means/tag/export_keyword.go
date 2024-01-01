package tag

import (
	"github.com/auho/go-etl/v2/job/means"
)

// all
// line
// flag

// NewExportKeyword
//
// keys []string, name of tokenize
// df: map[string]any, defaultValues
func NewExportKeyword(keys []string, df map[string]any, fn func(Results, means.Ruler) []map[string]any) *ExportResults {
	e := &ExportResults{
		keys:           keys,
		defaultValues:  df,
		resultsToToken: fn,
	}

	return e
}

func NewExportKeywordAll(rule means.Ruler) *ExportResults {
	var keys []string
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		keys = append(keys, _ta)
		values[_ta] = ""
	}

	keys = append(keys, rule.KeywordNameAlias(), rule.KeywordAmountNameAlias())
	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordAmountNameAlias()] = 0

	return NewExportKeyword(keys, values, func(results Results, rule means.Ruler) []map[string]any {
		return results.ToAll(rule)
	})
}

func NewExportKeywordLine(rule means.Ruler) *ExportResults {
	var keys []string
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		keys = append(keys, _ta)
		values[_ta] = ""
	}

	keys = append(keys, rule.KeywordNameAlias(), rule.KeywordNumNameAlias())
	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordNumNameAlias()] = 0

	return NewExportKeyword(keys, values, func(results Results, rule means.Ruler) []map[string]any {
		return results.ToLine(rule)
	})
}
