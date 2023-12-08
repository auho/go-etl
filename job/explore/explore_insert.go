package explore

import (
	"github.com/auho/go-etl/v2/job/mode"
)

var _ mode.InsertModer = (*Insert)(nil)

type Insert struct {
	Explore
}

func (i *Insert) Do(item map[string]any) []map[string]any {
	i.AddTotal(1)

	_ticket := i.collect.Pick(item, i.search.Search)
	if !_ticket.GetOk() {
		return nil
	}

	ret := _ticket.ToExport(i.exportWay).ToTokenize()

	i.AddAmount(int64(len(ret)))

	return ret
}
