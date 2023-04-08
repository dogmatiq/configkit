package configkit

import "github.com/dogmatiq/dogma"

type processConfigurer struct {
	handlerConfigurer
}

func (c *processConfigurer) Routes(routes ...dogma.ProcessRoute) {
	for _, r := range routes {
		c.route(r)
	}
}
