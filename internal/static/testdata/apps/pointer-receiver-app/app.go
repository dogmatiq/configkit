package app

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure is implemented using a pointer receiver, such that the *App
// implements dogma.Application, and not App itself.
func (a *App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<app>", "b754902b-47c8-48fc-84d2-d920c9cbdaec")
}
