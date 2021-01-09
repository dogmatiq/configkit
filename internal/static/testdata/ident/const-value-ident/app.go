package ident

import "github.com/dogmatiq/dogma"

const (
	// AppName is the application name.
	AppName = "<app>"
	// AppKey is the application key.
	AppKey = "04e12cf2-3c66-4414-9203-e045ddbe02c7"
)

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(AppName, AppKey)
}
