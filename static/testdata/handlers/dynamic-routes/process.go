package app

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
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
	c.Identity("<process>", "d00131d2-f99c-44e3-be11-104e69b04f77")

	routes := []dogma.ProcessRoute{
		dogma.HandlesEvent[fixtures.MessageA](),
		dogma.HandlesEvent[fixtures.MessageB](),
		dogma.ExecutesCommand[fixtures.MessageC](),
		dogma.ExecutesCommand[fixtures.MessageD](),
		dogma.SchedulesTimeout[fixtures.MessageE](),
		dogma.SchedulesTimeout[fixtures.MessageF](),
	}

	c.Routes(routes...)
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

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
func (ProcessHandler) TimeoutHint(dogma.Message) time.Duration {
	return 0
}
