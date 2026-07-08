package pipeline

import (
	"github.com/auho/go-toolkit-flow/v3/flow"
	"github.com/auho/go-toolkit-flow/v3/storage"
)

// Run builds the source and executes the flow.
// Group/dest orchestration is passed by the caller via opts (flow.WithGroup);
// the orchestration layer does not preset any group strategy, giving the
// scenario layer full freedom.
func Run[SE, DE storage.Entry](source storage.Source[SE], opts ...flow.Option[SE, DE]) error {
	opts = append([]flow.Option[SE, DE]{flow.WithSource[SE, DE](source)}, opts...)
	return flow.RunFlow[SE, DE](opts...)
}
