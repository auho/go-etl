package tag

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Exporter = (*ExportKeywordTag)(nil)
var _ search.Exporter = (*ExportKeywordLine)(nil)
var _ search.Exporter = (*ExportKeywordFlag)(nil)

type NewExportKeyword func(rs Results, rule means.Ruler) search.Exporter

type ExportKeyword struct {
	Ok      bool
	Results Results
	Rule    means.Ruler

	defaultValues map[string]any
}

func (e *ExportKeyword) IsOk() bool {
	return e.Ok
}

func (e *ExportKeyword) DefaultValues() map[string]any {
	return e.defaultValues
}

func (e *ExportKeyword) init() {
	e.Ok = true
	if e.Results == nil {
		e.Ok = false
	}
}

func (e *ExportKeyword) defaultValuesTagsWithKeywordAndNum() map[string]any {
	m := e.defaultValuesTagsWithKeyword()
	m[e.Rule.KeywordNumNameAlias()] = 0

	return m
}

func (e *ExportKeyword) defaultValuesTagsWithKeyword() map[string]any {
	m := e.defaultValuesTags()
	m[e.Rule.KeywordNameAlias()] = ""

	return m
}

func (e *ExportKeyword) defaultValuesTags() map[string]any {
	m := map[string]any{}

	m[e.Rule.NameAlias()] = ""
	for _, _la := range e.Rule.LabelsAlias() {
		m[_la] = ""
	}

	return m
}

type ExportKeywordTag struct {
	ExportKeyword
}

func NewExportKeywordTag(rs Results, rule means.Ruler) search.Exporter {
	e := &ExportKeywordTag{ExportKeyword{
		Results: rs,
		Rule:    rule,
	}}

	e.init()
	e.initial()

	return e
}

func (e *ExportKeywordTag) initial() {
	e.defaultValues = e.defaultValuesTagsWithKeywordAndNum()
}

func (e *ExportKeywordTag) ToTokenize() []map[string]any {
	return e.Results.ToTags(e.Rule)
}

type ExportKeywordLine struct {
	ExportKeyword
}

func NewExportKeywordLine(rs Results, rule means.Ruler) search.Exporter {
	e := &ExportKeywordLine{ExportKeyword{
		Results: rs,
		Rule:    rule,
	}}

	e.init()
	e.initial()

	return e
}

func (e *ExportKeywordLine) initial() {
	e.defaultValues = e.defaultValuesTagsWithKeyword()
}

func (e *ExportKeywordLine) ToTokenize() []map[string]any {
	return e.Results.ToLine(e.Rule)
}

type ExportKeywordFlag struct {
	ExportKeyword
}

func NewExportKeywordFlag(rs Results, rule means.Ruler) search.Exporter {
	e := &ExportKeywordFlag{ExportKeyword{
		Results: rs,
		Rule:    rule,
	}}

	e.init()
	e.initial()

	return e
}

func (e *ExportKeywordFlag) initial() {
	e.defaultValues = e.defaultValuesTagsWithKeyword()
}

func (e *ExportKeywordFlag) ToTokenize() []map[string]any {
	return e.Results.ToFlag(e.Rule)
}
