package app

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// Projection is a projection used for testing.
type Projection struct{}

// ProjectionHandler is a test implementation of dogma.ProjectionMessageHandler.
type ProjectionHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (ProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("<projection>", "7a5090e0-7248-4a58-8d70-a5dfd8c8abe1")

	c.Routes(
		dogma.HandlesEvent[fixtures.MessageB](),
		dogma.HandlesEvent[fixtures.MessageC](),
	)
}

// HandleEvent updates the projection to reflect the occurrence of an event.
func (ProjectionHandler) HandleEvent(
	ctx context.Context,
	r, c, n []byte,
	s dogma.ProjectionEventScope,
	m dogma.Message,
) (ok bool, err error) {
	return false, nil
}

// ResourceVersion returns the version of the resource r.
func (ProjectionHandler) ResourceVersion(
	ctx context.Context,
	r []byte,
) ([]byte, error) {
	return nil, nil
}

// CloseResource informs the projection that the resource r will not be
// used in any future calls to HandleEvent().
func (ProjectionHandler) CloseResource(ctx context.Context, r []byte) error {
	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
func (ProjectionHandler) TimeoutHint(m dogma.Message) time.Duration {
	return 0
}

// Compact reduces the size of the projection's data.
func (ProjectionHandler) Compact(ctx context.Context, s dogma.ProjectionCompactScope) error {
	return nil
}
