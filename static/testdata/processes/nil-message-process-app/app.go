package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<nil-message-process-app>", "1003442d-620a-47bb-8011-596a45401fc5")

	c.RegisterProcess(ProcessHandler{})
}
