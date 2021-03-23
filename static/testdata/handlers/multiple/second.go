package app

import (
	"context"
	"time"

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

// SecondProjection is a projection used for testing.
type SecondProjection struct{}

// SecondProjectionHandler is a test implementation of
// dogma.ProjectionMessageHandler.
type SecondProjectionHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (SecondProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("<second-projection>", "2e22850e-7c84-4b3f-b8b3-25ac743d90f2")

	c.ConsumesEventType(fixtures.MessageB{})
}

// HandleEvent updates the projection to reflect the occurrence of an event.
func (SecondProjectionHandler) HandleEvent(
	ctx context.Context,
	r, c, n []byte,
	s dogma.ProjectionEventScope,
	m dogma.Message,
) (ok bool, err error) {
	return false, nil
}

// ResourceVersion returns the version of the resource r.
func (SecondProjectionHandler) ResourceVersion(
	ctx context.Context,
	r []byte,
) ([]byte, error) {
	return nil, nil
}

// CloseResource informs the projection that the resource r will not be
// used in any future calls to HandleEvent().
func (SecondProjectionHandler) CloseResource(ctx context.Context, r []byte) error {
	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
func (SecondProjectionHandler) TimeoutHint(m dogma.Message) time.Duration {
	return 0
}

// Compact reduces the size of the projection's data.
func (SecondProjectionHandler) Compact(ctx context.Context, s dogma.ProjectionCompactScope) error {
	return nil
}

// SecondIntegration is an integration used for testing.
type SecondIntegration struct{}

// ApplyEvent updates the integration instance to reflect the occurrence of an
// event that was recorded against this instance.
func (SecondIntegration) ApplyEvent(m dogma.Message) {}

// SecondIntegrationHandler is a test implementation of
// dogma.IntegrationMessageHandler.
type SecondIntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (SecondIntegrationHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("<second-integration>", "6bed3fbc-30e2-44c7-9a5b-e440ffe370d9")

	c.ConsumesCommandType(fixtures.MessageA{})

	c.ProducesEventType(fixtures.MessageB{})
}

// RouteCommandToInstance returns the ID of the integration instance that is
// targetted by m.
func (SecondIntegrationHandler) RouteCommandToInstance(m dogma.Message) string {
	return "<second-integration>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (SecondIntegrationHandler) HandleCommand(
	ctx context.Context,
	s dogma.IntegrationCommandScope,
	m dogma.Message,
) error {
	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
func (SecondIntegrationHandler) TimeoutHint(m dogma.Message) time.Duration {
	return 0
}
