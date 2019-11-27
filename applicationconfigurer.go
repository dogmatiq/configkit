package configkit

import (
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

// applicationConfigurer is an implementation of dogma.ApplicationConfigurer.
type applicationConfigurer struct {
	entityConfigurer

	app *application
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
	if c.app.handlers == nil {
		c.app.handlers = HandlerSet{}
		c.app.richHandlers = RichHandlerSet{}
	}

	c.app.handlers.Add(h)
	c.app.richHandlers.Add(h)

	types := h.MessageTypes()

	for mt, r := range types.Roles {
		if c.entity.names.Roles == nil {
			c.entity.names.Roles = message.NameRoles{}
			c.entity.types.Roles = message.TypeRoles{}
		}

		c.entity.names.Roles.Add(mt.Name(), r)
		c.entity.types.Roles.Add(mt, r)
	}

	for mt, r := range types.Produced {
		if c.entity.names.Produced == nil {
			c.entity.names.Produced = message.NameRoles{}
			c.entity.types.Produced = message.TypeRoles{}
		}

		c.entity.names.Produced.Add(mt.Name(), r)
		c.entity.types.Produced.Add(mt, r)
	}

	for mt, r := range types.Consumed {
		if c.entity.names.Consumed == nil {
			c.entity.names.Consumed = message.NameRoles{}
			c.entity.types.Consumed = message.TypeRoles{}
		}

		c.entity.names.Consumed.Add(mt.Name(), r)
		c.entity.types.Consumed.Add(mt, r)
	}
}

func (c *applicationConfigurer) validate() {
	c.entityConfigurer.validate()

	for mt, r := range c.entity.types.Roles {
		if c.isForeign(mt, r) {
			if c.app.foreignNames == nil {
				c.app.foreignNames = message.NameRoles{}
				c.app.foreignTypes = message.TypeRoles{}
			}

			c.app.foreignNames.Add(mt.Name(), r)
			c.app.foreignTypes.Add(mt, r)
		}
	}
}

// isForeign returns true if mt is "foreign", meaning that it needs to be
// obtained from or sent to a different application.
func (c *applicationConfigurer) isForeign(mt message.Type, r message.Role) bool {
	produced := c.entity.types.Produced.Has(mt)
	consumed := c.entity.types.Consumed.Has(mt)

	switch r {
	case message.CommandRole:
		return produced && !consumed
	case message.EventRole:
		return consumed && !produced
	default:
		return false
	}
}
