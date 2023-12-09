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

	_export := i.collect.Do(item, i.search)
	if _export == nil || !_export.IsOk() {
		return nil
	}

	ret := _export.ToTokenize()

	i.AddAmount(int64(len(ret)))

	return ret
}
