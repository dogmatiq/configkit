package app

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// Aggregate is an aggregate used for testing.
type Aggregate struct{}

// ApplyEvent updates the aggregate instance to reflect the occurrence of an
// event that was recorded against this instance.
func (Aggregate) ApplyEvent(m dogma.Message) {}

// AggregateHandler is a test implementation of dogma.AggregateMessageHandler.
type AggregateHandler struct{}

// New returns a new account instance.
func (AggregateHandler) New() dogma.AggregateRoot {
	return Aggregate{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (AggregateHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("<aggregate>", "ee1814e3-194d-438a-916e-ee7766598646")

	c.Routes(
		dogma.HandlesCommand[fixtures.MessageA](),
		dogma.HandlesCommand[fixtures.MessageB](),
		dogma.RecordsEvent[fixtures.MessageC](),
		dogma.RecordsEvent[fixtures.MessageD](),
	)
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (AggregateHandler) RouteCommandToInstance(m dogma.Message) string {
	return "<aggregate>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (AggregateHandler) HandleCommand(
	r dogma.AggregateRoot,
	s dogma.AggregateCommandScope,
	m dogma.Message,
) {
}
