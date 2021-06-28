package app

import (
	"github.com/dogmatiq/dogma"
)

// App implements dogma.Application interface.
type App struct {
	Aggregate AggregateHandler
}

// Configure configures the behavior of the engine as it relates to this
// application.
func (a App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<app>", "7468a57f-20f0-4d11-9aad-48fcd553a908")

	c.RegisterAggregate(&a.Aggregate)
}
