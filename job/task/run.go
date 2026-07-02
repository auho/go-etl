package task

import (
	"fmt"

	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-etl/v3/tool/slicex"
	"github.com/auho/go-toolkit-flow/v3/flow"
	"github.com/auho/go-toolkit-flow/v3/storage"
	"github.com/auho/go-toolkit-flow/v3/storage/database/source"
)

type executor interface {
	options() ([]flow.Option[map[string]any, map[string]any], error)
}

func toProcessors[T processor](items []T) []processor {
	ps := make([]processor, len(items))
	for i, item := range items {
		ps[i] = item
	}
	return ps
}

func run(table job.Table, ps []processor, e executor, opts ...ConfigOption) error {
	r := &Runner{}
	r.prepare(opts)

	ds, err := r.source(table, ps)
	if err != nil {
		return fmt.Errorf("source: %w", err)
	}

	err = r.run(ds, e)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}

type Runner struct {
	config *Config
}

func (r *Runner) prepare(opts []ConfigOption) {
	r.config = &Config{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		opt(r.config)
	}

	r.config.Init()
}

func (r *Runner) source(table job.Table, ps []processor) (*source.Section[storage.MapEntry], error) {
	fields := []string{table.IDName()}
	for _, p := range ps {
		pFields, err := p.Fields()
		if err != nil {
			return nil, fmt.Errorf("SourceFields: %w", err)
		}

		fields = append(fields, pFields...)
	}

	fields = slicex.SliceDropDuplicates(fields)

	ds, err := source.NewSectionMapWithGorm(
		source.SectionConfig{
			Concurrency: r.config.source.Concurrency,
			MaxItems:    r.config.source.Maximum,
			StartID:     0,
			EndID:       0,
			PageSize:    r.config.source.PageSize,
		},
		source.ScanConfig{
			TableName:     table.TableName(),
			SegmentIDName: table.IDName(),
			Where:         "",
			Order:         "",
			SelectFields:  fields,
			WhereArgs:     nil,
		},
		table.GetDB().GormDB(),
	)

	if err != nil {
		return nil, fmt.Errorf("NewSectionMapWithGorm: %w", err)
	}

	return ds, nil
}

func (r *Runner) run(d *source.Section[storage.MapEntry], e executor) error {
	eOpts, err := e.options()
	if err != nil {
		return fmt.Errorf("options: %w", err)
	}

	opts := []flow.Option[map[string]any, map[string]any]{
		flow.WithSource[map[string]any, map[string]any](d),
	}

	opts = append(opts, eOpts...)

	err = flow.RunFlow[map[string]any](opts...)
	if err != nil {
		return fmt.Errorf("RunFlow: %w", err)
	}

	return nil
}
