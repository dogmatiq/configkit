package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<nil-integration-app>", "e0075708-ca66-40e1-baf8-f5787e8d0a38")

	c.RegisterIntegration(nil)
}
