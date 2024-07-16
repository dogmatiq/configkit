package configkit

import (
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

type processConfigurer struct {
	handlerConfigurer
}

func (c *processConfigurer) Routes(routes ...dogma.ProcessRoute) {
	for _, r := range routes {
		c.route(r)
	}
}

func (c *processConfigurer) mustValidate() {
	c.handlerConfigurer.mustValidate()
	c.mustConsume(message.EventRole)
	c.mustProduce(message.CommandRole)
}
