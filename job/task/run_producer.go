package task

import (
	"fmt"

	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-toolkit-flow/v3/exec"
	"github.com/auho/go-toolkit-flow/v3/exec/producer/item"
	"github.com/auho/go-toolkit-flow/v3/flow"
)

var _ executor = (*producerExecutor)(nil)

type producerExecutor struct {
	producers []itemProducer
}

func (p *producerExecutor) options() ([]flow.Option[map[string]any, map[string]any], error) {
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

func RunProducer(table job.Table, producers []itemProducer, opts ...RunnerOption) error {
	return run(table, toProcessors(producers), &producerExecutor{producers: producers}, opts...)
}
