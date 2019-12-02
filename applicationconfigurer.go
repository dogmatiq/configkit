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

func (c *applicationConfigurer) Identity(name, key string) {
	if h, ok := c.app.richHandlers.ByKey(key); ok {
		Panicf(
			`%s can not use the application key "%s", because it is already used by %s`,
			c.entity.rt,
			key,
			h.ReflectType(),
		)
	}

	c.entityConfigurer.Identity(name, key)
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

// register adds a handler configuration to the application.
func (c *applicationConfigurer) register(h RichHandler) {
	c.guardAgainstConflictingIdentities(h)
	c.guardAgainstConflictingRoles(h)
	c.guardAgainstConflictingRoutes(h)

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

// isForeign returns true if mt is "foreign", meaning that it needs to be
// obtained from or sent outside the application.
func (c *applicationConfigurer) isForeign(mt message.Type, r message.Role) bool {
	produced := c.entity.types.Produced.Has(mt)
	consumed := c.entity.types.Consumed.Has(mt)

	// Commands must always have a handler, therefore any command that is not
	// both produced and consumed by this application is considered foreign.
	if r == message.CommandRole {
		return produced != consumed
	}

	// Events that are only considered foreign if they need to be obtained from
	// elsewhere. There is no requirement that a given event MUST have a
	// handler, therefore any events produced by this app but not consumed are
	// not considered foreign.
	return consumed && !produced
}

// guardAgainstConflictingIdentities panics if h's identity conflicts with the
// application or any other handlers.
func (c *applicationConfigurer) guardAgainstConflictingIdentities(h RichHandler) {
	i := h.Identity()

	if i.Key == c.entity.ident.Key {
		Panicf(
			`%s can not use the handler key "%s", because it is already used by %s`,
			h.ReflectType(),
			i.Key,
			c.entity.rt,
		)
	}

	if x, ok := c.app.richHandlers.ByName(i.Name); ok {
		Panicf(
			`%s can not use the handler name "%s", because it is already used by %s`,
			h.ReflectType(),
			i.Name,
			x.ReflectType(),
		)
	}

	if x, ok := c.app.richHandlers.ByKey(i.Key); ok {
		Panicf(
			`%s can not use the handler key "%s", because it is already used by %s`,
			h.ReflectType(),
			i.Key,
			x.ReflectType(),
		)
	}
}

// guardAgainstConflictingRoles panics if h configures any messages in roles
// contrary to the way they are configured by any other handler.
func (c *applicationConfigurer) guardAgainstConflictingRoles(h RichHandler) {
	for mt, r := range h.MessageTypes().Roles {
		xr, ok := c.entity.types.Roles[mt]

		if !ok || xr == r {
			continue
		}

		// we know there's a conflict, now we just need to find a handler that
		// refers to this message type as some other role.
		xh, _ := c.app.richHandlers.Find(func(h RichHandler) bool {
			return h.MessageTypes().Roles[mt] != r
		})

		Panicf(
			`%s (%s) configures %s as a %s but %s (%s) configures it as a %s`,
			h.ReflectType(),
			h.Identity().Name,
			mt,
			r,
			xh.ReflectType(),
			xh.Identity().Name,
			xr,
		)
	}
}

// guardAgainstConflictingRoutes panics if an h consumes the same commands or
// produces the same events as some existing handler.
func (c *applicationConfigurer) guardAgainstConflictingRoutes(h RichHandler) {
	types := h.MessageTypes()

	for mt, r := range types.Consumed {
		if r != message.CommandRole {
			continue
		}

		for _, x := range c.app.richHandlers.ConsumersOf(mt) {
			Panicf(
				`%s (%s) can not consume %s commands because they are already consumed by %s (%s)`,
				h.ReflectType(),
				h.Identity().Name,
				mt,
				x.ReflectType(),
				x.Identity().Name,
			)
		}
	}

	for mt, r := range types.Produced {
		if r != message.EventRole {
			continue
		}

		for _, x := range c.app.richHandlers.ProducersOf(mt) {
			Panicf(
				`%s (%s) can not produce %s events because they are already produced by %s (%s)`,
				h.ReflectType(),
				h.Identity().Name,
				mt,
				x.ReflectType(),
				x.Identity().Name,
			)
		}
	}
}
