package tag

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

//type NewExporterLabelResult func(rs LabelResults, rule means.Ruler) *ExportLabelResult // TODO
// TODO WIP this file

var _ search.Exporter = (*ExportLabelTag)(nil)
var _ search.Exporter = (*ExportLabelLine)(nil)
var _ search.Exporter = (*ExportLabelFlag)(nil)

type NewExportLabel func(rs Results, rule means.Ruler) *ExportLabel

type ExportLabel struct {
	Ok      bool
	Results Results
	Rule    means.Ruler

	defaultValues map[string]any
}

func (e *ExportLabel) IsOk() bool {
	return e.Ok
}

func (e *ExportLabel) DefaultValues() map[string]any {
	return e.defaultValues
}

func (e *ExportLabel) init() {
	e.Ok = true
	if e.Results == nil {
		e.Ok = false
	}
}

func (e *ExportLabel) defaultValuesTagsWithKeywordAndNum() map[string]any {
	m := e.defaultValuesTagsWithKeyword()
	m[e.Rule.KeywordNumNameAlias()] = 0

	return m
}

func (e *ExportLabel) defaultValuesTagsWithKeyword() map[string]any {
	m := e.defaultValuesTags()
	m[e.Rule.KeywordNameAlias()] = ""

	return m
}

func (e *ExportLabel) defaultValuesTags() map[string]any {
	m := map[string]any{}

	m[e.Rule.NameAlias()] = ""
	for _, _la := range e.Rule.LabelsAlias() {
		m[_la] = ""
	}

	return m
}

type ExportLabelTag struct {
	ExportLabel
}

func NewExportLabelTag(rs Results, rule means.Ruler) *ExportLabelTag {
	e := &ExportLabelTag{ExportLabel{
		Results: rs,
		Rule:    rule,
	}}

	e.init()
	e.initial()

	return e
}

func (e *ExportLabelTag) initial() {
	e.defaultValues = e.defaultValuesTagsWithKeywordAndNum()
}

func (e *ExportLabelTag) ToTokenize() []map[string]any {
	return e.Results.ToTags(e.Rule)
}

type ExportLabelLine struct {
	ExportLabel
}

func NewExportLabelLine(rs Results, rule means.Ruler) *ExportLabelLine {
	e := &ExportLabelLine{ExportLabel{
		Results: rs,
		Rule:    rule,
	}}

	e.init()
	e.initial()

	return e
}

func (e *ExportLabelLine) initial() {
	e.defaultValues = e.defaultValuesTagsWithKeyword()
}

func (e *ExportLabelLine) ToTokenize() []map[string]any {
	return e.Results.ToLine(e.Rule)
}

type ExportLabelFlag struct {
	ExportLabel
}

func NewExportLabelFlag(rs Results, rule means.Ruler) *ExportLabelFlag {
	e := &ExportLabelFlag{ExportLabel{
		Results: rs,
		Rule:    rule,
	}}

	e.init()
	e.initial()

	return e
}

func (e *ExportLabelFlag) initial() {
	e.defaultValues = e.defaultValuesTagsWithKeyword()
}

func (e *ExportLabelFlag) ToTokenize() []map[string]any {
	return e.Results.ToFlag(e.Rule)
}
