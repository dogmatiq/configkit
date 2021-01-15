package ident

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct {
	Name string
	Key  string
}

// Configure sets the application identity using non-constant expressions.
func (a App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(a.Name, a.Name)
}
