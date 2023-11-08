package action

import (
	"fmt"

	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/mode"
)

var _ Actor = (*Insert)(nil)

type InsertConfig struct {
	NotTruncate bool
	ExtraKeys   []string // 附加写入到 target 的 source 字段
}

type Insert struct {
	action

	mode mode.InsertModer

	truncate  bool
	extraKeys []string // 附加写入到 target 的 source 字段
}

func WithInsertConfig(ic InsertConfig) func(*Insert) {
	return func(i *Insert) {
		i.truncate = !ic.NotTruncate
		i.extraKeys = ic.ExtraKeys
	}
}

// NewInsert
// insert
func NewInsert(target job.Target, moder mode.InsertModer, opts ...func(*Insert)) *Insert {
	i := &Insert{}
	i.mode = moder
	i.target = target

	for _, opt := range opts {
		opt(i)
	}

	i.Init()

	return i
}

// GetFields
// source data filed
func (i *Insert) GetFields() []string {
	return append(i.mode.GetFields(), i.extraKeys...)
}

func (i *Insert) Title() string {
	return fmt.Sprintf("Insert[%s] {%s}", i.target.TableName(), i.mode.GetTitle())
}

func (i *Insert) Prepare() error {
	err := i.mode.Prepare()
	if err != nil {
		return fmt.Errorf("insert action mode prepare error; %w", err)
	}

	if i.truncate {
		err = i.target.GetDB().Truncate(i.target.TableName())
		if err != nil {
			return fmt.Errorf("insert action target truncate error; %w", err)
		}
	}

	return nil
}

func (i *Insert) Do(item map[string]any) ([]map[string]any, bool) {
	newItems := i.mode.Do(item)
	if len(newItems) <= 0 {
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
