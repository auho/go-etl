package action

import (
	"fmt"

	"github.com/auho/go-etl/v2/mode"
)

var _ Actor = (*Insert)(nil)

type Insert struct {
	action
	target    Target
	mode      mode.InsertModer
	extraKeys []string
}

// NewInsert
// insert
func NewInsert(target Target, moder mode.InsertModer, extraKeys []string) *Insert {
	i := &Insert{}
	i.mode = moder
	i.target = target
	i.extraKeys = extraKeys

	i.Init()

	return i
}

// get Keys
// target data field
//func (i *Insert) getKeys() []string {
//	return append(i.mode.GetKeys(), i.extraKeys...)
//}

// GetFields
// source data filed
func (i *Insert) GetFields() []string {
	return append(i.mode.GetFields(), i.extraKeys...)
}

func (i *Insert) Title() string {
	return fmt.Sprintf("Insert[%s] {%s}", i.target.TableName(), i.mode.GetTitle())
}

func (i *Insert) Prepare() error {
	return nil
}

func (i *Insert) Do(item map[string]any) ([]map[string]any, bool) {
	newItems := i.mode.Do(item)
	if newItems == nil {
		return nil, false
	}

	if len(i.extraKeys) > 0 {
		for index := range newItems {
			for _, key := range i.extraKeys {
				newItems[index][key] = item[key]
			}
		}
	}

	return newItems, true
}

func (i *Insert) PostBatchDo(items []map[string]any) {
	err := i.target.GetDB().BulkInsertFromSliceMap(i.target.TableName(), items, batchSize)
	if err != nil {
		panic(err)
	}
}

func (i *Insert) PostDo() {}
