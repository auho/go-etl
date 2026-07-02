package task

import (
	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-toolkit-flow/v3/exec"
	"github.com/auho/go-toolkit-flow/v3/exec/consumer/item"
	"github.com/auho/go-toolkit-flow/v3/flow"
)

var _ executor = (*consumerExecutor)(nil)

type consumerExecutor struct {
	consumers []itemConsumer
}

func (c *consumerExecutor) options() ([]flow.Option[map[string]any, map[string]any], error) {
	var opts []flow.Option[map[string]any, map[string]any]

	for _, consumer := range c.consumers {
		opts = append(opts, flow.WithGroup(
			[]exec.Runner[map[string]any, map[string]any]{
				item.NewRunner[map[string]any, map[string]any](consumer),
			},
		))
	}

	return opts, nil
}

func RunConsumer(table job.Table, consumers []itemConsumer, opts ...ConfigOption) error {
	return run(table, toProcessors(consumers), &consumerExecutor{consumers: consumers}, opts...)
}
