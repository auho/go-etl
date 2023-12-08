package search

import (
	"github.com/auho/go-etl/v2/insight/assistant/model"
	"github.com/auho/go-etl/v2/job/explore/token"
)

type Searcher interface {
	GetTitle() string
	Search(s []string) token.Tokenizer
	DefaultTokenize() token.Tokenizer
}

type Way interface {
}

type Search struct {
	rule *model.Rule
}
