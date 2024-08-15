package configkit

import (
	"github.com/dogmatiq/configkit/internal/validation"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

type projectionConfigurer struct {
	handlerConfigurer
	projection *richProjection
}

func (c *projectionConfigurer) Routes(routes ...dogma.ProjectionRoute) {
	for _, r := range routes {
		c.route(r)
	}
}

func (c *projectionConfigurer) DeliveryPolicy(p dogma.ProjectionDeliveryPolicy) {
	c.configured = true

	if p == nil {
		validation.Panicf("delivery policy must not be nil")
	}

	c.projection.deliveryPolicy = p
}

func (c *projectionConfigurer) mustValidate() {
	c.handlerConfigurer.mustValidate()
	c.mustConsume(message.EventRole)
}
