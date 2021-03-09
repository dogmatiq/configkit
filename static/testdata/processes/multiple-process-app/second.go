package app

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// SecondProcess is a process used for testing.
type SecondProcess struct{}

// SecondProcessHandler is a test implementation of dogma.ProcessMessageHandler.
type SecondProcessHandler struct{}

// New constructs a new process instance initialized with any default values and
// returns the process root.
func (SecondProcessHandler) New() dogma.ProcessRoot {
	return SecondProcess{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (SecondProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Identity("<second-process>", "0311717c-cf51-4292-8ed9-95125302a18e")

	c.ConsumesEventType(fixtures.MessageD{})

	c.ProducesCommandType(fixtures.MessageE{})

	c.SchedulesTimeoutType(fixtures.MessageF{})
}

// RouteEventToInstance returns the ID of the process instance that is
// targeted by m.
func (SecondProcessHandler) RouteEventToInstance(
	ctx context.Context,
	m dogma.Message,
) (string, bool, error) {
	return "<second-process>", true, nil
}

// HandleEvent handles an event message.
func (SecondProcessHandler) HandleEvent(
	ctx context.Context,
	r dogma.ProcessRoot,
	s dogma.ProcessEventScope,
	m dogma.Message,
) error {
	return nil
}

// HandleTimeout handles a timeout message that has been scheduled with
// ProcessScope.ScheduleTimeout().
func (SecondProcessHandler) HandleTimeout(
	ctx context.Context,
	r dogma.ProcessRoot,
	s dogma.ProcessTimeoutScope,
	m dogma.Message,
) error {
	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
func (SecondProcessHandler) TimeoutHint(m dogma.Message) time.Duration {
	return 0
}
