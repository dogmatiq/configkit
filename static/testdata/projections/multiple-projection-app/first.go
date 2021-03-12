package app

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// FirstProjection is a projection used for testing.
type FirstProjection struct{}

// FirstProjectionHandler is a test implementation of
// dogma.ProjectionMessageHandler.
type FirstProjectionHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (FirstProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("<first-projection>", "9174783f-4f12-4619-b5c6-c4ab70bd0937")

	c.ConsumesEventType(fixtures.MessageA{})
}

// HandleEvent updates the projection to reflect the occurrence of an event.
func (FirstProjectionHandler) HandleEvent(
	ctx context.Context,
	r, c, n []byte,
	s dogma.ProjectionEventScope,
	m dogma.Message,
) (ok bool, err error) {
	return false, nil
}

// ResourceVersion returns the version of the resource r.
func (FirstProjectionHandler) ResourceVersion(
	ctx context.Context,
	r []byte,
) ([]byte, error) {
	return nil, nil
}

// CloseResource informs the projection that the resource r will not be
// used in any future calls to HandleEvent().
func (FirstProjectionHandler) CloseResource(ctx context.Context, r []byte) error {
	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
func (FirstProjectionHandler) TimeoutHint(m dogma.Message) time.Duration {
	return 0
}

// Compact reduces the size of the projection's data.
func (FirstProjectionHandler) Compact(ctx context.Context, s dogma.ProjectionCompactScope) error {
	return nil
}
