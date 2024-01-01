package match

import (
	"github.com/auho/go-etl/v2/job/means"
)

// all
// line
// flag

// NewExportLabel
//
// keys []string, name of tokenize
// df: map[string]any, defaultValues
func NewExportLabel(keys []string, df map[string]any, fn func(LabelResults, means.Ruler) []map[string]any) *ExportLabelResults {
	return NewExport[LabelResults](keys, df, fn)
}

func NewExportLabelAll(rule means.Ruler) *ExportLabelResults {
	var keys []string
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		keys = append(keys, _ta)
		values[_ta] = ""
	}

	keys = append(keys, rule.KeywordNameAlias(), rule.KeywordAmountNameAlias())
	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordAmountNameAlias()] = 0

	return NewExportLabel(keys, values, func(results LabelResults, rule means.Ruler) []map[string]any {
		return results.ToAll(rule)
	})
}

func NewExportLabelLine(rule means.Ruler) *ExportLabelResults {
	var keys []string
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		keys = append(keys, _ta)
		values[_ta] = ""
	}

	keys = append(keys, rule.KeywordNameAlias(), rule.LabelNumNameAlias(), rule.KeywordNumNameAlias(), rule.KeywordAmountNameAlias())
	values[rule.KeywordNameAlias()] = ""
	values[rule.LabelNumNameAlias()] = 0
	values[rule.KeywordNumNameAlias()] = 0
	values[rule.KeywordAmountNameAlias()] = 0

	return NewExportLabel(keys, values, func(results LabelResults, rule means.Ruler) []map[string]any {
		return results.ToLine(rule)
	})
}

func NewExportLabelFlag(rule means.Ruler) *ExportLabelResults {
	var keys []string
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		keys = append(keys, _ta)
		values[_ta] = ""
	}

	keys = append(keys, rule.KeywordNameAlias())
	values[rule.KeywordNameAlias()] = ""
	values[rule.NameAlias()] = 0

	return NewExportLabel(keys, values, func(results LabelResults, rule means.Ruler) []map[string]any {
		return results.ToFlag(rule)
	})
}
