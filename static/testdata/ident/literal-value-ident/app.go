package ident

import "github.com/dogmatiq/dogma"

// App implements dogma.Application interface.
type App struct{}

// Configure sets the application identity using literal string values.
func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<app>", "9d0af85d-f506-4742-b676-ce87730bb1a0")
}
