package app

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// Aggregate is an aggregate used for testing.
type Aggregate struct{}

// ApplyEvent updates the aggregate instance to reflect the occurrence of an
// event that was recorded against this instance.
func (Aggregate) ApplyEvent(dogma.Event) {}

// AggregateHandler is a test implementation of dogma.AggregateMessageHandler.
type AggregateHandler struct{}

// New returns a new account instance.
func (AggregateHandler) New() dogma.AggregateRoot {
	return Aggregate{}
}

var condition int

// Configure configures the behavior of the engine as it relates to this
// handler.
func (AggregateHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("<aggregate>", "c3b4b3c7-fbe6-4789-9358-e4f45b154d31")

	var routes []dogma.AggregateRoute
	if condition == 0 {
		routes = []dogma.AggregateRoute{
			dogma.HandlesCommand[fixtures.MessageA](),
		}
		routes = append(
			routes,
			dogma.RecordsEvent[fixtures.MessageC](),
		)
	} else {
		routes = append(
			routes,
			[]dogma.AggregateRoute{
				dogma.HandlesCommand[fixtures.MessageB](),
				dogma.RecordsEvent[fixtures.MessageD](),
			}...,
		)
	}

	c.Routes(routes...)
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (AggregateHandler) RouteCommandToInstance(dogma.Command) string {
	return "<aggregate>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (AggregateHandler) HandleCommand(
	dogma.AggregateRoot,
	dogma.AggregateCommandScope,
	dogma.Command,
) {
}
