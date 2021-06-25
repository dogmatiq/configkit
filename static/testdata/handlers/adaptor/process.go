package app

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// ProcessHandler is the type that provides the handler logic, but is not
// itself an implementation of dogma.ProcessMessageHandler.
type ProcessHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (ProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Identity("<process>", "5e839b73-170b-42c0-bf41-8feee4b5a583")

	c.ConsumesEventType(fixtures.MessageA{})
	c.ConsumesEventType(fixtures.MessageB{})

	c.ProducesCommandType(fixtures.MessageC{})
	c.ProducesCommandType(fixtures.MessageD{})

	c.SchedulesTimeoutType(fixtures.MessageE{})
	c.SchedulesTimeoutType(fixtures.MessageF{})
}

// AdaptProcess adapts h to the dogma.ProcessMessageHandler interface.
func AdaptProcess(h ProcessHandler) dogma.ProcessMessageHandler {
	panic("the implementation of this function is irrelevant to the analyzer")
}
