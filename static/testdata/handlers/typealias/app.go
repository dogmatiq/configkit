package app

import "github.com/dogmatiq/dogma"

type (
	// IntergrationHandlerAlias is a test type alias of IntegrationHandler.
	IntergrationHandlerAlias = IntegrationHandler
	// ProjectionHandlerAlias is a test type alias of ProjectionHandler.
	ProjectionHandlerAlias = ProjectionHandler
	// AggregateHandlerAlias is a test type alias of AggregateHandler.
	AggregateHandlerAlias = AggregateHandler
	// ProcessHandlerAlias is a test type alias of ProcessHandler.
	ProcessHandlerAlias = ProcessHandler
)

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<handler-as-typealias-app>", "1b828a1c-eba1-4e4c-88b8-e49f78ad15c7")

	c.RegisterIntegration(IntergrationHandlerAlias{})
	c.RegisterProjection(ProjectionHandlerAlias{})
	c.RegisterAggregate(AggregateHandlerAlias{})
	c.RegisterProcess(ProcessHandlerAlias{})
}
