package task

import (
	"fmt"
	"runtime"

	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-etl/v3/job/transform"
	"github.com/auho/go-toolkit-flow/v3/storage"
	"github.com/auho/go-toolkit-flow/v3/storage/database/destination"
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

	mode   transform.InsertOperator
	target job.Table
	config InsertConfig
}

// NewInsert
// insert
func NewInsert(target job.Table, moder transform.InsertOperator, opts ...func(*Insert)) *Insert {
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
	return fmt.Sprintf("Insert[%s] {%s}", i.target.TableName(), i.mode.Title())
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
	newItems, err := i.mode.Apply(item)
	if err != nil {
		return nil, false, fmt.Errorf("mode.Apply: %w", err)
	}
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

func (i *Insert) destinations() ([]storage.Destination[storage.MapEntry], error) {
	var ds []storage.Destination[storage.MapEntry]
	dest, err := destination.NewBulkInsertMapWithGorm(
		destination.BulkConfig{
			IsTruncate:  !i.config.NotTruncate,
			Concurrency: i.config.Concurrency,
			PageSize:    int64(i.config.BatchSize),
		},
		destination.WriteConfig{
			TableName: i.target.TableName(),
		},
		i.target.GetDB().GormDB(),
	)
	if err != nil {
		return nil, fmt.Errorf("NewBulkInsertMapWithGorm: %w", err)
	}

	ds = append(ds, dest)
	return ds, nil
}
