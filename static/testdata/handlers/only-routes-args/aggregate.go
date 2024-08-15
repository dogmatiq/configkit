package app

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

// Aggregate is an aggregate used for testing.
type Aggregate struct{}

// ApplyEvent updates the aggregate instance to reflect the occurrence of an
// event that was recorded against this instance.
func (Aggregate) ApplyEvent(dogma.Event) {}

// AggregateHandler is a test implementation of dogma.AggregateMessageHandler.
type AggregateHandler struct{}

// New returns a new account instance.
func (AggregateHandler) New() dogma.AggregateRoot {
	return Aggregate{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (AggregateHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("<aggregate>", "dcfdd034-e374-478b-8faa-bc688ff59f1f")

	c.Routes(
		dogma.HandlesCommand[fixtures.MessageA](),
		dogma.HandlesCommand[fixtures.MessageB](),
		dogma.RecordsEvent[fixtures.MessageC](),
		dogma.RecordsEvent[fixtures.MessageD](),
	)

	// These assignments should not be included into handler configuration
	// as they are not arguments to the Routes() method.
	var r1 dogma.AggregateRoute = dogma.RecordsEvent[fixtures.MessageE]()
	var r2 dogma.AggregateRoute = dogma.HandlesCommand[fixtures.MessageF]()
	// Assign to blank identifiers to avoid unused variable errors during
	// compilation.
	_ = r1
	_ = r2
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (AggregateHandler) RouteCommandToInstance(dogma.Command) string {
	return "<aggregate>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (AggregateHandler) HandleCommand(
	dogma.AggregateRoot,
	dogma.AggregateCommandScope,
	dogma.Command,
) {
}
