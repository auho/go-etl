package search

import (
	"github.com/auho/go-etl/v2/job/explore/token"
)

type Searcher interface {
	GetTitle() string
	Search(s string) token.Tokenizer
	DefaultTokenize() token.Tokenizer
}
