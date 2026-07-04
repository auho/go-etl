package task

import (
	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-toolkit-flow/v3/exec"
	"github.com/auho/go-toolkit-flow/v3/exec/consumer/item"
	"github.com/auho/go-toolkit-flow/v3/flow"
)

var _ executor = (*consumerExecutor)(nil)

// consumerExecutor is the executor for consumer tasks: it wires each consumer
// (no destinations) into the flow.
type consumerExecutor struct {
	consumers []itemConsumer
}

func (c *consumerExecutor) flowOptions() ([]flow.Option[map[string]any, map[string]any], error) {
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

func (c *consumerExecutor) processors() []processor {
	return toProcessors(c.consumers)
}

// RunConsumer runs consumer tasks over table with the given runner options.
func RunConsumer(table job.Table, consumers []itemConsumer, opts ...RunnerOption) error {
	return run(table, &consumerExecutor{consumers: consumers}, opts...)
}
