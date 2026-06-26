package transform

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collect"
	"github.com/auho/go-etl/v3/job/transform/filter"
)

var _ UpdateOperator = (*Update)(nil)

type Update struct {
	*Pipeline
}

func newUpdateFromPipeline(e *Pipeline) *Update {
	return NewUpdate(e.collect, e.search, e.condition)
}

func NewUpdate(collect collect.Collector, search extract.Extractor, expression filter.Predicate) *Update {
	return &Update{
		Pipeline: newPipeline(collect, search, expression),
	}
}

func (u *Update) Apply(item map[string]any) map[string]any {
	u.AddTotal(1)

	if !u.expressionOperation(item) {
		return nil
	}

	token := u.collect.Search(item, u.search)
	if !token.IsOK() {
		return nil
	}

	ret := token.Rows()

	u.AddAmount(int64(len(ret)))

	return ret[0]
}
