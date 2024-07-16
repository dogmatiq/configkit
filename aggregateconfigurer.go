package configkit

import (
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

type aggregateConfigurer struct {
	handlerConfigurer
}

func (c *aggregateConfigurer) Routes(routes ...dogma.AggregateRoute) {
	for _, r := range routes {
		c.route(r)
	}
}

func (c *aggregateConfigurer) mustValidate() {
	c.handlerConfigurer.mustValidate()
	c.mustConsume(message.CommandRole)
	c.mustProduce(message.EventRole)
}
