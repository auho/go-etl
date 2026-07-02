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

func (c *producerExecutor) options() ([]flow.Option[map[string]any, map[string]any], error) {
	var opts []flow.Option[map[string]any, map[string]any]

	for _, _p := range c.producers {
		dests, err := _p.destinations()
		if err != nil {
			return nil, fmt.Errorf("producer destinations: %w", err)
		}

		opts = append(opts, flow.WithGroup(
			[]exec.Runner[map[string]any, map[string]any]{
				item.NewRunner[map[string]any, map[string]any](_p),
			},
			dests...,
		))
	}

	return opts, nil
}

func RunProducer(js job.Table, producers []itemProducer, opts ...ConfigOption) {
	var ps []processor
	for _, _p := range producers {
		ps = append(ps, _p)
	}

	e := &producerExecutor{producers: producers}

	run(js, ps, e, opts...)
}
