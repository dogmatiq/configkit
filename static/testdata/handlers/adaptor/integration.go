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

	c.Routes(
		dogma.HandlesCommand[fixtures.MessageA](),
		dogma.HandlesCommand[fixtures.MessageB](),
		dogma.RecordsEvent[fixtures.MessageC](),
		dogma.RecordsEvent[fixtures.MessageD](),
	)
}

// PartialIntegrationMessageHandler is the subset of
// dogma.IntegrationMessageHandler that must be implemented for a type to be
// detected as a concrete implementation.
type PartialIntegrationMessageHandler interface {
	Configure(c dogma.IntegrationConfigurer)
}

// AdaptIntegration adapts h to the dogma.IntegrationMessageHandler interface.
func AdaptIntegration(h PartialIntegrationMessageHandler) dogma.IntegrationMessageHandler {
	panic("the implementation of this function is irrelevant to the analyzer")
}
