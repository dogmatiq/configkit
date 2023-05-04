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
	c.Identity("<process>", "f754da79-205b-4d65-889f-0d8ae86e3def")

	var routes []dogma.ProcessRoute
	if condition == 0 {
		routes = []dogma.ProcessRoute{
			dogma.HandlesEvent[fixtures.MessageA](),
			dogma.HandlesEvent[fixtures.MessageB](),
		}
		routes = append(
			routes,
			dogma.ExecutesCommand[fixtures.MessageC](),
		)
	} else {
		routes = append(
			routes,
			[]dogma.ProcessRoute{
				dogma.ExecutesCommand[fixtures.MessageD](),
				dogma.SchedulesTimeout[fixtures.MessageE](),
				dogma.SchedulesTimeout[fixtures.MessageF](),
			}...,
		)
	}

	c.Routes(routes...)
}

// RouteEventToInstance returns the ID of the process instance that is
// targeted by m.
func (ProcessHandler) RouteEventToInstance(
	ctx context.Context,
	m dogma.Message,
) (string, bool, error) {
	return "<process>", true, nil
}

// HandleEvent handles an event message.
func (ProcessHandler) HandleEvent(
	ctx context.Context,
	r dogma.ProcessRoot,
	s dogma.ProcessEventScope,
	m dogma.Message,
) error {
	return nil
}

// HandleTimeout handles a timeout message that has been scheduled with
// ProcessScope.ScheduleTimeout().
func (ProcessHandler) HandleTimeout(
	ctx context.Context,
	r dogma.ProcessRoot,
	s dogma.ProcessTimeoutScope,
	m dogma.Message,
) error {
	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
func (ProcessHandler) TimeoutHint(m dogma.Message) time.Duration {
	return 0
}
