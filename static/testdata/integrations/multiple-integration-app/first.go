package app

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// FirstIntegration is an integration used for testing.
type FirstIntegration struct{}

// ApplyEvent updates the integration instance to reflect the occurrence of an
// event that was recorded against this instance.
func (FirstIntegration) ApplyEvent(m dogma.Message) {}

// FirstIntegrationHandler is a test implementation of
// dogma.IntegrationMessageHandler.
type FirstIntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (FirstIntegrationHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("<first-integration>", "14cf2812-eead-43b3-9c9c-10db5b469e94")

	c.ConsumesCommandType(fixtures.MessageA{})

	c.ProducesEventType(fixtures.MessageB{})
}

// RouteCommandToInstance returns the ID of the integration instance that is
// targetted by m.
func (FirstIntegrationHandler) RouteCommandToInstance(m dogma.Message) string {
	return "<first-integration>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (FirstIntegrationHandler) HandleCommand(
	ctx context.Context,
	s dogma.IntegrationCommandScope,
	m dogma.Message,
) error {
	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
func (FirstIntegrationHandler) TimeoutHint(m dogma.Message) time.Duration {
	return 0
}
