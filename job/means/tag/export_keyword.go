package tag

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Exporter = (*ExportKeywordAll)(nil)
var _ search.Exporter = (*ExportKeywordLine)(nil)
var _ search.Exporter = (*ExportKeywordFlag)(nil)

type NewExportKeyword func(rs Results, rule means.Ruler) search.Exporter

type ExportKeyword struct {
	export
	Results Results
}

func newExportKeyword(rs Results, rule means.Ruler) ExportKeyword {
	e := ExportKeyword{
		export: export{
			Rule: rule,
		},
		Results: rs,
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

func (e *ExportKeyword) defaultValuesTagsWithKeywordAndNum() ([]string, map[string]any) {
	keys, m := e.defaultValuesTagsWithKeyword()
	m[e.Rule.KeywordNumNameAlias()] = 0

	return append(keys, e.Rule.KeywordNumNameAlias()), m
}

func (e *ExportKeyword) defaultValuesTagsWithKeyword() ([]string, map[string]any) {
	keys, m := e.defaultValuesTags()
	m[e.Rule.KeywordNameAlias()] = ""

	return append(keys, e.Rule.KeywordNameAlias()), m
}

func (e *ExportKeyword) defaultValuesTags() ([]string, map[string]any) {
	m := map[string]any{}

	m[e.Rule.NameAlias()] = ""
	for _, _la := range e.Rule.LabelsAlias() {
		m[_la] = ""
	}

	return append([]string{e.Rule.NameAlias()}, e.Rule.LabelsAlias()...), m
}

type ExportKeywordAll struct {
	ExportKeyword
}

func NewExportKeywordAll(rs Results, rule means.Ruler) search.Exporter {
	e := &ExportKeywordAll{
		ExportKeyword: newExportKeyword(rs, rule),
	}

	e.initial()

	return e
}

func (e *ExportKeywordAll) initial() {
	e.keys, e.defaultValues = e.defaultValuesTagsWithKeywordAndNum()
}

func (e *ExportKeywordAll) ToTokenize() []map[string]any {
	return e.Results.ToTags(e.Rule)
}

type ExportKeywordLine struct {
	ExportKeyword
}

func NewExportKeywordLine(rs Results, rule means.Ruler) search.Exporter {
	e := &ExportKeywordLine{
		ExportKeyword: newExportKeyword(rs, rule),
	}

	e.initial()

	return e
}

func (e *ExportKeywordLine) initial() {
	e.keys, e.defaultValues = e.defaultValuesTagsWithKeyword()
}

func (e *ExportKeywordLine) ToTokenize() []map[string]any {
	return e.Results.ToLine(e.Rule)
}

type ExportKeywordFlag struct {
	ExportKeyword
}

func NewExportKeywordFlag(rs Results, rule means.Ruler) search.Exporter {
	e := &ExportKeywordFlag{
		ExportKeyword: newExportKeyword(rs, rule),
	}

	e.initial()

	return e
}

func (e *ExportKeywordFlag) initial() {
	e.keys, e.defaultValues = e.defaultValuesTagsWithKeyword()
}

func (e *ExportKeywordFlag) ToTokenize() []map[string]any {
	return e.Results.ToFlag(e.Rule)
}
