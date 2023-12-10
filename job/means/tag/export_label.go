package tag

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

// all
// line
// flag

var _ search.Exporter = (*ExportLabelAll)(nil)
var _ search.Exporter = (*ExportLabelLine)(nil)
var _ search.Exporter = (*ExportLabelFlag)(nil)

type NewExportLabel func(rs LabelResults, rule means.Ruler) *ExportLabel

type ExportLabel struct {
	export
	Results LabelResults
}

func newExportLabel(rs LabelResults, rule means.Ruler) ExportLabel {
	e := ExportLabel{
		export: export{
			Rule: rule,
		},
		Results: rs,
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

func (e *ExportLabel) defaultValuesTagsWithKeywordAndLabelNum() ([]string, map[string]any) {
	keys, m := e.defaultValuesTagsWithKeyword()
	m[e.Rule.LabelNumNameAlias()] = 0

	return append(keys, e.Rule.LabelNumNameAlias()), m
}

func (e *ExportLabel) defaultValuesTagsWithKeyword() ([]string, map[string]any) {
	keys, m := e.defaultValuesTags()
	m[e.Rule.KeywordNameAlias()] = ""

	return append(keys, e.Rule.KeywordNameAlias()), m
}

func (e *ExportLabel) defaultValuesTags() ([]string, map[string]any) {
	m := map[string]any{}

	m[e.Rule.NameAlias()] = ""
	for _, _la := range e.Rule.LabelsAlias() {
		m[_la] = ""
	}

	return append([]string{e.Rule.NameAlias()}, e.Rule.LabelsAlias()...), m
}

type ExportLabelAll struct {
	ExportLabel
}

func NewExportLabelAll(rs LabelResults, rule means.Ruler) *ExportLabelAll {
	e := &ExportLabelAll{
		ExportLabel: newExportLabel(rs, rule),
	}

	e.keys, e.defaultValues = e.defaultValuesTagsWithKeywordAndLabelNum()

	return e
}

func (e *ExportLabelAll) ToTokenize() []map[string]any {
	return e.Results.ToTags(e.Rule)
}

type ExportLabelLine struct {
	ExportLabel
}

func NewExportLabelLine(rs LabelResults, rule means.Ruler) *ExportLabelLine {
	e := &ExportLabelLine{
		ExportLabel: newExportLabel(rs, rule),
	}

	e.keys, e.defaultValues = e.defaultValuesTagsWithKeyword()

	return e
}

func (e *ExportLabelLine) ToTokenize() []map[string]any {
	return e.Results.ToLine(e.Rule)
}

type ExportLabelFlag struct {
	ExportLabel
}

func NewExportLabelFlag(rs LabelResults, rule means.Ruler) *ExportLabelFlag {
	e := &ExportLabelFlag{
		ExportLabel: newExportLabel(rs, rule),
	}

	e.keys, e.defaultValues = e.defaultValuesTagsWithKeyword()

	return e
}

func (e *ExportLabelFlag) ToTokenize() []map[string]any {
	return e.Results.ToFlag(e.Rule)
}
