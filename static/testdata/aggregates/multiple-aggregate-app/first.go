package app

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// FirstAggregate is an aggregate used for testing.
type FirstAggregate struct{}

// ApplyEvent updates the aggregate instance to reflect the occurrence of an
// event that was recorded against this instance.
func (FirstAggregate) ApplyEvent(m dogma.Message) {}

// FirstAggregateHandler is a test implementation of dogma.AggregateMessageHandler.
type FirstAggregateHandler struct{}

// New returns a new account instance.
func (FirstAggregateHandler) New() dogma.AggregateRoot {
	return FirstAggregate{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (FirstAggregateHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("<first-aggregate>", "e6300d8d-6530-405e-9729-e9ca21df23d3")

	c.ConsumesCommandType(fixtures.MessageA{})
	c.ConsumesCommandType(fixtures.MessageB{})
	c.ConsumesCommandType(fixtures.MessageC{})

	c.ProducesEventType(fixtures.MessageD{})
	c.ProducesEventType(fixtures.MessageE{})
	c.ProducesEventType(fixtures.MessageF{})
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (FirstAggregateHandler) RouteCommandToInstance(m dogma.Message) string {
	return "<first-aggregate>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (FirstAggregateHandler) HandleCommand(
	r dogma.AggregateRoot,
	s dogma.AggregateCommandScope,
	m dogma.Message,
) {
}
