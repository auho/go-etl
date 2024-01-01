package match

import (
	"github.com/auho/go-etl/v2/job/means"
)

func NewSearchLabel(rule means.Ruler, export *ExportLabelResults, srf SearchResultsFun[LabelResults]) *SearchLabelResults {
	return NewSearch[LabelResults](rule, export, srf)
}

func NewSearchWholeLabels(rule means.Ruler, export *ExportLabelResults) *SearchLabelResults {
	return NewSearchLabel(rule, export, func(ctx *SearchContextLabelResults, c []string) LabelResults {
		return ctx.Matcher.MatchLabel(c)
	})
}

func NewSearchLabels(rule means.Ruler, export *ExportLabelResults) *SearchLabelResults {
	return NewSearchLabel(rule, export, func(ctx *SearchContextLabelResults, c []string) LabelResults {
		return ctx.Matcher.MatchLabel(c)
	})
}
