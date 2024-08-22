package app

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// SecondAggregate is an aggregate used for testing.
type SecondAggregate struct{}

// ApplyEvent updates the aggregate instance to reflect the occurrence of an
// event that was recorded against this instance.
func (SecondAggregate) ApplyEvent(dogma.Event) {}

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

	c.Routes(
		dogma.HandlesCommand[fixtures.MessageC](),
		dogma.RecordsEvent[fixtures.MessageD](),
	)
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (SecondAggregateHandler) RouteCommandToInstance(dogma.Command) string {
	return "<second-aggregate>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (SecondAggregateHandler) HandleCommand(
	dogma.AggregateRoot,
	dogma.AggregateCommandScope,
	dogma.Command,
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

	c.Routes(
		dogma.HandlesEvent[fixtures.MessageD](),
		dogma.ExecutesCommand[fixtures.MessageE](),
		dogma.SchedulesTimeout[fixtures.MessageF](),
	)
}

// RouteEventToInstance returns the ID of the process instance that is
// targeted by m.
func (SecondProcessHandler) RouteEventToInstance(
	context.Context,
	dogma.Event,
) (string, bool, error) {
	return "<second-process>", true, nil
}

// HandleEvent handles an event message.
func (SecondProcessHandler) HandleEvent(
	context.Context,
	dogma.ProcessRoot,
	dogma.ProcessEventScope,
	dogma.Event,
) error {
	return nil
}

// HandleTimeout handles a timeout message that has been scheduled with
// ProcessScope.ScheduleTimeout().
func (SecondProcessHandler) HandleTimeout(
	context.Context,
	dogma.ProcessRoot,
	dogma.ProcessTimeoutScope,
	dogma.Timeout,
) error {
	return nil
}

// SecondProjectionHandler is a test implementation of
// dogma.ProjectionMessageHandler.
type SecondProjectionHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (SecondProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("<second-projection>", "2e22850e-7c84-4b3f-b8b3-25ac743d90f2")

	c.Routes(
		dogma.HandlesEvent[fixtures.MessageB](),
	)
}

// HandleEvent updates the projection to reflect the occurrence of an event.
func (SecondProjectionHandler) HandleEvent(
	_ context.Context,
	_, _, _ []byte,
	_ dogma.ProjectionEventScope,
	_ dogma.Event,
) (ok bool, err error) {
	return false, nil
}

// ResourceVersion returns the version of the resource r.
func (SecondProjectionHandler) ResourceVersion(context.Context, []byte) ([]byte, error) {
	return nil, nil
}

// CloseResource informs the projection that the resource r will not be
// used in any future calls to HandleEvent().
func (SecondProjectionHandler) CloseResource(context.Context, []byte) error {
	return nil
}

// Compact reduces the size of the projection's data.
func (SecondProjectionHandler) Compact(context.Context, dogma.ProjectionCompactScope) error {
	return nil
}

// SecondIntegrationHandler is a test implementation of
// dogma.IntegrationMessageHandler.
type SecondIntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (SecondIntegrationHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("<second-integration>", "6bed3fbc-30e2-44c7-9a5b-e440ffe370d9")

	c.Routes(
		dogma.HandlesCommand[fixtures.MessageA](),
		dogma.RecordsEvent[fixtures.MessageB](),
	)
}

// RouteCommandToInstance returns the ID of the integration instance that is
// targetted by m.
func (SecondIntegrationHandler) RouteCommandToInstance(dogma.Command) string {
	return "<second-integration>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (SecondIntegrationHandler) HandleCommand(
	context.Context,
	dogma.IntegrationCommandScope,
	dogma.Command,
) error {
	return nil
}
