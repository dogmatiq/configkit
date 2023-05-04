package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<nil-message-handler-app>", "68081b6e-af25-4522-b54b-88ef9759c5f2")

	c.RegisterIntegration(IntegrationHandler{})
	c.RegisterProjection(ProjectionHandler{})
	c.RegisterAggregate(AggregateHandler{})
	c.RegisterProcess(ProcessHandler{})
}
