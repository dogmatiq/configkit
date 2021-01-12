package apps

import "github.com/dogmatiq/dogma"

// AppSecond implements dogma.Application interface.
type AppSecond struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (a AppSecond) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<app-second>", "6e97d403-3cb8-4a59-a7ec-74e8e219a7bc")
}
