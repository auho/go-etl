package etl

import (
	"fmt"

	"github.com/auho/go-etl/v3/task/pipeline"
	"github.com/auho/go-etl/v3/tool/slicex"
	"github.com/auho/go-toolkit-flow/v3/exec"
	consitem "github.com/auho/go-toolkit-flow/v3/exec/consumer/item"
	proditem "github.com/auho/go-toolkit-flow/v3/exec/producer/item"
	"github.com/auho/go-toolkit-flow/v3/flow"
	"github.com/auho/go-toolkit-flow/v3/processor/consumer"
	"github.com/auho/go-toolkit-flow/v3/processor/producer"
	"github.com/auho/go-toolkit-flow/v3/storage"
)

// fielder is the subset of a processor needed for source field collection.
type fielder interface {
	Fields() ([]string, error)
}

// collectFields builds the deduplicated source field list from the id column and
// each processor's Fields.
func collectFields(idName string, ps ...fielder) ([]string, error) {
	fields := []string{idName}
	for _, p := range ps {
		pf, err := p.Fields()
		if err != nil {
			return nil, fmt.Errorf("fields: %w", err)
		}
		fields = append(fields, pf...)
	}
	return slicex.SliceDropDuplicates(fields), nil
}

// defaultSourceConfig returns a SourceConfig with defaults applied.
func defaultSourceConfig() SourceConfig {
	sc := SourceConfig{}
	sc.Check()
	return sc
}

// producerItem is the contract a producer processor satisfies for runProducer.
type producerItem interface {
	fielder
	producer.Item[storage.MapEntry, storage.MapEntry]
	BuildDestinations() ([]storage.Destination[storage.MapEntry], error)
}

// consumerItem is the contract a consumer processor satisfies for runConsumer.
type consumerItem interface {
	fielder
	consumer.Item[storage.MapEntry]
}

// runProducer builds the source, wraps the producer into a runner with its
// destinations, and runs the flow.
func runProducer(source Table, p producerItem) error {
	fields, err := collectFields(source.IDName(), p)
	if err != nil {
		return err
	}

	src, err := newSource(source, fields, defaultSourceConfig())
	if err != nil {
		return err
	}

	dests, err := p.BuildDestinations()
	if err != nil {
		return err
	}

	runner := proditem.NewRunner[storage.MapEntry, storage.MapEntry](p)
	return pipeline.Run(src,
		flow.WithGroup[storage.MapEntry, storage.MapEntry](
			[]exec.Runner[storage.MapEntry, storage.MapEntry]{runner},
			dests...,
		),
	)
}

// runConsumer builds the source, wraps the consumer into a runner (no external
// destinations), and runs the flow.
func runConsumer(source Table, c consumerItem) error {
	fields, err := collectFields(source.IDName(), c)
	if err != nil {
		return err
	}

	src, err := newSource(source, fields, defaultSourceConfig())
	if err != nil {
		return err
	}

	runner := consitem.NewRunner[storage.MapEntry, storage.MapEntry](c)
	return pipeline.Run(src,
		flow.WithGroup[storage.MapEntry, storage.MapEntry](
			[]exec.Runner[storage.MapEntry, storage.MapEntry]{runner},
		),
	)
}
