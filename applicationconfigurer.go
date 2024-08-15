package configkit

import (
	"github.com/dogmatiq/configkit/internal/validation"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

// applicationConfigurer is an implementation of dogma.ApplicationConfigurer.
type applicationConfigurer struct {
	entityConfigurer

	app *richApplication
}

func (c *applicationConfigurer) Identity(name, key string) {
	if h, ok := c.app.richHandlers.ByKey(key); ok {
		validation.Panicf(
			`%s can not use the application key "%s", because it is already used by %s`,
			c.entity.rt,
			key,
			h.ReflectType(),
		)
	}

	c.entityConfigurer.Identity(name, key)
}

func (c *applicationConfigurer) RegisterAggregate(h dogma.AggregateMessageHandler, _ ...dogma.RegisterAggregateOption) {
	c.registerIfConfigured(fromAggregate(h))
}

func (c *applicationConfigurer) RegisterProcess(h dogma.ProcessMessageHandler, _ ...dogma.RegisterProcessOption) {
	c.registerIfConfigured(fromProcess(h))
}

func (c *applicationConfigurer) RegisterIntegration(h dogma.IntegrationMessageHandler, _ ...dogma.RegisterIntegrationOption) {
	c.registerIfConfigured(fromIntegration(h))
}

func (c *applicationConfigurer) RegisterProjection(h dogma.ProjectionMessageHandler, _ ...dogma.RegisterProjectionOption) {
	c.registerIfConfigured(fromProjection(h))
}

// register adds a handler configuration to the application.
func (c *applicationConfigurer) registerIfConfigured(
	h RichHandler,
	hc interface {
		isConfigured() bool
		mustValidate()
	},
) {
	if h.IsDisabled() && !hc.isConfigured() {
		return
	}

	hc.mustValidate()

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

// guardAgainstConflictingIdentities panics if h's identity conflicts with the
// application or any other handlers.
func (c *applicationConfigurer) guardAgainstConflictingIdentities(h RichHandler) {
	i := h.Identity()

	if i.Key == c.entity.ident.Key {
		validation.Panicf(
			`%s can not use the handler key "%s", because it is already used by %s`,
			h.ReflectType(),
			i.Key,
			c.entity.rt,
		)
	}

	if x, ok := c.app.richHandlers.ByName(i.Name); ok {
		validation.Panicf(
			`%s can not use the handler name "%s", because it is already used by %s`,
			h.ReflectType(),
			i.Name,
			x.ReflectType(),
		)
	}

	if x, ok := c.app.richHandlers.ByKey(i.Key); ok {
		validation.Panicf(
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
	for mt, r := range h.MessageTypes().All() {
		xr, ok := c.entity.types.RoleOf(mt)

		if !ok || xr == r {
			continue
		}

		// we know there's a conflict, now we just need to find a handler that
		// refers to this message type as some other role.
		xh, _ := c.app.richHandlers.Find(func(h RichHandler) bool {
			x, ok := h.MessageTypes().RoleOf(mt)
			return ok && x != r
		})

		validation.Panicf(
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
			validation.Panicf(
				`%s (%s) can not handle %s commands because they are already configured to be handled by %s (%s)`,
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
			validation.Panicf(
				`%s (%s) can not record %s events because they are already configured to be recorded by %s (%s)`,
				h.ReflectType(),
				h.Identity().Name,
				mt,
				x.ReflectType(),
				x.Identity().Name,
			)
		}
	}
}
