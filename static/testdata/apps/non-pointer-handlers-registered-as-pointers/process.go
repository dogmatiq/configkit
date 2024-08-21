package app

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

// Process is a process used for testing.
type Process struct{}

// ProcessHandler is a test implementation of dogma.ProcessMessageHandler.
type ProcessHandler struct{}

// New constructs a new process instance initialized with any default values and
// returns the process root.
func (ProcessHandler) New() dogma.ProcessRoot {
	return Process{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (ProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Identity("<process>", "39af6b34-5fa1-4f3a-b049-40a5e1d9b33b")

	c.Routes(
		dogma.HandlesEvent[stubs.EventStub[stubs.TypeA]](),
		dogma.ExecutesCommand[stubs.CommandStub[stubs.TypeB]](),
	)
}

// RouteEventToInstance returns the ID of the process instance that is
// targeted by m.
func (ProcessHandler) RouteEventToInstance(
	context.Context,
	dogma.Event,
) (string, bool, error) {
	return "<process>", true, nil
}

// HandleEvent handles an event message.
func (ProcessHandler) HandleEvent(
	context.Context,
	dogma.ProcessRoot,
	dogma.ProcessEventScope,
	dogma.Event,
) error {
	return nil
}

// HandleTimeout handles a timeout message that has been scheduled with
// ProcessScope.ScheduleTimeout().
func (ProcessHandler) HandleTimeout(
	context.Context,
	dogma.ProcessRoot,
	dogma.ProcessTimeoutScope,
	dogma.Timeout,
) error {
	return nil
}
