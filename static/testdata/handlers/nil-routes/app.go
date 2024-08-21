package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<app>", "c100edcc-6dcc-42ed-ac75-69eecb3d0ec4")
	c.RegisterIntegration(IntegrationHandler{})
}
