package configkit

import (
	"strings"

	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

// handlerConfigurer is an implementation of the configurer interfaces for
// all of the Dogma handler types.
//
// - dogma.AggregateConfigurer
// - dogma.ProcessConfigurer
// - dogma.IntegrationConfigurer
// - dogma.ProjectionConfigurer
type handlerConfigurer struct {
	entityConfigurer
}

// ConsumesCommandTypes marks the handler as a consumer of command messages of
// the same type as m.
func (c *handlerConfigurer) ConsumesCommandType(m dogma.Message) {
	c.consumes(m, message.CommandRole, "consume")
}

// ConsumesEventType marks the handler as a consumer of event messages of the
// same type as m.
func (c *handlerConfigurer) ConsumesEventType(m dogma.Message) {
	c.consumes(m, message.EventRole, "consume")
}

// ProducesCommandType marks the handler as a producer of command messages of
// the same type as m.
func (c *handlerConfigurer) ProducesCommandType(m dogma.Message) {
	c.produces(m, message.CommandRole, "produce")
}

// ProducesEventType marks the handler as a producer of event messages of the
// same type as m.
func (c *handlerConfigurer) ProducesEventType(m dogma.Message) {
	c.produces(m, message.EventRole, "produce")
}

// SchedulesTimeoutType marks the handler as a scheduler of timeout messages of
// the same type as m.
func (c *handlerConfigurer) SchedulesTimeoutType(m dogma.Message) {
	c.consumes(m, message.TimeoutRole, "schedule")
	c.produces(m, message.TimeoutRole, "schedule")
}

// consumes marks the handler as a consumer of messages of the same type as m.
func (c *handlerConfigurer) consumes(m dogma.Message, r message.Role, verb string) {
	mt := message.TypeOf(m)
	c.guardAgainstRoleMismatch(mt, r)

	if c.entity.types.Consumed.Has(mt) {
		Panicf(
			"%s is configured to %s the %s %s more than once, should this refer to different message types?",
			c.displayName(),
			verb,
			mt,
			r,
		)
	}

	if c.entity.names.Roles == nil {
		c.entity.names.Roles = message.NameRoles{}
		c.entity.types.Roles = message.TypeRoles{}
	}

	if c.entity.names.Consumed == nil {
		c.entity.names.Consumed = message.NameRoles{}
		c.entity.types.Consumed = message.TypeRoles{}
	}

	n := mt.Name()
	c.entity.names.Roles.Add(n, r)
	c.entity.names.Consumed.Add(n, r)
	c.entity.types.Roles.Add(mt, r)
	c.entity.types.Consumed.Add(mt, r)
}

// produces marks the handler as a consumer of messages of the same type as m.
func (c *handlerConfigurer) produces(m dogma.Message, r message.Role, verb string) {
	mt := message.TypeOf(m)
	c.guardAgainstRoleMismatch(mt, r)

	if c.entity.types.Produced.Has(mt) {
		Panicf(
			"%s is configured to %s the %s %s more than once, should this refer to different message types?",
			c.displayName(),
			verb,
			mt,
			r,
		)
	}
	if c.entity.names.Roles == nil {
		c.entity.names.Roles = message.NameRoles{}
		c.entity.types.Roles = message.TypeRoles{}
	}

	if c.entity.names.Produced == nil {
		c.entity.names.Produced = message.NameRoles{}
		c.entity.types.Produced = message.TypeRoles{}
	}

	n := mt.Name()
	c.entity.names.Roles.Add(n, r)
	c.entity.names.Produced.Add(n, r)
	c.entity.types.Roles.Add(mt, r)
	c.entity.types.Produced.Add(mt, r)
}

// guardAgainstRoleMismatch panics if mt is already used in some role other than r.
func (c *handlerConfigurer) guardAgainstRoleMismatch(mt message.Type, r message.Role) {
	x, ok := c.entity.types.Roles[mt]

	if !ok || x == r {
		return
	}

	Panicf(
		"%s is configured to use %s as both a %s and a %s",
		c.displayName(),
		mt,
		x,
		r,
	)
}

// mustConsume panics if the handler does not consume any messages of the given role.
func (c *handlerConfigurer) mustConsume(r message.Role) {
	for mt := range c.entity.names.Consumed {
		if r == c.entity.names.Roles[mt] {
			return
		}
	}

	Panicf(
		`%s (%s) is not configured to consume any %ss, Consumes%sType() must be called at least once within Configure()`,
		c.entity.rt,
		c.entity.ident.Name,
		r,
		strings.Title(r.String()),
	)
}

// mustProduce panics if the handler does not produce any messages of the given role.
func (c *handlerConfigurer) mustProduce(r message.Role) {
	for mt := range c.entity.names.Produced {
		if r == c.entity.names.Roles[mt] {
			return
		}
	}

	Panicf(
		`%s is not configured to produce any %ss, Produces%sType() must be called at least once within Configure()`,
		c.displayName(),
		r,
		strings.Title(r.String()),
	)
}

func (c *handlerConfigurer) displayName() string {
	s := c.entity.rt.String()

	if !c.entity.ident.IsZero() {
		s += " (" + c.entity.ident.Name + ")"
	}

	return s
}
