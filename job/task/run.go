package task

import (
	"fmt"

	"github.com/auho/go-etl/v2/job"
	slices "github.com/auho/go-etl/v2/tool/slicex"
	"github.com/auho/go-toolkit-flow/flow"
	"github.com/auho/go-toolkit-flow/storage"
	"github.com/auho/go-toolkit-flow/storage/database/source"
)

type executor interface {
	exec() []flow.Option[map[string]any, map[string]any]
}

type Runner struct {
	config *Config
}

func run(jb job.Table, ps []processor, e executor, opts ...ConfigOption) {
	r := &Runner{}
	r.prepare(opts)
	ds, err := r.source(jb, ps)
	if err != nil {
		panic(fmt.Sprintf("source: %v", err))
	}

	err = r.run(ds, e)
	if err != nil {
		panic(fmt.Sprintf("run: %v", err))
	}
}

func (r *Runner) prepare(opts []ConfigOption) {
	r.config = &Config{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		opt(r.config)
	}

	r.config.Check()
}

func (r *Runner) source(s job.Table, ps []processor) (*source.Section[storage.MapEntry], error) {
	fields := []string{s.GetIDName()}
	for _, p := range ps {
		fields = append(fields, p.GetFields()...)
	}

	fields = slices.SliceDropDuplicates(fields)

	ds, err := source.NewSectionMapWithGorm(
		source.SectionConfig{
			Concurrency: r.config.source.Concurrency,
			MaxItems:    r.config.source.Maximum,
			StartID:     0,
			EndID:       0,
			PageSize:    r.config.source.PageSize,
		},
		source.ScanConfig{
			TableName:     s.TableName(),
			SegmentIDName: s.GetIDName(),
			Where:         "",
			Order:         "",
			SelectFields:  fields,
			WhereArgs:     nil,
		},
		s.GetDB().GormDB(),
	)

	if err != nil {
		return nil, fmt.Errorf("source.NewSectionMapWithGorm: %w", err)
	}

	return ds, nil
}

func (r *Runner) run(d *source.Section[storage.MapEntry], e executor) error {
	opts := []flow.Option[map[string]any, map[string]any]{
		flow.WithSource[map[string]any, map[string]any](d),
	}

	opts = append(opts, e.exec()...)

	err := flow.RunFlow[map[string]any](opts...)
	if err != nil {
		return fmt.Errorf("flow.RunFlow: %w", err)
	}

	return nil
}
