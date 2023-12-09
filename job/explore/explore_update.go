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

	_export := u.collect.Do(item, u.search)
	if _export == nil || !_export.IsOk() {
		return nil
	}

	ret := _export.ToTokenize()

	u.AddAmount(int64(len(ret)))

	return ret[0]
}
