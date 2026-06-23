package match

import (
	"github.com/auho/go-etl/v2/job/extract"
)

// NewFirstText
// the leftmost text matched
func NewFirstText(rule means.Rule) *SearchResults {
	return NewSearchFirstText(NewExportKeywordAll(rule))
}

// NewMostText
// most text
func NewMostText(rule means.Rule) *SearchResults {
	return NewSearchMostText(NewExportKeywordAll(rule))
}

// NewKey
// keyword
func NewKey(rule means.Rule) *SearchResults {
	return NewSearchKey(NewExportKeywordAll(rule))
}

// NewFirstKey
// the first keyword matched
func NewFirstKey(rule means.Rule) *SearchResults {
	return NewSearchFirstKey(NewExportKeywordAll(rule))
}

// NewMostKey
// most key
func NewMostKey(rule means.Rule) *SearchResults {
	return NewSearchMostKey(NewExportKeywordAll(rule))
}

// NewWholeLabels
// merge all labels together
// label1|label2|label3
// keyword1|keyword2|keyword3|
func NewWholeLabels(rule means.Rule) *SearchLabelResults {
	return NewSearchWholeLabels(NewExportLabelLine(rule))
}

// NewLabel
// label tags
func NewLabel(rule means.Rule) *SearchLabelResults {
	return NewSearchLabels(NewExportLabelAll(rule))
}
