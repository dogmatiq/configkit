package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<multiple-integration-app>", "c01a4fef-0979-4468-bd7e-c9d371f82cd2")

	c.RegisterIntegration(FirstIntegrationHandler{})
	c.RegisterIntegration(SecondIntegrationHandler{})
}
