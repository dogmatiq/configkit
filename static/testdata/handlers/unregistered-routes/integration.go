package app

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

// IntegrationHandler is a test implementation of
// dogma.IntegrationMessageHandler.
type IntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("<integration>", "ac391765-da58-4e7c-a478-e4725eb2b0e9")

	// Create a route that is never passed to c.Routes().
	dogma.HandlesCommand[stubs.CommandStub[stubs.TypeX]]()

	// Ensure there is still _some_ call to Routes().
	c.Routes(
		dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
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
