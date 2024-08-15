package configkit

import (
	"github.com/dogmatiq/configkit/internal/validation"
	"github.com/dogmatiq/dogma"
)

type projectionConfigurer struct {
	config *richProjection
}

func (c *projectionConfigurer) Identity(name, key string) {
	configureIdentity(
		&c.config.ident,
		name,
		key,
		c.config.ReflectType(),
	)
}

func (c *projectionConfigurer) Routes(routes ...dogma.ProjectionRoute) {
	for _, route := range routes {
		configureRoute(
			&c.config.types,
			route,
			c.config.ident,
			c.config.ReflectType(),
		)
	}
}

func (c *projectionConfigurer) Disable(...dogma.DisableOption) {
	c.config.isDisabled = true
}

func (c *projectionConfigurer) DeliveryPolicy(p dogma.ProjectionDeliveryPolicy) {
	if p == nil {
		validation.Panicf("delivery policy must not be nil")
	}
	c.config.deliveryPolicy = p
}
