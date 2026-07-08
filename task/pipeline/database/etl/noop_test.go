package etl

import (
	"testing"

	"github.com/auho/go-etl/v3/task/pipeline"
	"github.com/auho/go-toolkit-flow/v3/exec"
	consitem "github.com/auho/go-toolkit-flow/v3/exec/consumer/item"
	"github.com/auho/go-toolkit-flow/v3/flow"
	"github.com/auho/go-toolkit-flow/v3/storage"
	mocksource "github.com/auho/go-toolkit-flow/v3/storage/mock/source"
)

// TestNoopConsumer runs an in-memory flow: a mock source feeds the Noop
// consumer, which discards every item. This test does not require MySQL and
// validates the ETL pipeline wiring without external infrastructure.
func TestNoopConsumer(t *testing.T) {
	src, err := mocksource.NewMap(mocksource.Config{
		IDName:      "id",
		PageSize:    10,
		Total:       100,
		Concurrency: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	noop := NewNoop()
	runner := consitem.NewRunner[storage.MapEntry, storage.MapEntry](noop)

	err = pipeline.Run(src,
		flow.WithGroup[storage.MapEntry, storage.MapEntry](
			[]exec.Runner[storage.MapEntry, storage.MapEntry]{runner},
		),
	)
	if err != nil {
		t.Error(err)
	}
}
