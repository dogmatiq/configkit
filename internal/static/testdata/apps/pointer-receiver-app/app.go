package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
//
// Note that this method uses a pointer receiver to test the detection of the
// dogma.Application implementation with pointer receivers.
func (a *App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<app>", "b754902b-47c8-48fc-84d2-d920c9cbdaec")
}
