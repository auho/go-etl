package explore

import (
	"github.com/auho/go-etl/v2/job/mode"
)

var _ mode.UpdateModer = (*Update)(nil)

type Update struct {
	Explore
}

func (u *Update) Do(item map[string]any) map[string]any {
	u.AddTotal(1)

	_ticket := u.collect.Pick(item, u.search.Search)
	if !_ticket.GetOk() {
		return nil
	}

	ret := _ticket.ToExport(u.exportWay).ToTokenize()

	u.AddAmount(int64(len(ret)))

	return ret[0]
}
