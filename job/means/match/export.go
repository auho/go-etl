package match

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Exporter = (*Export[Results])(nil)
var _ search.Exporter = (*Export[LabelResults])(nil)

type Export[T ResultsEntity] struct {
	resultsToToken func(T, means.Ruler) []map[string]any

	keys          []string
	defaultValues map[string]any
}

func (e *Export[T]) GetKeys() []string {
	return e.keys
}

func (e *Export[T]) GetDefaultValues() map[string]any {
	return e.defaultValues
}

func (e *Export[T]) ToToken(results T, rule means.Ruler) search.Token {
	token := search.Token{
		Ok:        false,
		Tokenizer: nil,
	}

	if len(results) <= 0 {
		return token
	}

	return search.Token{
		Ok: true,
		Tokenizer: func() []map[string]any {
			return e.resultsToToken(results, rule)
		},
	}
}
