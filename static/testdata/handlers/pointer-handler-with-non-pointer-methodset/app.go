package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<non-pointer-handler-registered-as-pointer>", "282653ad-9343-44f1-889e-a8b2b095b54b")

	c.RegisterIntegration(&IntegrationHandler{})
	c.RegisterProjection(&ProjectionHandler{})
	c.RegisterAggregate(&AggregateHandler{})
	c.RegisterProcess(&ProcessHandler{})
}
