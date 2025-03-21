package app

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

// ProjectionHandler is the type that provides the handler logic, but is not
// itself an implementation of dogma.ProjectionMessageHandler.
type ProjectionHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (ProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("<projection>", "823e61d3-ace1-469d-b0a6-778e84c0a508")

	c.Routes(
		dogma.HandlesEvent[stubs.EventStub[stubs.TypeR]](),
	)
}

// PartialProjectionMessageHandler is the subset of
// dogma.ProjectionMessageHandler that must be implemented for a type to be
// detected as a concrete implementation.
type PartialProjectionMessageHandler interface {
	Configure(c dogma.ProjectionConfigurer)
}

// AdaptProjection adapts the argument to the dogma.ProjectionMessageHandler interface.
func AdaptProjection(PartialProjectionMessageHandler) dogma.ProjectionMessageHandler {
	panic("the implementation of this function is irrelevant to the analyzer")
}
