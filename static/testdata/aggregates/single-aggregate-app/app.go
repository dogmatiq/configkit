package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<single-aggregate-app>", "3bc3849b-abe0-4c4e-9db4-e48dc28c9a26")

	c.RegisterAggregate(AggregateHandler{})
}
