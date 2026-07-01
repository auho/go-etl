package tag

import (
	"github.com/auho/go-etl/v3/job/extract"
)

// NewFirstText
// the leftmost text matched
func NewFirstText(rule extract.Rule) *MatcherResults {
	return newMatcherFirstText(rule, keywordAllToMaps, keywordAllKeys(rule), keywordAllDefaults(rule))
}

// NewMostText
// most text
func NewMostText(rule extract.Rule) *MatcherResults {
	return newMatcherMostText(rule, keywordAllToMaps, keywordAllKeys(rule), keywordAllDefaults(rule))
}

// NewKey
// keyword
func NewKey(rule extract.Rule) *MatcherResults {
	return newMatcherKey(rule, keywordAllToMaps, keywordAllKeys(rule), keywordAllDefaults(rule))
}

// NewFirstKey
// the first keyword matched
func NewFirstKey(rule extract.Rule) *MatcherResults {
	return newMatcherFirstKey(rule, keywordAllToMaps, keywordAllKeys(rule), keywordAllDefaults(rule))
}

// NewMostKey
// most key
func NewMostKey(rule extract.Rule) *MatcherResults {
	return newMatcherMostKey(rule, keywordAllToMaps, keywordAllKeys(rule), keywordAllDefaults(rule))
}

// NewWholeLabels
// merge all labels together
// label1|label2|label3
// keyword1|keyword2|keyword3|
func NewWholeLabels(rule extract.Rule) *MatcherLabelResults {
	return newMatcherWholeLabels(rule, labelLineToMaps, labelLineKeys(rule), labelLineDefaults(rule))
}

// NewLabel
// label tags
func NewLabel(rule extract.Rule) *MatcherLabelResults {
	return newMatcherLabels(rule, labelAllToMaps, labelAllKeys(rule), labelAllDefaults(rule))
}
