package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<nil-message-projection-app>", "d6ff7e9a-f070-4447-8177-fa29fa637de7")

	c.RegisterProjection(ProjectionHandler{})
}
