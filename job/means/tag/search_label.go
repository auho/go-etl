package tag

import (
	"github.com/auho/go-etl/v2/job/means"
)

func GenSearchLabel(rule means.Ruler, gek GenExportLabel, ste SearchResults[LabelResults]) *Search[LabelResults] {
	return NewSearch[LabelResults](rule, gek(rule), ste)
}

func NewSearchWholeLabels(rule means.Ruler) *Search[LabelResults] {
	return GenSearchLabel(rule, NewExportLabelLine, func(ctx *SearchContext[LabelResults], c []string) LabelResults {
		return ctx.Matcher.MatchLabel(c)
	})
}

func NewSearchLabels(rule means.Ruler, gel GenExportLabel) *Search[LabelResults] {
	return GenSearchLabel(rule, gel, func(ctx *SearchContext[LabelResults], c []string) LabelResults {
		return ctx.Matcher.MatchLabel(c)
	})
}
