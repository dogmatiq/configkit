package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<only-routes-args>", "f2c08525-623e-4c76-851c-3172953269e3")

	c.RegisterIntegration(IntegrationHandler{})
	c.RegisterProjection(ProjectionHandler{})
	c.RegisterAggregate(AggregateHandler{})
	c.RegisterProcess(ProcessHandler{})
}
