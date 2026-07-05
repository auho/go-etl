package runner

import (
	"fmt"

	"github.com/auho/go-etl/v3/task"
	"github.com/auho/go-etl/v3/tool/slicex"
	"github.com/auho/go-toolkit-flow/v3/flow"
	"github.com/auho/go-toolkit-flow/v3/storage"
	"github.com/auho/go-toolkit-flow/v3/storage/database/source"
)

// executor builds the flow options that wire processors/destinations into a flow.
type executor interface {
	flowOptions() ([]flow.Option[map[string]any, map[string]any], error)
	processors() []processor
}

// toProcessors converts a typed slice of processors into a []processor.
func toProcessors[T processor](items []T) []processor {
	ps := make([]processor, len(items))
	for i, item := range items {
		ps[i] = item
	}
	return ps
}

// run builds a Runner from opts, constructs the source over table, and executes
// the flow driven by executor e.
func run(table task.Table, e executor, opts ...RunnerOption) error {
	r := &Runner{}
	r.prepare(opts)

	s, err := r.source(table, e.processors())
	if err != nil {
		return fmt.Errorf("source: %w", err)
	}

	err = r.run(s, e)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}

// Runner drives the source scan and the flow execution.
type Runner struct {
	sourceConfig SourceConfig
}

// prepare applies the given options to the Runner and validates the config.
func (r *Runner) prepare(opts []RunnerOption) {
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		opt(r)
	}

	r.sourceConfig.check()
}

// source builds the paginated source section over table, selecting the union of
// fields required by all processors.
func (r *Runner) source(table task.Table, ps []processor) (*source.Section[storage.MapEntry], error) {
	fields := []string{table.IDName()}
	for _, p := range ps {
		pFields, err := p.Fields()
		if err != nil {
			return nil, fmt.Errorf("fields: %w", err)
		}

		fields = append(fields, pFields...)
	}

	fields = slicex.SliceDropDuplicates(fields)

	ds, err := source.NewSectionMapWithGorm(
		source.SectionConfig{
			Concurrency: r.sourceConfig.Concurrency,
			MaxItems:    r.sourceConfig.Maximum,
			StartID:     0,
			EndID:       0,
			PageSize:    r.sourceConfig.PageSize,
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

// run assembles the flow from the source section and executor options, then runs it.
func (r *Runner) run(s *source.Section[storage.MapEntry], e executor) error {
	eOpts, err := e.flowOptions()
	if err != nil {
		return fmt.Errorf("options: %w", err)
	}

	opts := []flow.Option[map[string]any, map[string]any]{
		flow.WithSource[map[string]any, map[string]any](s),
	}

	opts = append(opts, eOpts...)

	err = flow.RunFlow[map[string]any](opts...)
	if err != nil {
		return fmt.Errorf("RunFlow: %w", err)
	}

	return nil
}
