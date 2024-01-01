package tag

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Exporter = (*Export[Results])(nil)
var _ search.Exporter = (*Export[LabelResults])(nil)
var _ search.Exporter = (*ExportResults)(nil)
var _ search.Exporter = (*ExportLabelResults)(nil)

type ExportResults = Export[Results]
type ExportLabelResults = Export[LabelResults]

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
	token := search.Token{}

	if len(results) > 0 {
		token.SetOk()
		token.SetTokenizerFunc(func() []map[string]any {
			return e.resultsToToken(results, rule)
		})
	}

	return token
}
