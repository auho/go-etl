package match

func NewSearchLabel(export *ExportLabelResults, srf SearchResultsFunc[LabelResults]) *SearchLabelResults {
	return NewSearch[LabelResults](export, srf)
}

func NewSearchWholeLabels(export *ExportLabelResults) *SearchLabelResults {
	return NewSearchLabel(export, func(ctx *SearchContextLabelResults, c []string) LabelResults {
		return ctx.Matcher.MatchLabel(c)
	})
}

func NewSearchLabels(export *ExportLabelResults) *SearchLabelResults {
	return NewSearchLabel(export, func(ctx *SearchContextLabelResults, c []string) LabelResults {
		return ctx.Matcher.MatchLabel(c)
	})
}
