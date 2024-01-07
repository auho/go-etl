package match

import (
	"github.com/auho/go-etl/v2/job/means"
)

// NewFirstText
// the leftmost text matched
func NewFirstText(rule means.Ruler) *SearchResults {
	return NewSearchFirstText(rule, NewExportKeywordAll(rule))
}

// NewMostText
// most text
func NewMostText(rule means.Ruler) *SearchResults {
	return NewSearchMostText(rule, NewExportKeywordAll(rule))
}

// NewKey
// keyword
func NewKey(rule means.Ruler) *SearchResults {
	return NewSearchKey(rule, NewExportKeywordAll(rule))
}

// NewFirstKey
// the first keyword matched
func NewFirstKey(rule means.Ruler) *SearchResults {
	return NewSearchFirstKey(rule, NewExportKeywordAll(rule))
}

// NewMostKey
// most key
func NewMostKey(rule means.Ruler) *SearchResults {
	return NewSearchMostKey(rule, NewExportKeywordAll(rule))
}

// NewWholeLabels
// merge all labels together
// label1|label2|label3
// keyword1|keyword2|keyword3|
func NewWholeLabels(rule means.Ruler) *SearchLabelResults {
	return NewSearchWholeLabels(rule, NewExportLabelLine(rule))
}

// NewLabel
// label tags
func NewLabel(rule means.Ruler) *SearchLabelResults {
	return NewSearchLabels(rule, NewExportLabelAll(rule))
}
