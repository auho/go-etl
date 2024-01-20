package explore

import (
	"github.com/auho/go-etl/v2/job/mode"
)

// TODO need implement
var _ mode.InsertModer = (*InsertCross)(nil)

// InsertCross
// cross means 交叉
//
// 1，2
// 3，4
// =>
// 1，3
// 1，4
// 2，3
// 2，4
type InsertCross struct {
}

func (ic *InsertCross) GetTitle() string {
	//TODO implement me
	panic("implement me")
}

func (ic *InsertCross) GetFields() []string {
	//TODO implement me
	panic("implement me")
}

func (ic *InsertCross) Prepare() error {
	//TODO implement me
	panic("implement me")
}

func (ic *InsertCross) Close() error {
	//TODO implement me
	panic("implement me")
}

func (ic *InsertCross) GetKeys() []string {
	//TODO implement me
	panic("implement me")
}

func (ic *InsertCross) DefaultValues() map[string]any {
	//TODO implement me
	panic("implement me")
}

func (ic *InsertCross) Do(contents map[string]any) []map[string]any {
	//TODO implement me
	panic("implement me")
}

func (ic *InsertCross) State() []string {
	//TODO implement me
	panic("implement me")
}
