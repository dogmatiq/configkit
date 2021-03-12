package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<nil-projection-app>", "aa8263fd-6333-41f0-8cd4-6a72d467cedf")

	c.RegisterProjection(nil)
}
