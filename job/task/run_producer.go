package task

import (
	"fmt"

	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-toolkit-flow/v3/exec"
	"github.com/auho/go-toolkit-flow/v3/exec/producer/item"
	"github.com/auho/go-toolkit-flow/v3/flow"
)

var _ executor = (*producerExecutor)(nil)

// producerExecutor is the executor for producer tasks: it wires each producer
// with its destinations into the flow.
type producerExecutor struct {
	producers []itemProducer
}

func (p *producerExecutor) flowOptions() ([]flow.Option[map[string]any, map[string]any], error) {
	var opts []flow.Option[map[string]any, map[string]any]

	for _, producer := range p.producers {
		dests, err := producer.destinations()
		if err != nil {
			return nil, fmt.Errorf("producer destinations: %w", err)
		}

		opts = append(opts, flow.WithGroup(
			[]exec.Runner[map[string]any, map[string]any]{
				item.NewRunner[map[string]any, map[string]any](producer),
			},
			dests...,
		))
	}

	return opts, nil
}

func (p *producerExecutor) processors() []processor {
	return toProcessors(p.producers)
}

// RunProducer runs producer tasks over table with the given runner options.
func RunProducer(table job.Table, producers []itemProducer, opts ...RunnerOption) error {
	return run(table, &producerExecutor{producers: producers}, opts...)
}
