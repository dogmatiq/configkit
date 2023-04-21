package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<single-handler-of-a-kind-app-calling-deprecated-methods>", "4485488f-7fa8-42d7-9b60-7462ade4db9b")

	c.RegisterIntegration(IntegrationHandler{})
	c.RegisterProjection(ProjectionHandler{})
	c.RegisterAggregate(AggregateHandler{})
	c.RegisterProcess(ProcessHandler{})
}
