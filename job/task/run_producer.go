package task

import (
	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-toolkit-flow/exec"
	"github.com/auho/go-toolkit-flow/exec/producer/item"
	"github.com/auho/go-toolkit-flow/flow"
	"github.com/auho/go-toolkit-flow/processor/producer"
)

var _ executor = (*producerExecutor)(nil)

type producerExecutor struct {
	items []producer.Item[map[string]any, map[string]any]
}

func (c *producerExecutor) options() []flow.Option[map[string]any, map[string]any] {
	var opts []flow.Option[map[string]any, map[string]any]

	for _, a := range c.items {
		opts = append(opts, flow.WithGroup(
			[]exec.Runner[map[string]any, map[string]any]{
				item.NewRunner[map[string]any, map[string]any](a),
			},
		))
	}

	return opts
}

func RunProducer(js job.Table, tasks []itemProducer, opts ...ConfigOption) {
	var ps []processor
	var items []producer.Item[map[string]any, map[string]any]
	for _, t := range tasks {
		ps = append(ps, t)
		items = append(items, t)
	}

	e := &producerExecutor{items: items}

	run(js, ps, e, opts...)
}
