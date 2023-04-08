package configkit

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dogmatiq/configkit/internal/validation"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

// handlerConfigurer is a partial implementation of the configurer interfaces for
// all of the Dogma handler types.
//
//   - [dogma.AggregateConfigurer]
//   - [dogma.ProcessConfigurer]
//   - [dogma.IntegrationConfigurer]
//   - [dogma.ProjectionConfigurer]
type handlerConfigurer struct {
	entityConfigurer
}

func (c *handlerConfigurer) route(r dogma.Route) {
	switch r := r.(type) {
	case dogma.HandlesCommandRoute:
		c.consumes(r.Type, message.CommandRole, "consume")
	case dogma.RecordsEventRoute:
		c.produces(r.Type, message.EventRole, "produce")
	case dogma.HandlesEventRoute:
		c.consumes(r.Type, message.EventRole, "consume")
	case dogma.ExecutesCommandRoute:
		c.produces(r.Type, message.CommandRole, "produce")
	case dogma.SchedulesTimeoutRoute:
		c.consumes(r.Type, message.TimeoutRole, "schedule")
		c.produces(r.Type, message.TimeoutRole, "schedule")
	default:
		panic(fmt.Sprintf("unsupported route type: %T", r))
	}
}

func (c *handlerConfigurer) ConsumesCommandType(m dogma.Message) {
	c.route(dogma.HandlesCommandRoute{Type: reflect.TypeOf(m)})
}

func (c *handlerConfigurer) ConsumesEventType(m dogma.Message) {
	c.route(dogma.HandlesEventRoute{Type: reflect.TypeOf(m)})
}

func (c *handlerConfigurer) ProducesCommandType(m dogma.Message) {
	c.route(dogma.ExecutesCommandRoute{Type: reflect.TypeOf(m)})
}

func (c *handlerConfigurer) ProducesEventType(m dogma.Message) {
	c.route(dogma.RecordsEventRoute{Type: reflect.TypeOf(m)})
}

func (c *handlerConfigurer) SchedulesTimeoutType(m dogma.Message) {
	c.route(dogma.SchedulesTimeoutRoute{Type: reflect.TypeOf(m)})
}

func (c *handlerConfigurer) consumes(t reflect.Type, r message.Role, verb string) {
	mt := message.TypeFromReflect(t)
	c.guardAgainstConflictingRoles(mt, r)

	if c.entity.types.Consumed.Has(mt) {
		validation.Panicf(
			"%s is configured to %s the %s %s more than once, should this refer to different message types?",
			c.displayName(),
			verb,
			mt,
			r,
		)
	}

	if c.entity.names.Consumed == nil {
		c.entity.names.Consumed = message.NameRoles{}
		c.entity.types.Consumed = message.TypeRoles{}
	}

	n := mt.Name()
	c.entity.names.Consumed.Add(n, r)
	c.entity.types.Consumed.Add(mt, r)
}

func (c *handlerConfigurer) produces(t reflect.Type, r message.Role, verb string) {
	mt := message.TypeFromReflect(t)
	c.guardAgainstConflictingRoles(mt, r)

	if c.entity.types.Produced.Has(mt) {
		validation.Panicf(
			"%s is configured to %s the %s %s more than once, should this refer to different message types?",
			c.displayName(),
			verb,
			mt,
			r,
		)
	}

	if c.entity.names.Produced == nil {
		c.entity.names.Produced = message.NameRoles{}
		c.entity.types.Produced = message.TypeRoles{}
	}

	n := mt.Name()
	c.entity.names.Produced.Add(n, r)
	c.entity.types.Produced.Add(mt, r)
}

// guardAgainstConflictingRoles panics if mt is already used in some role other than r.
func (c *handlerConfigurer) guardAgainstConflictingRoles(mt message.Type, r message.Role) {
	x, ok := c.entity.types.RoleOf(mt)

	if !ok || x == r {
		return
	}

	validation.Panicf(
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
		if x, ok := c.entity.names.RoleOf(mt); ok {
			if x == r {
				return
			}
		}
	}

	validation.Panicf(
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
		if x, ok := c.entity.names.RoleOf(mt); ok {
			if x == r {
				return
			}
		}
	}

	validation.Panicf(
		`%s is not configured to produce any %ss, Produces%sType() must be called at least once within Configure()`,
		c.displayName(),
		r,
		strings.Title(r.String()),
	)
}
