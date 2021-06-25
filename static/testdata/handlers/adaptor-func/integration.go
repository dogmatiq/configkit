package app

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// IntegrationHandler is the type that provides the handler logic, but is not
// itself an implementation of dogma.IntegrationMessageHandler.
type IntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("<integration>", "099b5b8d-9e04-422f-bcc3-bb0d451158c7")

	c.ConsumesCommandType(fixtures.MessageA{})
	c.ConsumesCommandType(fixtures.MessageB{})

	c.ProducesEventType(fixtures.MessageC{})
	c.ProducesEventType(fixtures.MessageD{})
}

// AdaptIntegration adapts h to the dogma.IntegrationMessageHandler interface.
func AdaptIntegration(h IntegrationHandler) dogma.IntegrationMessageHandler {
	panic("the implementation of this function is irrelevant to the analyzer")
}
