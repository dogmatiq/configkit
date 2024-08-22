package app

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// FirstAggregate is an aggregate used for testing.
type FirstAggregate struct{}

// ApplyEvent updates the aggregate instance to reflect the occurrence of an
// event that was recorded against this instance.
func (FirstAggregate) ApplyEvent(dogma.Event) {}

// FirstAggregateHandler is a test implementation of dogma.AggregateMessageHandler.
type FirstAggregateHandler struct{}

// New returns a new account instance.
func (FirstAggregateHandler) New() dogma.AggregateRoot {
	return FirstAggregate{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (FirstAggregateHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("<first-aggregate>", "e6300d8d-6530-405e-9729-e9ca21df23d3")

	c.Routes(
		dogma.HandlesCommand[fixtures.MessageA](),
		dogma.RecordsEvent[fixtures.MessageB](),
	)
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (FirstAggregateHandler) RouteCommandToInstance(dogma.Command) string {
	return "<first-aggregate>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (FirstAggregateHandler) HandleCommand(
	dogma.AggregateRoot,
	dogma.AggregateCommandScope,
	dogma.Command,
) {
}

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

	c.Routes(
		dogma.HandlesEvent[fixtures.MessageA](),
		dogma.ExecutesCommand[fixtures.MessageB](),
		dogma.SchedulesTimeout[fixtures.MessageC](),
	)
}

// RouteEventToInstance returns the ID of the process instance that is
// targeted by m.
func (FirstProcessHandler) RouteEventToInstance(
	context.Context,
	dogma.Event,
) (string, bool, error) {
	return "<first-process>", true, nil
}

// HandleEvent handles an event message.
func (FirstProcessHandler) HandleEvent(
	context.Context,
	dogma.ProcessRoot,
	dogma.ProcessEventScope,
	dogma.Event,
) error {
	return nil
}

// HandleTimeout handles a timeout message that has been scheduled with
// ProcessScope.ScheduleTimeout().
func (FirstProcessHandler) HandleTimeout(
	context.Context,
	dogma.ProcessRoot,
	dogma.ProcessTimeoutScope,
	dogma.Timeout,
) error {
	return nil
}

// FirstProjectionHandler is a test implementation of
// dogma.ProjectionMessageHandler.
type FirstProjectionHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (FirstProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("<first-projection>", "9174783f-4f12-4619-b5c6-c4ab70bd0937")

	c.Routes(
		dogma.HandlesEvent[fixtures.MessageA](),
	)
}

// HandleEvent updates the projection to reflect the occurrence of an event.
func (FirstProjectionHandler) HandleEvent(
	_ context.Context,
	_, _, _ []byte,
	_ dogma.ProjectionEventScope,
	_ dogma.Event,
) (ok bool, err error) {
	return false, nil
}

// ResourceVersion returns the version of the resource r.
func (FirstProjectionHandler) ResourceVersion(context.Context, []byte) ([]byte, error) {
	return nil, nil
}

// CloseResource informs the projection that the resource r will not be
// used in any future calls to HandleEvent().
func (FirstProjectionHandler) CloseResource(context.Context, []byte) error {
	return nil
}

// Compact reduces the size of the projection's data.
func (FirstProjectionHandler) Compact(context.Context, dogma.ProjectionCompactScope) error {
	return nil
}

// FirstIntegrationHandler is a test implementation of
// dogma.IntegrationMessageHandler.
type FirstIntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (FirstIntegrationHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("<first-integration>", "14cf2812-eead-43b3-9c9c-10db5b469e94")

	c.Routes(
		dogma.HandlesCommand[fixtures.MessageA](),
		dogma.RecordsEvent[fixtures.MessageB](),
	)
}

// RouteCommandToInstance returns the ID of the integration instance that is
// targetted by m.
func (FirstIntegrationHandler) RouteCommandToInstance(dogma.Command) string {
	return "<first-integration>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (FirstIntegrationHandler) HandleCommand(
	context.Context,
	dogma.IntegrationCommandScope,
	dogma.Command,
) error {
	return nil
}
