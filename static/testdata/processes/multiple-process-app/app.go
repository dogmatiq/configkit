package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<multiple-process-app>", "66abab58-e0cd-4d50-90c2-b0bd093ea2ba")

	c.RegisterProcess(FirstProcessHandler{})
	c.RegisterProcess(SecondProcessHandler{})
}
