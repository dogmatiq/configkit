package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<nil-message-aggregate-app>", "77199eb1-c893-4151-aab8-5088f6db4ea6")

	c.RegisterAggregate(AggregateHandler{})
}
