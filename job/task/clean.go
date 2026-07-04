package task

import (
	"fmt"
	"maps"
	"strings"

	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-etl/v3/job/transform"
	"github.com/auho/go-etl/v3/tool/slicex"
	"github.com/auho/go-toolkit-flow/v3/processor/consumer"
	"github.com/auho/go-toolkit-flow/v3/storage"
	"github.com/auho/go-toolkit-flow/v3/storage/database/destination"
)

// CleanConfig configures a Clean task.
type CleanConfig struct {
	baseConfig
	SkipTruncate bool
	ExtraTags    bool     // tags to deleted data
	Keys         []string // source columns name, priority of use this keys
}

// WithCleanConfig returns an option that sets the Clean config.
func WithCleanConfig(cc CleanConfig) func(*Clean) {
	return func(c *Clean) {
		c.config = cc
	}
}

var _ itemConsumer = (*Clean)(nil)
var _ consumer.DestinationHolder[storage.MapEntry] = (*Clean)(nil)

// Clean is a consumer that filters source rows: rows matching any operator are
// routed to the deleted table, the rest to the data table.
type Clean struct {
	consumerTask

	resource  job.CleanResource
	operators []transform.UpdateOperator

	config CleanConfig

	dataDest    *destination.Bulk[storage.MapEntry]
	deletedDest *destination.Bulk[storage.MapEntry]
}

// NewClean creates a Clean consumer over the given resource and operators.
func NewClean(cr job.CleanResource, operators []transform.UpdateOperator, opts ...func(*Clean)) *Clean {
	c := &Clean{}
	c.resource = cr
	c.operators = operators

	for _, opt := range opts {
		opt(c)
	}

	c.config.check()

	return c
}

func (c *Clean) Fields() ([]string, error) {
	var keys []string
	if len(c.config.Keys) > 0 {
		keys = append(keys, c.config.Keys...)
		for _, op := range c.operators {
			keys = append(keys, op.Fields()...)
		}

		keys = slicex.SliceDropDuplicates(keys)
	} else {
		var err error
		keys, err = c.resource.Deleted().GetDB().GetTableColumns(c.resource.Deleted().TableName())
		if err != nil {
			return nil, fmt.Errorf("GetTableColumns: %w", err)
		}
	}

	return keys, nil
}

func (c *Clean) Summary() string {
	s := make([]string, 0)
	for _, op := range c.operators {
		s = append(s, op.Title())
	}

	return fmt.Sprintf("Clean[%s, %s] {%s}",
		c.resource.Data().TableName(),
		c.resource.Deleted().TableName(),
		strings.Join(s, ", "))
}

func (c *Clean) Prepare() error {
	var err error
	for _, op := range c.operators {
		err = op.Prepare()
		if err != nil {
			return fmt.Errorf("prepare: %w", err)
		}
	}

	bConfig := destination.BulkConfig{
		IsTruncate:  !c.config.SkipTruncate,
		Concurrency: c.config.Concurrency,
		PageSize:    int64(c.config.BatchSize),
	}

	c.dataDest, err = destination.NewBulkInsertMapWithGorm(
		bConfig,
		destination.WriteConfig{TableName: c.resource.Data().TableName()},
		c.resource.Data().GetDB().GormDB(),
	)
	if err != nil {
		return fmt.Errorf("NewBulkInsertMapWithGorm DataTarget: %w", err)
	}

	c.deletedDest, err = destination.NewBulkInsertMapWithGorm(
		bConfig,
		destination.WriteConfig{TableName: c.resource.Deleted().TableName()},
		c.resource.Deleted().GetDB().GormDB(),
	)
	if err != nil {
		return fmt.Errorf("NewBulkInsertMapWithGorm DeletedTarget: %w", err)
	}

	return nil
}

func (c *Clean) BeforeRun() error {
	return nil
}

func (c *Clean) Exec(item map[string]any) (bool, error) {
	needDeleted := false
	for _, op := range c.operators {
		res, err := op.Apply(item)
		if err != nil {
			return false, fmt.Errorf("apply: %w", err)
		}

		if len(res) > 0 {
			needDeleted = true

			if c.config.ExtraTags {
				maps.Copy(item, res)
			}

			break
		}
	}

	var err error

	if needDeleted {
		err = c.deletedDest.Receive([]map[string]any{item})
		if err != nil {
			return false, fmt.Errorf("deletedDest.Receive: %w", err)
		}
	} else {
		err = c.dataDest.Receive([]map[string]any{item})
		if err != nil {
			return false, fmt.Errorf("dataDest.Receive: %w", err)
		}
	}

	return true, nil
}

func (c *Clean) AfterRun() error {
	return nil
}

func (c *Clean) AppendState() {}

func (c *Clean) Close() error {
	for _, op := range c.operators {
		err := op.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Clean) Destinations() ([]storage.Destination[storage.MapEntry], error) {
	return []storage.Destination[storage.MapEntry]{c.dataDest, c.deletedDest}, nil
}
