package action

import (
	"fmt"
	"runtime"

	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/mode"
)

var _ Actor = (*Insert)(nil)

type InsertConfig struct {
	NotTruncate      bool
	BatchSize        int
	Concurrency      int
	AllowInsertEmpty bool
	ExtraKeys        []string // 附加写入到 target 的 source 字段
}

func (ic *InsertConfig) check() {
	if ic.BatchSize <= 0 {
		ic.BatchSize = batchSize
	}

	if ic.Concurrency <= 0 {
		ic.Concurrency = runtime.NumCPU()
	}
}

func WithInsertConfig(ic InsertConfig) func(*Insert) {
	return func(i *Insert) {
		i.config = ic
	}
}

type Insert struct {
	TargetAction

	mode mode.InsertModer

	config InsertConfig
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

	i.config.check()

	return i
}

// GetFields
// source data filed
func (i *Insert) GetFields() []string {
	return append(i.mode.GetFields(), i.config.ExtraKeys...)
}

func (i *Insert) Title() string {
	return fmt.Sprintf("Insert[%s] {%s}", i.target.TableName(), i.mode.GetTitle())
}

func (i *Insert) Prepare() error {
	err := i.mode.Prepare()
	if err != nil {
		return fmt.Errorf("insert action mode prepare error; %w", err)
	}

	return nil
}

func (i *Insert) PreDo() error {
	if !i.config.NotTruncate {
		err := i.target.GetDB().Truncate(i.target.TableName())
		if err != nil {
			return fmt.Errorf("insert action target truncate error; %w", err)
		}
	}

	return nil
}

func (i *Insert) Do(item map[string]any) ([]map[string]any, bool) {
	newItems := i.mode.Do(item)
	if len(newItems) <= 0 {
		if i.config.AllowInsertEmpty {
			newItems = []map[string]any{i.mode.DefaultValues()}
		} else {
			return nil, false
		}
	}

	if len(i.config.ExtraKeys) > 0 {
		for index := range newItems {
			for _, key := range i.config.ExtraKeys {
				newItems[index][key] = item[key]
			}
		}
	}

	return newItems, true
}

func (i *Insert) PostBatchDo(items []map[string]any) {
	err := i.target.GetDB().BulkInsertFromSliceMap(i.target.TableName(), items, batchSize)
	if err != nil {
		s := err.Error()
		if len(s) > 300 {
			s = s[0:300]
		}

		panic(fmt.Errorf("%s; %s", i.Title(), s))
	}
}

func (i *Insert) Blink()        {}
func (i *Insert) PostDo() error { return nil }
func (i *Insert) Close() error {
	for _, s := range i.mode.State() {
		i.Println(s)
	}

	i.Println("")

	return i.mode.Close()
}
