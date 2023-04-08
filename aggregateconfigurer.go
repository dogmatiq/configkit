package configkit

import "github.com/dogmatiq/dogma"

type aggregateConfigurer struct {
	handlerConfigurer
}

func (c *aggregateConfigurer) Routes(routes ...dogma.AggregateRoute) {
	for _, r := range routes {
		c.route(r)
	}
}
