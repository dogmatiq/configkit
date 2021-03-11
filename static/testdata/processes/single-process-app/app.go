package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<single-process-app>", "662cfe49-5fff-4924-a5d6-a32255608686")

	c.RegisterProcess(ProcessHandler{})
}
