package app

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

// ProcessHandler is the type that provides the handler logic, but is not
// itself an implementation of dogma.ProcessMessageHandler.
type ProcessHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (ProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Identity("<process>", "5e839b73-170b-42c0-bf41-8feee4b5a583")

	c.Routes(
		dogma.HandlesEvent[stubs.EventStub[stubs.TypeP]](),
		dogma.ExecutesCommand[stubs.CommandStub[stubs.TypeP]](),
	)
}

// PartialProcessMessageHandler is the subset of dogma.ProcessMessageHandler
// that must be implemented for a type to be detected as a concrete
// implementation.
type PartialProcessMessageHandler interface {
	Configure(c dogma.ProcessConfigurer)
}

// AdaptProcess adapts the argument to the dogma.ProcessMessageHandler interface.
func AdaptProcess(PartialProcessMessageHandler) dogma.ProcessMessageHandler {
	panic("the implementation of this function is irrelevant to the analyzer")
}
