package app01

import "github.com/dogmatiq/dogma"

const (
	// AppKey is the application key.
	AppKey = "b754902b-47c8-48fc-84d2-d920c9cbdaec"
)

// App implements dogma.Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app01", AppKey)
}
