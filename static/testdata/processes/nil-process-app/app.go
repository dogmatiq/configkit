package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<nil-process-app>", "ff9a8436-2d8f-4cb8-ae1b-639a3e4d859f")

	c.RegisterProcess(nil)
}
