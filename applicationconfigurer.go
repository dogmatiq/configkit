package configkit

import "github.com/dogmatiq/dogma"

// applicationConfigurer is an implementation of dogma.ApplicationConfigurer.
type applicationConfigurer struct {
	entityConfigurer
}

func (c *applicationConfigurer) RegisterAggregate(h dogma.AggregateMessageHandler) {

}

func (c *applicationConfigurer) RegisterProcess(h dogma.ProcessMessageHandler) {

}

func (c *applicationConfigurer) RegisterIntegration(h dogma.IntegrationMessageHandler) {

}

func (c *applicationConfigurer) RegisterProjection(h dogma.ProjectionMessageHandler) {

}
