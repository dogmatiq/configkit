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
		c.config.ReflectType(),
		&c.config.ident,
		name,
		key,
	)
}

func (c *aggregateConfigurer) Routes(routes ...dogma.AggregateRoute) {
	for _, r := range routes {
		configureRoute(
			c.config.ReflectType(),
			c.config.ident,
			&c.config.types,
			r,
		)
	}
}

func (c *aggregateConfigurer) Disable(...dogma.DisableOption) {
	c.config.isDisabled = true
}
