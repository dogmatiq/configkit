package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<single-integration-app>", "cb5558eb-451a-4df3-8290-96e29ed793e7")

	c.RegisterIntegration(IntegrationHandler{})
}
