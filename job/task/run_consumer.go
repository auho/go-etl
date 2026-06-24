package task

import (
	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-toolkit-flow/exec"
	"github.com/auho/go-toolkit-flow/exec/consumer/item"
	"github.com/auho/go-toolkit-flow/flow"
	"github.com/auho/go-toolkit-flow/processor/consumer"
)

var _ executor = (*consumerExecutor)(nil)

type consumerExecutor struct {
	items []consumer.Item[map[string]any]
}

func (c *consumerExecutor) exec() []flow.Option[map[string]any, map[string]any] {
	var opts []flow.Option[map[string]any, map[string]any]

	for _, _item := range c.items {
		opts = append(opts, flow.WithGroup(
			[]exec.Runner[map[string]any, map[string]any]{
				item.NewRunner[map[string]any, map[string]any](_item),
			},
		))
	}

	return opts
}

func RunConsumer(js job.Table, tasks []itemConsumer, opts ...ConfigOption) {
	var ps []processor
	var items []consumer.Item[map[string]any]
	for _, t := range tasks {
		ps = append(ps, t)
		items = append(items, t)
	}

	e := &consumerExecutor{items: items}

	run(js, ps, e, opts...)
}
