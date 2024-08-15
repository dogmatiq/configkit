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
	c.Identity("<integration>", "92cce461-8d30-409b-8d5a-406f656cef2d")

	var routes []dogma.IntegrationRoute
	if condition == 0 {
		routes = []dogma.IntegrationRoute{
			dogma.HandlesCommand[fixtures.MessageA](),
		}
		routes = append(
			routes,
			dogma.RecordsEvent[fixtures.MessageC](),
		)
	} else {
		routes = append(
			routes,
			[]dogma.IntegrationRoute{
				dogma.HandlesCommand[fixtures.MessageB](),
				dogma.RecordsEvent[fixtures.MessageD](),
			}...,
		)
	}

	c.Routes(routes...)
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
