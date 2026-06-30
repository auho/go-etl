package match

import (
	"github.com/auho/go-etl/v3/job/extract"
)

// NewFirstText
// the leftmost text matched
func NewFirstText(rule extract.Rule) *MatcherResults {
	return newMatcherFirstText(rule, keywordToMapsAll, keywordKeysAll)
}

// NewMostText
// most text
func NewMostText(rule extract.Rule) *MatcherResults {
	return newMatcherMostText(rule, keywordToMapsAll, keywordKeysAll)
}

// NewKey
// keyword
func NewKey(rule extract.Rule) *MatcherResults {
	return newMatcherKey(rule, keywordToMapsAll, keywordKeysAll)
}

// NewFirstKey
// the first keyword matched
func NewFirstKey(rule extract.Rule) *MatcherResults {
	return newMatcherFirstKey(rule, keywordToMapsAll, keywordKeysAll)
}

// NewMostKey
// most key
func NewMostKey(rule extract.Rule) *MatcherResults {
	return newMatcherMostKey(rule, keywordToMapsAll, keywordKeysAll)
}

// NewWholeLabels
// merge all labels together
// label1|label2|label3
// keyword1|keyword2|keyword3|
func NewWholeLabels(rule extract.Rule) *MatcherLabelResults {
	return newMatcherWholeLabels(rule, labelToMapsLine, labelKeysLine)
}

// NewLabel
// label tags
func NewLabel(rule extract.Rule) *MatcherLabelResults {
	return newMatcherLabels(rule, labelToMapsAll, labelKeysAll)
}
