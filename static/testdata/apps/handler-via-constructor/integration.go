package app

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

// IntegrationHandler is a test implementation of
// dogma.IntegrationMessageHandler.
type IntegrationHandler struct{}

// NewIntegrationHandler returns a new IntegrationHandler.
func NewIntegrationHandler() IntegrationHandler {
	panic("the implementation of this function is irrelevant to the analyzer")
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("<integration>", "099b5b8d-9e04-422f-bcc3-bb0d451158c7")

	c.Routes(
		dogma.HandlesCommand[stubs.CommandStub[stubs.TypeB]](),
	)
}

// HandleCommand handles a command message that has been routed to this handler.
func (IntegrationHandler) HandleCommand(
	context.Context,
	dogma.IntegrationCommandScope,
	dogma.Command,
) error {
	return nil
}
