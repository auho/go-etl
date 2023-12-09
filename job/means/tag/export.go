package tag

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Exporter = (*ExportResultTag)(nil)
var _ search.Exporter = (*ExportResultLine)(nil)
var _ search.Exporter = (*ExportResultFlag)(nil)

type NewExporterResult func(rs Results, rule means.Ruler) search.Exporter
type NewExporterLabelResult func(rs LabelResults, rule means.Ruler) search.Exporter

type ExportResult struct {
	Ok      bool
	Results Results
	Rule    means.Ruler

	defaultValues map[string]any
}

func (e *ExportResult) IsOk() bool {
	return e.Ok
}

func (e *ExportResult) DefaultValues() map[string]any {
	return e.defaultValues
}

func (e *ExportResult) init() {
	e.Ok = true
	if e.Results == nil {
		e.Ok = false
	}
}

func (e *ExportResult) defaultValuesTagsWithKeywordAndNum() map[string]any {
	m := e.defaultValuesTagsWithKeyword()
	m[e.Rule.KeywordNumNameAlias()] = 0

	return m
}

func (e *ExportResult) defaultValuesTagsWithKeyword() map[string]any {
	m := e.defaultValuesTags()
	m[e.Rule.KeywordNameAlias()] = ""

	return m
}

func (e *ExportResult) defaultValuesTags() map[string]any {
	m := map[string]any{}

	m[e.Rule.NameAlias()] = ""
	for _, _la := range e.Rule.LabelsAlias() {
		m[_la] = ""
	}

	return m
}

type ExportResultTag struct {
	ExportResult
}

func NewExportResultTag(rs Results, rule means.Ruler) *ExportResultTag {
	e := &ExportResultTag{ExportResult{
		Results: rs,
		Rule:    rule,
	}}

	e.init()
	e.initial()

	return e
}

func (e *ExportResultTag) initial() {
	e.defaultValues = e.defaultValuesTagsWithKeywordAndNum()
}

func (e *ExportResultTag) ToTokenize() []map[string]any {
	return e.Results.ToTags(e.Rule)
}

type ExportResultLine struct {
	ExportResult
}

func NewExportResultLine(rs Results, rule means.Ruler) *ExportResultLine {
	e := &ExportResultLine{ExportResult{
		Results: rs,
		Rule:    rule,
	}}

	e.init()
	e.initial()

	return e
}

func (e *ExportResultLine) initial() {
	e.defaultValues = e.defaultValuesTagsWithKeyword()
}

func (e *ExportResultLine) ToTokenize() []map[string]any {
	return e.Results.ToLine(e.Rule)
}

type ExportResultFlag struct {
	ExportResult
}

func NewExportResultFlag(rs Results, rule means.Ruler) *ExportResultFlag {
	e := &ExportResultFlag{ExportResult{
		Results: rs,
		Rule:    rule,
	}}

	e.init()
	e.initial()

	return e
}

func (e *ExportResultFlag) initial() {
	e.defaultValues = e.defaultValuesTagsWithKeyword()
}

func (e *ExportResultFlag) ToTokenize() []map[string]any {
	return e.Results.ToFlag(e.Rule)
}
