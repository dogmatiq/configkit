package app

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// ProjectionHandler is a test implementation of dogma.ProjectionMessageHandler.
type ProjectionHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (ProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("<projection>", "3dfcd7cd-1f63-47a1-9be7-3242bd252423")

	c.Routes(
		dogma.HandlesEvent[fixtures.MessageA](),
		dogma.HandlesEvent[fixtures.MessageB](),
	)
}

// HandleEvent updates the projection to reflect the occurrence of an event.
func (ProjectionHandler) HandleEvent(
	_ context.Context,
	_, _, _ []byte,
	_ dogma.ProjectionEventScope,
	_ dogma.Event,
) (ok bool, err error) {
	return false, nil
}

// ResourceVersion returns the version of the resource r.
func (ProjectionHandler) ResourceVersion(context.Context, []byte) ([]byte, error) {
	return nil, nil
}

// CloseResource informs the projection that the resource r will not be
// used in any future calls to HandleEvent().
func (ProjectionHandler) CloseResource(context.Context, []byte) error {
	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
func (ProjectionHandler) TimeoutHint(dogma.Message) time.Duration {
	return 0
}

// Compact reduces the size of the projection's data.
func (ProjectionHandler) Compact(context.Context, dogma.ProjectionCompactScope) error {
	return nil
}
