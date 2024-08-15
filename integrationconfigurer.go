package configkit

import (
	"github.com/dogmatiq/dogma"
)

type integrationConfigurer struct {
	config *richIntegration
}

func (c *integrationConfigurer) Identity(name, key string) {
	configureIdentity(
		&c.config.ident,
		name,
		key,
		c.config.ReflectType(),
	)
}

func (c *integrationConfigurer) Routes(routes ...dogma.IntegrationRoute) {
	for _, route := range routes {
		configureRoute(
			&c.config.types,
			route,
			c.config.ident,
			c.config.ReflectType(),
		)
	}
}

func (c *integrationConfigurer) Disable(...dogma.DisableOption) {
	c.config.isDisabled = true
}
