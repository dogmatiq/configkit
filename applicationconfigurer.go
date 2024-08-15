package configkit

import (
	"github.com/dogmatiq/configkit/internal/validation"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

// applicationConfigurer is the default implementation of
// [dogma.ApplicationConfigurer].
type applicationConfigurer struct {
	config *richApplication
}

func (c *applicationConfigurer) Identity(name, key string) {
	if h, ok := c.config.handlers.ByKey(key); ok {
		validation.Panicf(
			`%s can not use the application key "%s", because it is already used by %s`,
			c.config.ReflectType(),
			key,
			h.ReflectType(),
		)
	}

	configureIdentity(
		c.config.ReflectType(),
		&c.config.ident,
		name,
		key,
	)
}

func (c *applicationConfigurer) RegisterAggregate(h dogma.AggregateMessageHandler, _ ...dogma.RegisterAggregateOption) {
	c.registerIfConfigured(fromAggregate(h))
}

func (c *applicationConfigurer) RegisterProcess(h dogma.ProcessMessageHandler, _ ...dogma.RegisterProcessOption) {
	c.registerIfConfiguredX(fromProcess(h))
}

func (c *applicationConfigurer) RegisterIntegration(h dogma.IntegrationMessageHandler, _ ...dogma.RegisterIntegrationOption) {
	c.registerIfConfiguredX(fromIntegration(h))
}

func (c *applicationConfigurer) RegisterProjection(h dogma.ProjectionMessageHandler, _ ...dogma.RegisterProjectionOption) {
	c.registerIfConfiguredX(fromProjection(h))
}

func (c *applicationConfigurer) registerIfConfigured(
	h interface {
		RichHandler
		isConfigured() bool
		mustValidate()
	},
) {
	if h.IsDisabled() && !h.isConfigured() {
		return
	}

	h.mustValidate()

	c.guardAgainstConflictingIdentities(h)
	c.guardAgainstConflictingRoles(h)
	c.guardAgainstConflictingRoutes(h)

	if c.config.handlers == nil {
		c.config.handlers = RichHandlerSet{}
	}
	c.config.handlers.Add(h)

	types := h.MessageTypes()

	for mt, r := range types.Produced {
		if c.config.types.Produced == nil {
			c.config.types.Produced = message.TypeRoles{}
		}
		c.config.types.Produced.Add(mt, r)
	}

	for mt, r := range types.Consumed {
		if c.config.types.Consumed == nil {
			c.config.types.Consumed = message.TypeRoles{}
		}
		c.config.types.Consumed.Add(mt, r)
	}
}

// register adds a handler configuration to the application.
// deprecated: use registerIfConfigured() instead.
func (c *applicationConfigurer) registerIfConfiguredX(
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

	if c.config.handlers == nil {
		c.config.handlers = RichHandlerSet{}
	}
	c.config.handlers.Add(h)

	types := h.MessageTypes()

	for mt, r := range types.Produced {
		if c.config.types.Produced == nil {
			c.config.types.Produced = message.TypeRoles{}
		}
		c.config.types.Produced.Add(mt, r)
	}

	for mt, r := range types.Consumed {
		if c.config.types.Consumed == nil {
			c.config.types.Consumed = message.TypeRoles{}
		}
		c.config.types.Consumed.Add(mt, r)
	}
}

// guardAgainstConflictingIdentities panics if h's identity conflicts with the
// application or any other handlers.
func (c *applicationConfigurer) guardAgainstConflictingIdentities(h RichHandler) {
	appIdent := c.config.Identity()
	handlerIdent := h.Identity()

	if handlerIdent.Key == appIdent.Key {
		validation.Panicf(
			`%s can not use the handler key "%s", because it is already used by %s`,
			h.ReflectType(),
			handlerIdent.Key,
			c.config.ReflectType(),
		)
	}

	if x, ok := c.config.handlers.ByName(handlerIdent.Name); ok {
		validation.Panicf(
			`%s can not use the handler name "%s", because it is already used by %s`,
			h.ReflectType(),
			handlerIdent.Name,
			x.ReflectType(),
		)
	}

	if x, ok := c.config.handlers.ByKey(handlerIdent.Key); ok {
		validation.Panicf(
			`%s can not use the handler key "%s", because it is already used by %s`,
			h.ReflectType(),
			handlerIdent.Key,
			x.ReflectType(),
		)
	}
}

// guardAgainstConflictingRoles panics if h configures any messages in roles
// contrary to the way they are configured by any other handler.
func (c *applicationConfigurer) guardAgainstConflictingRoles(h RichHandler) {
	for mt, r := range h.MessageTypes().All() {
		xr, ok := c.config.types.RoleOf(mt)

		if !ok || xr == r {
			continue
		}

		// we know there's a conflict, now we just need to find a handler that
		// refers to this message type as some other role.
		xh, _ := c.config.handlers.Find(func(h RichHandler) bool {
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

		for _, x := range c.config.handlers.ConsumersOf(mt) {
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

		for _, x := range c.config.handlers.ProducersOf(mt) {
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
