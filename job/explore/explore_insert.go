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

	token := i.collect.Do(item, i.search)
	if !token.IsOk() {
		return nil
	}

	ret := token.ToToken()

	i.AddAmount(int64(len(ret)))

	return ret
}
