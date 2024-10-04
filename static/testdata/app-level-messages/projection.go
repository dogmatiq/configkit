package app

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
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
		dogma.HandlesEvent[stubs.EventStub[stubs.TypeA]](),
		dogma.HandlesEvent[stubs.EventStub[stubs.TypeC]](),
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

// Compact reduces the size of the projection's data.
func (ProjectionHandler) Compact(context.Context, dogma.ProjectionCompactScope) error {
	return nil
}
