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

// NewAggregateHandler returns a new AggregateHandler.
func NewAggregateHandler() AggregateHandler {
	panic("the implementation of this function is irrelevant to the analyzer")
}

// New returns a new account instance.
func (AggregateHandler) New() dogma.AggregateRoot {
	return Aggregate{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (AggregateHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("<aggregate>", "ef16c9d1-d7b6-4c99-a0e7-a59218e544fc")

	c.Routes(
		dogma.HandlesCommand[fixtures.MessageA](),
		dogma.HandlesCommand[fixtures.MessageB](),
		dogma.RecordsEvent[fixtures.MessageC](),
		dogma.RecordsEvent[fixtures.MessageD](),
	)
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
