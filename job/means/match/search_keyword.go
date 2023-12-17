package match

import (
	"github.com/auho/go-etl/v2/job/means"
)

func GenSearchKeyword(rule means.Ruler, gek GenExportKeyword, ste SearchResults[Results]) *Search[Results] {
	return NewSearch[Results](rule, gek(rule), ste)
}

func NewSearchFirstText(rule means.Ruler, gek GenExportKeyword) *Search[Results] {
	return GenSearchKeyword(rule, gek, func(ctx *SearchContext[Results], c []string) Results {
		return ctx.Matcher.MatchFirstText(c)
	})
}

func NewSearchLastText(rule means.Ruler, gek GenExportKeyword) *Search[Results] {
	return GenSearchKeyword(rule, gek, func(ctx *SearchContext[Results], c []string) Results {
		return ctx.Matcher.MatchLastText(c)
	})
}

func NewSearchMostText(rule means.Ruler, gek GenExportKeyword) *Search[Results] {
	return GenSearchKeyword(rule, gek, func(ctx *SearchContext[Results], c []string) Results {
		return ctx.Matcher.MatchMostText(c)
	})
}

func NewSearchKey(rule means.Ruler, gek GenExportKeyword) *Search[Results] {
	return GenSearchKeyword(rule, gek, func(ctx *SearchContext[Results], c []string) Results {
		return ctx.Matcher.MatchKey(c)
	})
}

func NewSearchFirstKey(rule means.Ruler, gek GenExportKeyword) *Search[Results] {
	return GenSearchKeyword(rule, gek, func(ctx *SearchContext[Results], c []string) Results {
		return ctx.Matcher.MatchFirstKey(c)
	})
}

func NewSearchLastKey(rule means.Ruler, gek GenExportKeyword) *Search[Results] {
	return GenSearchKeyword(rule, gek, func(ctx *SearchContext[Results], c []string) Results {
		return ctx.Matcher.MatchLastKey(c)
	})
}

func NewSearchMostKey(rule means.Ruler, gek GenExportKeyword) *Search[Results] {
	return GenSearchKeyword(rule, gek, func(ctx *SearchContext[Results], c []string) Results {
		return ctx.Matcher.MatchMostKey(c)
	})
}
