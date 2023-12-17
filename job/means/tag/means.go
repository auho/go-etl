package tag

import (
	"github.com/auho/go-etl/v2/job/means"
)

// NewFirstText
// the leftmost text matched
func NewFirstText(rule means.Ruler) *Search[Results] {
	return NewSearchFirstText(rule, NewExportKeywordAll)
}

// NewKey
// keyword
func NewKey(rule means.Ruler) *Search[Results] {
	return NewSearchKey(rule, NewExportKeywordAll)
}

// NewFirstKey
// the first keyword matched
func NewFirstKey(rule means.Ruler) *Search[Results] {
	return NewSearchFirstKey(rule, NewExportKeywordAll)
}

// NewWholeLabels
// merge all labels together
// label1|label2|label3
// keyword1|keyword2|keyword3|
func NewWholeLabels(rule means.Ruler) *Search[LabelResults] {
	return NewSearchWholeLabels(rule)
}

// NewLabel
// label tags
func NewLabel(rule means.Ruler) *Search[LabelResults] {
	return NewSearchLabels(rule, NewExportLabelAll)
}
