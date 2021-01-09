package ident

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct {
	Name string
	Key  string
}

// Configure configures the behavior of the engine as it relates to this
// application.
func (a App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(a.Name, a.Name)
}
