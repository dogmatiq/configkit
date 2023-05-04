package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<conditional-branch-app>", "7e34538e-c407-4af8-8d3c-960e09cde98a")

	c.RegisterIntegration(IntegrationHandler{})
	c.RegisterProjection(ProjectionHandler{})
	c.RegisterAggregate(AggregateHandler{})
	c.RegisterProcess(ProcessHandler{})
}
