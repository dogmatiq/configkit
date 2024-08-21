package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<adaptor-func-app>", "f610eae4-f5d0-4eea-a9c9-6cbbfa9b2060")

	c.RegisterAggregate(AdaptAggregate(AggregateHandler{}))
	c.RegisterProcess(AdaptProcess(ProcessHandler{}))
	c.RegisterProjection(AdaptProjection(ProjectionHandler{}))
	c.RegisterIntegration(AdaptIntegration(IntegrationHandler{}))
}
