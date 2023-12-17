package tag

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

func GenSearchLabel(rule means.Ruler, gek GenExportLabel, matchFn func(*Matcher, []string) LabelResults) *Search[LabelResults] {
	_gek := func(results LabelResults, ruler means.Ruler) search.Exporter {
		return gek(results, ruler)
	}

	return NewSearch[LabelResults](rule, _gek, func(mc SearchContext[LabelResults], contents []string) LabelResults {
		rets := matchFn(mc.Matcher, contents)
		if rets == nil {
			return nil
		}

		return rets
	})
}

func NewSearchWholeLabels(rule means.Ruler) *Search[LabelResults] {
	return GenSearchLabel(rule, NewExportLabelLine, func(m *Matcher, c []string) LabelResults {
		return m.MatchLabel(c)
	})
}

func NewSearchLabels(rule means.Ruler, gel GenExportLabel) *Search[LabelResults] {
	return GenSearchLabel(rule, gel, func(m *Matcher, c []string) LabelResults {
		return m.MatchLabel(c)
	})
}
