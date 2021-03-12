package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<multiple-projection-app>", "e9d3c0df-f527-4604-a6c6-c8bc96a56869")

	c.RegisterProjection(FirstProjectionHandler{})
	c.RegisterProjection(SecondProjectionHandler{})
}
