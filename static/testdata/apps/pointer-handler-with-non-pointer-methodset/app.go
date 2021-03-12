package app

import (
	"github.com/dogmatiq/dogma"
)

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (a App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<app>", "dc3ef830-ee59-472a-a8d7-523f0570a918")

	c.RegisterAggregate(&AggregateHandler{})
}
