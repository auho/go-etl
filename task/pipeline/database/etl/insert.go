package etl

import (
	"fmt"

	"github.com/auho/go-etl/v3/task/transform"
	"github.com/auho/go-toolkit-flow/v3/processor/producer"
	"github.com/auho/go-toolkit-flow/v3/storage"
	"github.com/auho/go-toolkit-flow/v3/storage/database/destination"
)

var _ producer.Item[storage.MapEntry, storage.MapEntry] = (*Insert)(nil)

// InsertConfig configures an Insert task.
type InsertConfig struct {
	BaseConfig
	SkipTruncate     bool
	AllowInsertEmpty bool
	ExtraKeys        []string // source fields appended to target
}

// WithInsertConfig returns an option that sets the Insert config.
func WithInsertConfig(ic InsertConfig) func(*Insert) {
	return func(i *Insert) {
		i.config = ic
	}
}

// Insert is a producer that transforms each source row via operator and inserts
// the result into the target table.
type Insert struct {
	producer.Processor
	TaskBase

	operator transform.InsertOperator
	target   Table
	config   InsertConfig
}

// NewInsert creates an Insert producer writing to target via operator.
func NewInsert(target Table, operator transform.InsertOperator, opts ...func(*Insert)) *Insert {
	i := &Insert{}
	i.operator = operator
	i.target = target

	for _, opt := range opts {
		opt(i)
	}

	i.config.Check()

	return i
}

func (i *Insert) Fields() ([]string, error) {
	return append(i.operator.Fields(), i.config.ExtraKeys...), nil
}

func (i *Insert) Summary() string {
	return fmt.Sprintf("Insert[%s] {%s}", i.target.TableName(), i.operator.Title())
}

func (i *Insert) Prepare() error {
	err := i.operator.Prepare()
	if err != nil {
		return fmt.Errorf("operator.Prepare: %w", err)
	}

	return nil
}

func (i *Insert) BeforeRun() error { return nil }

func (i *Insert) Exec(item storage.MapEntry) ([]storage.MapEntry, bool, error) {
	newItems, err := i.operator.Apply(item)
	if err != nil {
		return nil, false, fmt.Errorf("operator.Apply: %w", err)
	}
	if len(newItems) <= 0 {
		if i.config.AllowInsertEmpty {
			newItems = []map[string]any{i.operator.DefaultValues()}
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
	return i.operator.Close()
}

// BuildDestinations constructs and returns the insert destination for the target table.
func (i *Insert) BuildDestinations() ([]storage.Destination[storage.MapEntry], error) {
	dest, err := destination.NewBulkInsertMapWithGorm(
		i.config.BulkConfig(!i.config.SkipTruncate),
		destination.WriteConfig{TableName: i.target.TableName()},
		i.target.GetDB().GormDB(),
	)
	if err != nil {
		return nil, fmt.Errorf("NewBulkInsertMapWithGorm: %w", err)
	}

	return []storage.Destination[storage.MapEntry]{dest}, nil
}
