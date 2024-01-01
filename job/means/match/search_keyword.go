package match

import (
	"github.com/auho/go-etl/v2/job/means"
)

func NewSearchKeyword(rule means.Ruler, export *ExportResults, ste SearchResultsFun[Results]) *SearchResults {
	return NewSearch[Results](rule, export, ste)
}

func NewSearchFirstText(rule means.Ruler, export *ExportResults) *SearchResults {
	return NewSearchKeyword(rule, export, func(ctx *SearchContextResults, c []string) Results {
		return ctx.Matcher.MatchFirstText(c)
	})
}

func NewSearchLastText(rule means.Ruler, export *ExportResults) *SearchResults {
	return NewSearchKeyword(rule, export, func(ctx *SearchContextResults, c []string) Results {
		return ctx.Matcher.MatchLastText(c)
	})
}

func NewSearchMostText(rule means.Ruler, export *ExportResults) *SearchResults {
	return NewSearchKeyword(rule, export, func(ctx *SearchContextResults, c []string) Results {
		return ctx.Matcher.MatchMostText(c)
	})
}

func NewSearchKey(rule means.Ruler, export *ExportResults) *SearchResults {
	return NewSearchKeyword(rule, export, func(ctx *SearchContextResults, c []string) Results {
		return ctx.Matcher.MatchKey(c)
	})
}

func NewSearchFirstKey(rule means.Ruler, export *ExportResults) *SearchResults {
	return NewSearchKeyword(rule, export, func(ctx *SearchContextResults, c []string) Results {
		return ctx.Matcher.MatchFirstKey(c)
	})
}

func NewSearchLastKey(rule means.Ruler, export *ExportResults) *SearchResults {
	return NewSearchKeyword(rule, export, func(ctx *SearchContextResults, c []string) Results {
		return ctx.Matcher.MatchLastKey(c)
	})
}

func NewSearchMostKey(rule means.Ruler, export *ExportResults) *SearchResults {
	return NewSearchKeyword(rule, export, func(ctx *SearchContextResults, c []string) Results {
		return ctx.Matcher.MatchMostKey(c)
	})
}
