package apps

import "github.com/dogmatiq/dogma"

// AppFirst implements dogma.Application interface.
type AppFirst struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (AppFirst) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<app-first>", "4fec74a1-6ed4-46f4-8417-01e0910be8f1")
}
