package task

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-etl/v3/job/transform"
	slices "github.com/auho/go-etl/v3/tool/slicex"
	"github.com/auho/go-toolkit-flow/v3/storage"
	"github.com/auho/go-toolkit-flow/v3/storage/database/destination"
)

var _ itemProducer = (*Update)(nil)

type UpdateConfig struct {
	NotTruncate bool
	BatchSize   int
	Concurrency int
}

func (uc *UpdateConfig) check() {
	if uc.BatchSize <= 0 {
		uc.BatchSize = batchSize
	}

	if uc.Concurrency <= 0 {
		uc.Concurrency = runtime.NumCPU()
	}
}

func WithUpdateConfig(cc UpdateConfig) func(update *Update) {
	return func(c *Update) {
		c.config = cc
	}
}

type Update struct {
	producerTask

	source job.Table
	modes  []transform.UpdateOperator

	config UpdateConfig
	dst    *destination.Bulk[storage.MapEntry]
}

func NewUpdate(source job.Table, modes []transform.UpdateOperator, opts ...func(*Update)) *Update {
	u := &Update{}
	u.source = source
	u.modes = modes

	for _, opt := range opts {
		opt(u)
	}

	u.config.check()

	return u
}

func (u *Update) GetFields() []string {
	fields := make([]string, 0)
	fields = append(fields, u.source.IDName())

	for _, m := range u.modes {
		fields = append(fields, m.GetFields()...)
	}

	fields = slices.SliceDropDuplicates(fields)

	return fields
}

func (u *Update) Summary() string {
	s := make([]string, 0)
	for _, m := range u.modes {
		s = append(s, m.Title())
	}

	return fmt.Sprintf("Update[%s] {%s}", u.source.TableName(), strings.Join(s, ", "))
}

func (u *Update) Prepare() error {
	var err error
	for _, m := range u.modes {
		err = m.Prepare()
		if err != nil {
			return fmt.Errorf("update action prepare error; %w", err)
		}
	}

	return nil
}

func (u *Update) BeforeRun() error {
	return nil
}

func (u *Update) Exec(item map[string]any) ([]map[string]any, bool, error) {
	_does := make(map[string]any)
	for _, m := range u.modes {
		_do := m.Do(item)
		for k, v := range _do {
			_does[k] = v
		}
	}

	if len(_does) <= 0 {
		return nil, false, nil
	}

	newItem := make(map[string]any)
	newItem[u.source.IDName()] = item[u.source.IDName()]

	for k, v := range _does {
		newItem[k] = v
	}

	return []map[string]any{newItem}, true, nil
}

func (u *Update) AppendState() {}

func (u *Update) AfterRun() error { return nil }

func (u *Update) Close() error {
	for _, m := range u.modes {
		err := m.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
