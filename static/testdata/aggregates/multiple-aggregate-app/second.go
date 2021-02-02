package app

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// SecondAggregate is an aggregate used for testing.
type SecondAggregate struct{}

// ApplyEvent updates the aggregate instance to reflect the occurrence of an
// event that was recorded against this instance.
func (SecondAggregate) ApplyEvent(m dogma.Message) {}

// SecondAggregateHandler is a test implementation of
// dogma.AggregateMessageHandler.
type SecondAggregateHandler struct{}

// New returns a new account instance.
func (SecondAggregateHandler) New() dogma.AggregateRoot {
	return SecondAggregate{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (SecondAggregateHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("<second-aggregate>", "feeb96d0-c56b-4e58-9cd0-d393683c2ec7")

	c.ConsumesCommandType(fixtures.MessageC{})

	c.ProducesEventType(fixtures.MessageD{})
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (SecondAggregateHandler) RouteCommandToInstance(m dogma.Message) string {
	return "<second-aggregate>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (SecondAggregateHandler) HandleCommand(
	r dogma.AggregateRoot,
	s dogma.AggregateCommandScope,
	m dogma.Message,
) {
}
