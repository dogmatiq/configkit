package app01

import "github.com/dogmatiq/dogma"

const (
	// AppKey is the application key.
	AppKey = "7c3c67dd-6c0b-4952-97d5-54c75fc7a1c6"
)

// Below is the deliberate illegal Go syntax to test loading of the packages
// with errors.
This is deliberate illegal syntax


// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app01", AppKey)
}
