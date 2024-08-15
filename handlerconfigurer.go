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
	handler *handlerEntity
}

func (c *handlerConfigurer) Disable(...dogma.DisableOption) {
	c.handler.isDisabled = true
}

func (c *handlerConfigurer) route(r dogma.Route) {
	c.configured = true

	switch r := r.(type) {
	case dogma.HandlesCommandRoute:
		c.consumes(r.Type, message.CommandRole, "HandlesCommand")
	case dogma.RecordsEventRoute:
		c.produces(r.Type, message.EventRole, "RecordsEvent")
	case dogma.HandlesEventRoute:
		c.consumes(r.Type, message.EventRole, "HandlesEvent")
	case dogma.ExecutesCommandRoute:
		c.produces(r.Type, message.CommandRole, "ExecutesCommand")
	case dogma.SchedulesTimeoutRoute:
		c.consumes(r.Type, message.TimeoutRole, "SchedulesTimeout")
		c.produces(r.Type, message.TimeoutRole, "SchedulesTimeout")
	default:
		panic(fmt.Sprintf("unsupported route type: %T", r))
	}
}

func (c *handlerConfigurer) consumes(t reflect.Type, r message.Role, route string) {
	mt := message.TypeFromReflect(t)
	c.guardAgainstConflictingRoles(mt, r)

	if c.entity.types.Consumed.Has(mt) {
		validation.Panicf(
			"%s is configured with multiple %s() routes for %s, should these refer to different message types?",
			c.displayName(),
			route,
			mt,
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

func (c *handlerConfigurer) produces(t reflect.Type, r message.Role, route string) {
	mt := message.TypeFromReflect(t)
	c.guardAgainstConflictingRoles(mt, r)

	if c.entity.types.Produced.Has(mt) {
		validation.Panicf(
			"%s is configured with multiple %s() routes for %s, should these refer to different message types?",
			c.displayName(),
			route,
			mt,
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

// mustConsume panics if the handler is not configured to handle any messages of
// the given role.
func (c *handlerConfigurer) mustConsume(r message.Role) {
	for mt := range c.entity.names.Consumed {
		if x, ok := c.entity.names.RoleOf(mt); ok {
			if x == r {
				return
			}
		}
	}

	validation.Panicf(
		`%s (%s) is not configured to handle any %ss, at least one Handles%s() route must be added within Configure()`,
		c.entity.rt,
		c.entity.ident.Name,
		r,
		strings.Title(r.String()),
	)
}

// mustProduce panics if the handler is not configured to produce any messages
// of the given role.
func (c *handlerConfigurer) mustProduce(r message.Role) {
	for mt := range c.entity.names.Produced {
		if x, ok := c.entity.names.RoleOf(mt); ok {
			if x == r {
				return
			}
		}
	}

	verb := ""
	route := ""
	switch r {
	case message.CommandRole:
		verb = "execute"
		route = "ExecutesCommand"
	case message.EventRole:
		verb = "record"
		route = "RecordsEvent"
	case message.TimeoutRole:
		verb = "schedule"
		route = "SchedulesTimeout"
	}

	validation.Panicf(
		`%s is not configured to %s any %ss, at least one %s() route must be added within Configure()`,
		c.displayName(),
		verb,
		r,
		route,
	)
}
