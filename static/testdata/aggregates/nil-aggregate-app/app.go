package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<nil-aggregate-app>", "250738eb-7728-4414-a850-1afefb04dbb4")

	c.RegisterAggregate(nil)
}
