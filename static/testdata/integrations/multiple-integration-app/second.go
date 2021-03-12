package app

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// SecondIntegration is an integration used for testing.
type SecondIntegration struct{}

// ApplyEvent updates the integration instance to reflect the occurrence of an
// event that was recorded against this instance.
func (SecondIntegration) ApplyEvent(m dogma.Message) {}

// SecondIntegrationHandler is a test implementation of
// dogma.IntegrationMessageHandler.
type SecondIntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (SecondIntegrationHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("<second-integration>", "6bed3fbc-30e2-44c7-9a5b-e440ffe370d9")

	c.ConsumesCommandType(fixtures.MessageA{})

	c.ProducesEventType(fixtures.MessageB{})
}

// RouteCommandToInstance returns the ID of the integration instance that is
// targetted by m.
func (SecondIntegrationHandler) RouteCommandToInstance(m dogma.Message) string {
	return "<second-integration>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (SecondIntegrationHandler) HandleCommand(
	ctx context.Context,
	s dogma.IntegrationCommandScope,
	m dogma.Message,
) error {
	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
func (SecondIntegrationHandler) TimeoutHint(m dogma.Message) time.Duration {
	return 0
}
