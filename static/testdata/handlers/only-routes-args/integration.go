package app

import (
	"context"
	"time"

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

// RouteCommandToInstance returns the ID of the integration instance that is
// targetted by m.
func (IntegrationHandler) RouteCommandToInstance(dogma.Command) string {
	return "<integration>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (IntegrationHandler) HandleCommand(
	context.Context,
	dogma.IntegrationCommandScope,
	dogma.Command,
) error {
	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
func (IntegrationHandler) TimeoutHint(dogma.Message) time.Duration {
	return 0
}
