package configkit

import (
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

type integrationConfigurer struct {
	handlerConfigurer
}

func (c *integrationConfigurer) Routes(routes ...dogma.IntegrationRoute) {
	for _, r := range routes {
		c.route(r)
	}
}

func (c *integrationConfigurer) mustValidate() {
	c.handlerConfigurer.mustValidate()
	c.mustConsume(message.CommandRole)
}
