package app

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// AggregateHandler is the type that provides the handler logic, but is not
// itself an implementation of dogma.AggregateMessageHandler.
type AggregateHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (AggregateHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("<aggregate>", "ef16c9d1-d7b6-4c99-a0e7-a59218e544fc")

	c.ConsumesCommandType(fixtures.MessageA{})
	c.ConsumesCommandType(fixtures.MessageB{})

	c.ProducesEventType(fixtures.MessageC{})
	c.ProducesEventType(fixtures.MessageD{})
}

// AdaptAggregate adapts h to the dogma.AggregateMessageHandler interface.
func AdaptAggregate(h AggregateHandler) dogma.AggregateMessageHandler {
	panic("the implementation of this function is irrelevant to the analyzer")
}
