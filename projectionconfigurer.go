package configkit

import (
	"github.com/dogmatiq/configkit/internal/validation"
	"github.com/dogmatiq/dogma"
)

type projectionConfigurer struct {
	handlerConfigurer
	deliveryPolicy dogma.ProjectionDeliveryPolicy
}

func (c *projectionConfigurer) Routes(routes ...dogma.ProjectionRoute) {
	for _, r := range routes {
		c.route(r)
	}
}

func (c *projectionConfigurer) DeliveryPolicy(p dogma.ProjectionDeliveryPolicy) {
	if p == nil {
		validation.Panicf("delivery policy must not be nil")
	}

	c.deliveryPolicy = p
}
