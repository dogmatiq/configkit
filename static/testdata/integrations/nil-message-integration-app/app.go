package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<nil-message-integraion-app>", "dc0c1986-953b-445e-a5d0-3240c22c1abf")

	c.RegisterIntegration(IntegrationHandler{})
}
