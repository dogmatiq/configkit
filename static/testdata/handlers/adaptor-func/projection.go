package app

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// ProjectionHandler is the type that provides the handler logic, but is not
// itself an implementation of dogma.ProjectionMessageHandler.
type ProjectionHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (ProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("<projection>", "823e61d3-ace1-469d-b0a6-778e84c0a508")

	c.ConsumesEventType(fixtures.MessageA{})
	c.ConsumesEventType(fixtures.MessageB{})
}

// AdaptProjection adapts h to the dogma.ProjectionMessageHandler interface.
func AdaptProjection(h ProjectionHandler) dogma.ProjectionMessageHandler {
	panic("the implementation of this function is irrelevant to the analyzer")
}
