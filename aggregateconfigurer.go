package configkit

import (
	"github.com/dogmatiq/dogma"
)

// aggregateConfigurer is the default implementation of
// [dogma.AggregateConfigurer].
type aggregateConfigurer struct {
	config *richAggregate
}

func (c *aggregateConfigurer) Identity(name, key string) {
	configureIdentity(
		&c.config.ident,
		name,
		key,
		c.config.ReflectType(),
	)
}

func (c *aggregateConfigurer) Routes(routes ...dogma.AggregateRoute) {
	for _, route := range routes {
		configureRoute(
			&c.config.types,
			route,
			c.config.ident,
			c.config.ReflectType(),
		)
	}
}

func (c *aggregateConfigurer) Disable(...dogma.DisableOption) {
	c.config.isDisabled = true
}
