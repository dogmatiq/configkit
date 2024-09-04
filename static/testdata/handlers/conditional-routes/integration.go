package app

import (
	"context"
	"math/rand"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

// IntegrationHandler is a test implementation of
// dogma.IntegrationMessageHandler.
type IntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("<integration>", "92cce461-8d30-409b-8d5a-406f656cef2d")

	if rand.Int() == 0 {
		c.Routes(
			dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
			dogma.RecordsEvent[stubs.EventStub[stubs.TypeA]](),
		)
	} else {
		c.Routes(
			dogma.HandlesCommand[stubs.CommandStub[stubs.TypeB]](),
			dogma.RecordsEvent[stubs.EventStub[stubs.TypeB]](),
		)
	}
}

// HandleCommand handles a command message that has been routed to this handler.
func (IntegrationHandler) HandleCommand(
	context.Context,
	dogma.IntegrationCommandScope,
	dogma.Command,
) error {
	return nil
}
