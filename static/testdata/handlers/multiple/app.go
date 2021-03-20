package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<multiple-handler-of-a-kind-app>", "8961f548-1afc-4996-894c-956835c83199")

	c.RegisterAggregate(FirstAggregateHandler{})
	c.RegisterAggregate(SecondAggregateHandler{})

	c.RegisterProcess(FirstProcessHandler{})
	c.RegisterProcess(SecondProcessHandler{})

	c.RegisterProjection(FirstProjectionHandler{})
	c.RegisterProjection(SecondProjectionHandler{})

	c.RegisterIntegration(FirstIntegrationHandler{})
	c.RegisterIntegration(SecondIntegrationHandler{})
}
