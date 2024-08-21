package app

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

// Integration is an integration used for testing.
type Integration struct{}

// ApplyEvent updates the integration instance to reflect the occurrence of an
// event that was recorded against this instance.
func (Integration) ApplyEvent(dogma.Event) {}

// IntegrationHandler is a test implementation of
// dogma.IntegrationMessageHandler.
type IntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("<integration>", "099b5b8d-9e04-422f-bcc3-bb0d451158c7")

	c.Routes(
		dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
		dogma.RecordsEvent[stubs.EventStub[stubs.TypeB]](),
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
