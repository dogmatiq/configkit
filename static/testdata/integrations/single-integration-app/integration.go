package app

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// Integration is an integration used for testing.
type Integration struct{}

// ApplyEvent updates the integration instance to reflect the occurrence of an
// event that was recorded against this instance.
func (Integration) ApplyEvent(m dogma.Message) {}

// IntegrationHandler is a test implementation of
// dogma.IntegrationMessageHandler.
type IntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("<integration>", "ef16c9d1-d7b6-4c99-a0e7-a59218e544fc")

	c.ConsumesCommandType(fixtures.MessageA{})
	c.ConsumesCommandType(fixtures.MessageB{})

	c.ProducesEventType(fixtures.MessageC{})
	c.ProducesEventType(fixtures.MessageD{})
}

// RouteCommandToInstance returns the ID of the integration instance that is
// targetted by m.
func (IntegrationHandler) RouteCommandToInstance(m dogma.Message) string {
	return "<integration>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (IntegrationHandler) HandleCommand(
	ctx context.Context,
	s dogma.IntegrationCommandScope,
	m dogma.Message,
) error {
	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
func (IntegrationHandler) TimeoutHint(m dogma.Message) time.Duration {
	return 0
}
