package task

import (
	"fmt"
	"runtime"

	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/mode"
)

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

var _ itemProducer = (*Insert)(nil)

type Insert struct {
	producerTask

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

	i.config.check()

	return i
}

// GetFields
// source data filed
func (i *Insert) GetFields() []string {
	return append(i.mode.GetFields(), i.config.ExtraKeys...)
}

func (i *Insert) Summary() string {
	return fmt.Sprintf("Insert[%s] {%s}", i.target.TableName(), i.mode.GetTitle())
}

func (i *Insert) Prepare() error {
	err := i.mode.Prepare()
	if err != nil {
		return fmt.Errorf("mode.Prepare: %w", err)
	}

	return nil
}

func (i *Insert) BeforeRun() error {
	return nil
}

func (i *Insert) Exec(item map[string]any) ([]map[string]any, bool, error) {
	newItems := i.mode.Do(item)
	if len(newItems) <= 0 {
		if i.config.AllowInsertEmpty {
			newItems = []map[string]any{i.mode.DefaultValues()}
		} else {
			return nil, false, nil
		}
	}

	if len(i.config.ExtraKeys) > 0 {
		for index := range newItems {
			for _, key := range i.config.ExtraKeys {
				newItems[index][key] = item[key]
			}
		}
	}

	return newItems, true, nil
}

func (i *Insert) AfterRun() error { return nil }

func (i *Insert) AppendState() {}

func (i *Insert) Close() error {
	return i.mode.Close()
}
