package app

import (
	"github.com/dogmatiq/dogma"
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

// Configure configures the behavior of the engine as it relates to this
// handler.
func (AggregateHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("<aggregate>", "df648962-7d96-427e-8bc2-5a4efdb4cc4b")

	c.Routes(nil)
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
