package tag

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

func GenSearchKeyword(rule means.Ruler, gek GenExportKeyword, matchFn func(*Matcher, []string) Results) *Search[Results] {
	_gek := func(results Results, ruler means.Ruler) search.Exporter {
		return gek(results, ruler)
	}

	return NewSearch[Results](rule, _gek, func(mc SearchContext[Results], contents []string) Results {
		rets := matchFn(mc.Matcher, contents)
		if rets == nil {
			return nil
		}

		return rets
	})
}

func NewSearchFirstText(rule means.Ruler, gek GenExportKeyword) *Search[Results] {
	return GenSearchKeyword(rule, gek, func(m *Matcher, c []string) Results {
		return m.MatchFirstText(c)
	})
}

func NewSearchLastText(rule means.Ruler, gek GenExportKeyword) *Search[Results] {
	return GenSearchKeyword(rule, gek, func(m *Matcher, c []string) Results {
		return m.MatchLastText(c)
	})
}

func NewSearchMostText(rule means.Ruler, gek GenExportKeyword) *Search[Results] {
	return GenSearchKeyword(rule, gek, func(m *Matcher, c []string) Results {
		return m.MatchMostText(c)
	})
}

func NewSearchKey(rule means.Ruler, gek GenExportKeyword) *Search[Results] {
	return GenSearchKeyword(rule, gek, func(m *Matcher, c []string) Results {
		return m.MatchKey(c)
	})
}

func NewSearchFirstKey(rule means.Ruler, gek GenExportKeyword) *Search[Results] {
	return GenSearchKeyword(rule, gek, func(m *Matcher, c []string) Results {
		return m.MatchFirstKey(c)
	})
}

func NewSearchLastKey(rule means.Ruler, gek GenExportKeyword) *Search[Results] {
	return GenSearchKeyword(rule, gek, func(m *Matcher, c []string) Results {
		return m.MatchLastKey(c)
	})
}

func NewSearchMostKey(rule means.Ruler, gek GenExportKeyword) *Search[Results] {
	return GenSearchKeyword(rule, gek, func(m *Matcher, c []string) Results {
		return m.MatchMostKey(c)
	})
}
