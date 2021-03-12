package app

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// SecondProjection is a projection used for testing.
type SecondProjection struct{}

// SecondProjectionHandler is a test implementation of
// dogma.ProjectionMessageHandler.
type SecondProjectionHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (SecondProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("<second-projection>", "2e22850e-7c84-4b3f-b8b3-25ac743d90f2")

	c.ConsumesEventType(fixtures.MessageB{})
}

// HandleEvent updates the projection to reflect the occurrence of an event.
func (SecondProjectionHandler) HandleEvent(
	ctx context.Context,
	r, c, n []byte,
	s dogma.ProjectionEventScope,
	m dogma.Message,
) (ok bool, err error) {
	return false, nil
}

// ResourceVersion returns the version of the resource r.
func (SecondProjectionHandler) ResourceVersion(
	ctx context.Context,
	r []byte,
) ([]byte, error) {
	return nil, nil
}

// CloseResource informs the projection that the resource r will not be
// used in any future calls to HandleEvent().
func (SecondProjectionHandler) CloseResource(ctx context.Context, r []byte) error {
	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
func (SecondProjectionHandler) TimeoutHint(m dogma.Message) time.Duration {
	return 0
}

// Compact reduces the size of the projection's data.
func (SecondProjectionHandler) Compact(ctx context.Context, s dogma.ProjectionCompactScope) error {
	return nil
}
