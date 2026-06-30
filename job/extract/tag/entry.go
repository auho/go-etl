package tag

import (
	"github.com/auho/go-etl/v3/job/extract"
)

// NewFirstText
// the leftmost text matched
func NewFirstText(rule extract.Rule) *searchResults {
	return newSearchFirstText(rule, keywordAllToMaps, keywordAllKeys(rule), keywordAllDefaults(rule))
}

// NewMostText
// most text
func NewMostText(rule extract.Rule) *searchResults {
	return newSearchMostText(rule, keywordAllToMaps, keywordAllKeys(rule), keywordAllDefaults(rule))
}

// NewKey
// keyword
func NewKey(rule extract.Rule) *searchResults {
	return newSearchKey(rule, keywordAllToMaps, keywordAllKeys(rule), keywordAllDefaults(rule))
}

// NewFirstKey
// the first keyword matched
func NewFirstKey(rule extract.Rule) *searchResults {
	return newSearchFirstKey(rule, keywordAllToMaps, keywordAllKeys(rule), keywordAllDefaults(rule))
}

// NewMostKey
// most key
func NewMostKey(rule extract.Rule) *searchResults {
	return newSearchMostKey(rule, keywordAllToMaps, keywordAllKeys(rule), keywordAllDefaults(rule))
}

// NewWholeLabels
// merge all labels together
// label1|label2|label3
// keyword1|keyword2|keyword3|
func NewWholeLabels(rule extract.Rule) *searchLabelResults {
	return newSearchWholeLabels(rule, labelLineToMaps, labelLineKeys(rule), labelLineDefaults(rule))
}

// NewLabel
// label tags
func NewLabel(rule extract.Rule) *searchLabelResults {
	return newSearchLabels(rule, labelAllToMaps, labelAllKeys(rule), labelAllDefaults(rule))
}
