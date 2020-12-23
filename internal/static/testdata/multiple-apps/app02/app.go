package app02

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
//
// Note that this method has POINTER RECEIVER, it is required to check the
// functionality of detecting the implementation of the `dogma.Application`
// interface. Please do not change it.
func (a *App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app02", "bfaf2a16-23a0-495d-8098-051d77635822")
}
