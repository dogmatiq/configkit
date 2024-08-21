package app

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

// AggregateHandler is the type that provides the handler logic, but is not
// itself an implementation of dogma.AggregateMessageHandler.
type AggregateHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (AggregateHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("<aggregate>", "ef16c9d1-d7b6-4c99-a0e7-a59218e544fc")

	c.Routes(
		dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
		dogma.RecordsEvent[stubs.EventStub[stubs.TypeA]](),
	)
}

// PartialAggregateMessageHandler is the subset of dogma.AggregateMessageHandler
// that must be implemented for a type to be detected as a concrete
// implementation.
type PartialAggregateMessageHandler interface {
	Configure(c dogma.AggregateConfigurer)
}

// AdaptAggregate adapts the argument to the dogma.AggregateMessageHandler interface.
func AdaptAggregate(PartialAggregateMessageHandler) dogma.AggregateMessageHandler {
	panic("the implementation of this function is irrelevant to the analyzer")
}
