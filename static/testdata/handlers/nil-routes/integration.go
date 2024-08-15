package app

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
)

// IntegrationHandler is a test implementation of
// dogma.IntegrationMessageHandler.
type IntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("<integration>", "363039e5-2938-4b2c-9bec-dcb29dee2da1")

	c.Routes(nil)
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
