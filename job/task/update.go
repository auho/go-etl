package task

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-etl/v3/job/transform"
	"github.com/auho/go-etl/v3/tool/slicex"
	"github.com/auho/go-toolkit-flow/v3/storage"
	"github.com/auho/go-toolkit-flow/v3/storage/database/destination"
)

var _ itemProducer = (*Update)(nil)

type UpdateConfig struct {
	baseConfig
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

func (u *Update) Fields() ([]string, error) {
	fields := make([]string, 0)
	fields = append(fields, u.source.IDName())

	for _, m := range u.modes {
		fields = append(fields, m.Fields()...)
	}

	fields = slicex.SliceDropDuplicates(fields)

	return fields, nil
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
			return fmt.Errorf("prepare: %w", err)
		}
	}

	return nil
}

func (u *Update) BeforeRun() error {
	return nil
}

func (u *Update) Exec(item map[string]any) ([]map[string]any, bool, error) {
	does := make(map[string]any)
	for _, m := range u.modes {
		do, err := m.Apply(item)
		if err != nil {
			return nil, false, fmt.Errorf("apply: %w", err)
		}
		for k, v := range do {
			does[k] = v
		}
	}

	if len(does) <= 0 {
		return nil, false, nil
	}

	newItem := make(map[string]any)
	newItem[u.source.IDName()] = item[u.source.IDName()]

	for k, v := range does {
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

func (u *Update) destinations() ([]storage.Destination[storage.MapEntry], error) {
	var ds []storage.Destination[storage.MapEntry]

	target, err := u.source.GetDB().Clone()
	if err != nil {
		return nil, fmt.Errorf("GetDB.Clone: %w", err)
	}

	dest, err := destination.NewBulkUpdateMapWithGorm(
		destination.BulkConfig{
			IsTruncate:  false,
			Concurrency: u.config.Concurrency,
			PageSize:    int64(u.config.BatchSize),
		},
		destination.WriteConfig{
			TableName: u.source.TableName(),
		},
		u.source.IDName(), target.GormDB(),
	)
	if err != nil {
		return nil, fmt.Errorf("NewBulkUpdateMapWithGorm: %w", err)
	}

	ds = append(ds, dest)
	return ds, nil
}
