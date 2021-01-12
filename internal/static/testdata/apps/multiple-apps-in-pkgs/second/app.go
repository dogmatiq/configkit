package second

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (a App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<app-second>", "bfaf2a16-23a0-495d-8098-051d77635822")
}
