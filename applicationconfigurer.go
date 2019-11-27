package configkit

import "github.com/dogmatiq/dogma"

// applicationConfigurer is an implementation of dogma.ApplicationConfigurer.
type applicationConfigurer struct {
	entityConfigurer

	target *application
}

func (c *applicationConfigurer) RegisterAggregate(h dogma.AggregateMessageHandler) {
	cfg := FromAggregate(h)
	c.register(cfg)
}

func (c *applicationConfigurer) RegisterProcess(h dogma.ProcessMessageHandler) {
	cfg := FromProcess(h)
	c.register(cfg)
}

func (c *applicationConfigurer) RegisterIntegration(h dogma.IntegrationMessageHandler) {
	cfg := FromIntegration(h)
	c.register(cfg)
}

func (c *applicationConfigurer) RegisterProjection(h dogma.ProjectionMessageHandler) {
	cfg := FromProjection(h)
	c.register(cfg)
}

func (c *applicationConfigurer) register(h RichHandler) {
	if c.target.handlers == nil {
		c.target.handlers = HandlerSet{}
		c.target.richHandlers = RichHandlerSet{}
	}

	c.target.handlers.Add(h)
	c.target.richHandlers.Add(h)
}
