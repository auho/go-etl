package task

import (
	"fmt"
	"maps"
	"runtime"
	"strings"

	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/mode"
	slices "github.com/auho/go-etl/v2/tool/slicex"
	"github.com/auho/go-toolkit-flow/processor/consumer"
	"github.com/auho/go-toolkit-flow/storage"
	"github.com/auho/go-toolkit-flow/storage/database/destination"
)

type CleanConfig struct {
	NotTruncate  bool
	AddExtraTags bool // tags to deleted data
	BatchSize    int
	Concurrency  int
	Keys         []string // source columns name，priority of use this keys
}

func (cc *CleanConfig) check() {
	if cc.BatchSize <= 0 {
		cc.BatchSize = batchSize
	}

	if cc.Concurrency <= 0 {
		cc.Concurrency = runtime.NumCPU()
	}
}

func WithCleanConfig(cc CleanConfig) func(*Clean) {
	return func(c *Clean) {
		c.config = cc
	}
}

var _ itemConsumer = (*Clean)(nil)
var _ consumer.DestinationHolder[storage.MapEntry] = (*Clean)(nil)

// Clean
// filter
type Clean struct {
	consumerTask

	cleanTarget job.CleanResource
	keys        []string
	modes       []mode.UpdateModer

	config CleanConfig

	dataDest    *destination.Bulk[storage.MapEntry]
	deletedDest *destination.Bulk[storage.MapEntry]

	dataDstLine    int
	deletedDstLine int
}

func NewClean(cr job.CleanResource, modes []mode.UpdateModer, opts ...func(*Clean)) *Clean {
	c := &Clean{}
	c.cleanTarget = cr
	c.modes = modes

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Clean) GetFields() []string {
	if len(c.config.Keys) > 0 {
		c.keys = append(c.keys, c.config.Keys...)
		for _, m := range c.modes {
			c.keys = append(c.keys, m.GetFields()...)
		}

		c.keys = slices.SliceDropDuplicates(c.keys)

	} else {
		var err error
		c.keys, err = c.cleanTarget.Deleted().GetDB().GetTableColumns(c.cleanTarget.Deleted().TableName())
		if err != nil {
			panic(fmt.Errorf("GetTableColumns error; %w", err))
		}
	}

	return c.keys
}

func (c *Clean) Summary() string {
	s := make([]string, 0)
	for _, m := range c.modes {
		s = append(s, m.GetTitle())
	}

	return fmt.Sprintf("Clean[%s, %s] {%s}",
		c.cleanTarget.Data().TableName(),
		c.cleanTarget.Deleted().TableName(),
		strings.Join(s, ", "))
}

func (c *Clean) Prepare() error {
	var err error
	for _, m := range c.modes {
		err = m.Prepare()
		if err != nil {
			return fmt.Errorf("prepare: %w", err)
		}
	}

	bConfig := destination.BulkConfig{
		IsTruncate: true,
		PageSize:   batchSize,
	}

	c.dataDest, err = destination.NewBulkInsertMapWithGorm(
		bConfig,
		destination.WriteConfig{TableName: c.cleanTarget.Data().TableName()},
		c.cleanTarget.Data().GetDB().GormDB(),
	)
	if err != nil {
		return fmt.Errorf("NewBulkInsertMapWithGorm DataTarget: %w", err)
	}

	c.deletedDest, err = destination.NewBulkInsertMapWithGorm(
		bConfig,
		destination.WriteConfig{TableName: c.cleanTarget.Deleted().TableName()},
		c.cleanTarget.Deleted().GetDB().GormDB(),
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
	_needDeleted := false
	for _, m := range c.modes {
		_res := m.Do(item)
		if len(_res) > 0 {
			_needDeleted = true

			if c.config.AddExtraTags {
				maps.Copy(item, _res)
			}

			break
		}
	}

	var err error

	if _needDeleted {
		err = c.deletedDest.Receive([]map[string]any{item})
		if err != nil {
			return false, fmt.Errorf("deletedDest.Receive")
		}
	} else {
		err = c.dataDest.Receive([]map[string]any{item})
		if err != nil {
			return false, fmt.Errorf("dataDest.Receive")
		}
	}

	return true, nil
}

func (c *Clean) AfterRun() error {
	return nil
}

func (c *Clean) AppendState() {}

func (c *Clean) Close() error {
	for _, m := range c.modes {
		err := m.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Clean) Destinations() []storage.Destination[storage.MapEntry] {
	return []storage.Destination[storage.MapEntry]{c.dataDest, c.deletedDest}
}
