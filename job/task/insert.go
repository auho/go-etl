package task

import (
	"fmt"

	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-etl/v3/job/transform"
	"github.com/auho/go-toolkit-flow/v3/storage"
	"github.com/auho/go-toolkit-flow/v3/storage/database/destination"
)

type InsertConfig struct {
	baseConfig
	SkipTruncate     bool
	AllowInsertEmpty bool
	ExtraKeys        []string // source fields appended to target
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
func NewInsert(target job.Table, mode transform.InsertOperator, opts ...func(*Insert)) *Insert {
	i := &Insert{}
	i.mode = mode
	i.target = target

	for _, opt := range opts {
		opt(i)
	}

	i.config.check()

	return i
}

func (i *Insert) Fields() ([]string, error) {
	return append(i.mode.Fields(), i.config.ExtraKeys...), nil
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
			IsTruncate:  !i.config.SkipTruncate,
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
