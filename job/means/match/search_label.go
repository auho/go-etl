package match

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

func GenSearchLabel(rule means.Ruler, gek GenExportLabel, ste SearchToExport[LabelResults]) *Search[LabelResults] {
	_gek := func(results LabelResults, ruler means.Ruler) search.Exporter {
		return gek(results, ruler)
	}

	return NewSearch[LabelResults](rule, _gek, ste)
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
