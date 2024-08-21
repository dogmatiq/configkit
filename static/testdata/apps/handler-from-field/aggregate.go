package app

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

// Aggregate is an aggregate used for testing.
type Aggregate struct{}

// ApplyEvent updates the aggregate instance to reflect the occurrence of an
// event that was recorded against this instance.
func (Aggregate) ApplyEvent(dogma.Event) {}

// AggregateHandler is a test implementation of dogma.AggregateMessageHandler.
type AggregateHandler struct{}

// New returns a new account instance.
func (*AggregateHandler) New() dogma.AggregateRoot {
	return Aggregate{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (*AggregateHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("<aggregate>", "195ede4a-3f26-4d19-a8fe-41b2a5f92d06")

	c.Routes(
		dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
		dogma.RecordsEvent[stubs.CommandStub[stubs.TypeB]](),
	)
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (*AggregateHandler) RouteCommandToInstance(dogma.Command) string {
	return "<aggregate>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (*AggregateHandler) HandleCommand(
	dogma.AggregateRoot,
	dogma.AggregateCommandScope,
	dogma.Command,
) {
}
