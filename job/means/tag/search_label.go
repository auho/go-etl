package tag

import (
	"github.com/auho/go-etl/v2/job/means"
)

func GenSearchLabel(rule means.Ruler, gek *Export[LabelResults], ste SearchResultsFun[LabelResults]) *SearchLabelResults {
	return NewSearch[LabelResults](rule, gek, ste)
}

func NewSearchWholeLabels(rule means.Ruler) *SearchLabelResults {
	return GenSearchLabel(rule, NewExportLabelLine(rule), func(ctx *SearchContext[LabelResults], c []string) LabelResults {
		return ctx.Matcher.MatchLabel(c)
	})
}

func NewSearchLabels(rule means.Ruler, gel *Export[LabelResults]) *SearchLabelResults {
	return GenSearchLabel(rule, gel, func(ctx *SearchContext[LabelResults], c []string) LabelResults {
		return ctx.Matcher.MatchLabel(c)
	})
}
