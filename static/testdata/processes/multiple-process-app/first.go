package app

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// FirstProcess is a process used for testing.
type FirstProcess struct{}

// FirstProcessHandler is a test implementation of dogma.ProcessMessageHandler.
type FirstProcessHandler struct{}

// New constructs a new process instance initialized with any default values and
// returns the process root.
func (FirstProcessHandler) New() dogma.ProcessRoot {
	return FirstProcess{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (FirstProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Identity("<first-process>", "d33198e0-f1f7-4c2d-8ac2-98f68a44414e")

	c.ConsumesEventType(fixtures.MessageA{})

	c.ProducesCommandType(fixtures.MessageB{})

	c.SchedulesTimeoutType(fixtures.MessageC{})
}

// RouteEventToInstance returns the ID of the process instance that is
// targeted by m.
func (FirstProcessHandler) RouteEventToInstance(
	ctx context.Context,
	m dogma.Message,
) (string, bool, error) {
	return "<first-process>", true, nil
}

// HandleEvent handles an event message.
func (FirstProcessHandler) HandleEvent(
	ctx context.Context,
	r dogma.ProcessRoot,
	s dogma.ProcessEventScope,
	m dogma.Message,
) error {
	return nil
}

// HandleTimeout handles a timeout message that has been scheduled with
// ProcessScope.ScheduleTimeout().
func (FirstProcessHandler) HandleTimeout(
	ctx context.Context,
	r dogma.ProcessRoot,
	s dogma.ProcessTimeoutScope,
	m dogma.Message,
) error {
	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
func (FirstProcessHandler) TimeoutHint(m dogma.Message) time.Duration {
	return 0
}
