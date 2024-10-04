package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<app-level-messages-app>", "2ecbd06f-6a1c-4dd8-81fc-23af38910f2b")

	c.RegisterIntegration(IntegrationHandler{})
	c.RegisterProjection(ProjectionHandler{})
	c.RegisterAggregate(AggregateHandler{})
	c.RegisterProcess(ProcessHandler{})
}
