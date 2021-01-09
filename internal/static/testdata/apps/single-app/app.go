package app

import "github.com/dogmatiq/dogma"

const (
	// AppKey is the application key.
	AppKey = "8a6baab1-ee64-402e-a081-e43f4bebc243"
)

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<app>", AppKey)
}
