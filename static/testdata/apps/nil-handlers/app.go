package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<nil-handler-app>", "0726ae0d-67e4-4a50-8a19-9f58eae38e51")

	c.RegisterAggregate(nil)
	c.RegisterProcess(nil)
	c.RegisterProjection(nil)
	c.RegisterIntegration(nil)
}
