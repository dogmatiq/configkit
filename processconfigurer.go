package configkit

import (
	"github.com/dogmatiq/dogma"
)

type processConfigurer struct {
	config *richProcess
}

func (c *processConfigurer) Identity(name, key string) {
	configureIdentity(
		&c.config.ident,
		name,
		key,
		c.config.ReflectType(),
	)
}

func (c *processConfigurer) Routes(routes ...dogma.ProcessRoute) {
	for _, route := range routes {
		configureRoute(
			&c.config.types,
			route,
			c.config.ident,
			c.config.ReflectType(),
		)
	}
}

func (c *processConfigurer) Disable(...dogma.DisableOption) {
	c.config.isDisabled = true
}
