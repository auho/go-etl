package match

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

// all
// line
// flag

var _ search.Exporter = (*ExportLabel)(nil)

type GenExportLabel func(LabelResults, means.Ruler) *ExportLabel

type ExportLabel struct {
	Export
	Results           LabelResults
	ResultsToTokenize func(results LabelResults) []map[string]any
}

func NewExportLabel(rs LabelResults, rule means.Ruler, keys []string, df map[string]any, fn func(LabelResults) []map[string]any) *ExportLabel {
	e := &ExportLabel{
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

func (e *ExportLabel) init() {
	e.Ok = true
	if e.Results == nil {
		e.Ok = false
	}
}

func (e *ExportLabel) ToTokenize() []map[string]any {
	return e.ResultsToTokenize(e.Results)
}

func NewExportLabelAll(rs LabelResults, rule means.Ruler) *ExportLabel {
	var keys []string
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		keys = append(keys, _ta)
		values[_ta] = ""
	}

	keys = append(keys, rule.KeywordNameAlias(), rule.KeywordAmountNameAlias())
	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordAmountNameAlias()] = 0

	return NewExportLabel(rs, rule, keys, values, func(results LabelResults) []map[string]any {
		return results.ToAll(rule)
	})
}

func NewExportLabelLine(rs LabelResults, rule means.Ruler) *ExportLabel {
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

	return NewExportLabel(rs, rule, keys, values, func(results LabelResults) []map[string]any {
		return results.ToLine(rule)
	})
}

func NewExportLabelFlag(rs LabelResults, rule means.Ruler) *ExportLabel {
	var keys []string
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		keys = append(keys, _ta)
		values[_ta] = ""
	}

	keys = append(keys, rule.KeywordNameAlias())
	values[rule.KeywordNameAlias()] = ""

	return NewExportLabel(rs, rule, keys, values, func(results LabelResults) []map[string]any {
		return results.ToFlag(rule)
	})
}
