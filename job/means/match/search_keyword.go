package match

func NewSearchKeyword(export *ExportResults, ste SearchResultsFun[Results]) *SearchResults {
	return NewSearch[Results](export, ste)
}

func NewSearchFirstText(export *ExportResults) *SearchResults {
	return NewSearchKeyword(export, func(ctx *SearchContextResults, c []string) Results {
		return ctx.Matcher.MatchFirstText(c)
	})
}

func NewSearchLastText(export *ExportResults) *SearchResults {
	return NewSearchKeyword(export, func(ctx *SearchContextResults, c []string) Results {
		return ctx.Matcher.MatchLastText(c)
	})
}

func NewSearchMostText(export *ExportResults) *SearchResults {
	return NewSearchKeyword(export, func(ctx *SearchContextResults, c []string) Results {
		return ctx.Matcher.MatchMostText(c)
	})
}

func NewSearchKey(export *ExportResults) *SearchResults {
	return NewSearchKeyword(export, func(ctx *SearchContextResults, c []string) Results {
		return ctx.Matcher.MatchKey(c)
	})
}

func NewSearchFirstKey(export *ExportResults) *SearchResults {
	return NewSearchKeyword(export, func(ctx *SearchContextResults, c []string) Results {
		return ctx.Matcher.MatchFirstKey(c)
	})
}

func NewSearchLastKey(export *ExportResults) *SearchResults {
	return NewSearchKeyword(export, func(ctx *SearchContextResults, c []string) Results {
		return ctx.Matcher.MatchLastKey(c)
	})
}

func NewSearchMostKey(export *ExportResults) *SearchResults {
	return NewSearchKeyword(export, func(ctx *SearchContextResults, c []string) Results {
		return ctx.Matcher.MatchMostKey(c)
	})
}
