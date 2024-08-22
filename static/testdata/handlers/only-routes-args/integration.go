package app

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// IntegrationHandler is a test implementation of
// dogma.IntegrationMessageHandler.
type IntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("<integration>", "ac391765-da58-4e7c-a478-e4725eb2b0e9")

	c.Routes(
		dogma.HandlesCommand[fixtures.MessageA](),
		dogma.HandlesCommand[fixtures.MessageB](),
		dogma.RecordsEvent[fixtures.MessageC](),
		dogma.RecordsEvent[fixtures.MessageD](),
	)

	// These assignments should not be included into handler configuration
	// as they are not arguments to the Routes() method.
	var r1 dogma.IntegrationRoute = dogma.RecordsEvent[fixtures.MessageE]()
	var r2 dogma.IntegrationRoute = dogma.HandlesCommand[fixtures.MessageF]()
	// Assign to blank identifiers to avoid unused variable errors during
	// compilation.
	_ = r1
	_ = r2
}

// HandleCommand handles a command message that has been routed to this handler.
func (IntegrationHandler) HandleCommand(
	context.Context,
	dogma.IntegrationCommandScope,
	dogma.Command,
) error {
	return nil
}
