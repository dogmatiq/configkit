package configkit

import (
	"strings"

	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

// handlerConfigurer is an implementation of the the configurer interfaces for
// all of the Dogma handler types.
//
// - dogma.AggregateConfigurer
// - dogma.ProcessConfigurer
// - dogma.IntegrationConfigurer
// - dogma.ProjectionConfigurer
type handlerConfigurer struct {
	// interfaceName is the name of the Dogma interface that the configurer is being used as.
	// It is used to construct error messages.
	interfaceName string

	// target is the handler to populate with the configuration values.
	target *handler
}

func (c *handlerConfigurer) Identity(name string, key string) {
	if !c.target.ident.IsZero() {
		Panicf(
			"%s is configured with multiple identities (%s and %s/%s), Identity() must be called exactly once within Configure()",
			c.target.rt.String(),
			c.target.ident,
			name,
			key,
		)
	}

	var err error
	c.target.ident, err = NewIdentity(name, key)

	if err != nil {
		Panicf(
			"%s is configured with an invalid identity, %s",
			c.target.rt.String(),
			err,
		)
	}
}

func (c *handlerConfigurer) ConsumesCommandType(m dogma.Message) {
	c.consumes(m, message.CommandRole)
}

func (c *handlerConfigurer) ConsumesEventType(m dogma.Message) {
	c.consumes(m, message.EventRole)
}

func (c *handlerConfigurer) ProducesCommandType(m dogma.Message) {
	c.produces(m, message.CommandRole)
}

func (c *handlerConfigurer) ProducesEventType(m dogma.Message) {
	c.produces(m, message.EventRole)
}

func (c *handlerConfigurer) SchedulesTimeoutType(m dogma.Message) {
	c.produces(m, message.TimeoutRole)
	c.consumes(m, message.TimeoutRole)
}

func (c *handlerConfigurer) consumes(m dogma.Message, r message.Role) {
	mt := message.TypeOf(m)
	c.guardAgainstRoleMismatch(mt, r)

	if c.target.types.Consumed.Has(mt) {
		Panicf(
			"%s is configured to consume %s more than once, should this refer to different message types?",
			c.target.rt.String(),
			mt,
		)
	}

	if c.target.names.Roles == nil {
		c.target.names.Roles = message.NameRoles{}
		c.target.types.Roles = message.TypeRoles{}
	}

	if c.target.names.Consumed == nil {
		c.target.names.Consumed = message.NameSet{}
		c.target.types.Consumed = message.TypeSet{}
	}

	n := mt.Name()
	c.target.names.Roles[n] = r
	c.target.names.Consumed.Add(n)
	c.target.types.Roles[mt] = r
	c.target.types.Consumed.Add(mt)
}

func (c *handlerConfigurer) produces(m dogma.Message, r message.Role) {
	mt := message.TypeOf(m)
	c.guardAgainstRoleMismatch(mt, r)

	if c.target.types.Produced.Has(mt) {
		Panicf(
			"%s is configured to produce %s more than once, should this refer to different message types?",
			c.target.rt.String(),
			mt,
		)
	}
	if c.target.names.Roles == nil {
		c.target.names.Roles = message.NameRoles{}
		c.target.types.Roles = message.TypeRoles{}
	}

	if c.target.names.Produced == nil {
		c.target.names.Produced = message.NameSet{}
		c.target.types.Produced = message.TypeSet{}
	}

	n := mt.Name()
	c.target.names.Roles[n] = r
	c.target.names.Produced.Add(n)
	c.target.types.Roles[mt] = r
	c.target.types.Produced.Add(mt)
}

func (c *handlerConfigurer) guardAgainstRoleMismatch(mt message.Type, r message.Role) {
	x, ok := c.target.types.Roles[mt]

	if !ok || x == r {
		return
	}

	Panicf(
		"%s is configured to use %s as both a %s and a %s",
		c.target.rt.String(),
		mt,
		x,
		r,
	)
}

// validate panics if the configuration is invalid.
func (c *handlerConfigurer) validate() {
	if c.target.ident.IsZero() {
		Panicf(
			"%s is configured without an identity, Identity() must be called exactly once within Configure()",
			c.target.rt.String(),
		)
	}
}

// mustConsume panics if the handler does not consume any messages of the given role.
func (c *handlerConfigurer) mustConsume(r message.Role) {
	for mt := range c.target.names.Consumed {
		if r == c.target.names.Roles[mt] {
			return
		}
	}

	Panicf(
		`%s is not configured to consume any %ss, Consumes%sType() must be called at least once within Configure()`,
		c.target.rt.String(),
		r.String(),
		strings.Title(r.String()),
	)
}

// mustProduce panics if the handler does not produce any messages of the given role.
func (c *handlerConfigurer) mustProduce(r message.Role) {
	for mt := range c.target.names.Produced {
		if r == c.target.names.Roles[mt] {
			return
		}
	}

	Panicf(
		`%s is not configured to produce any %ss, Produces%sType() must be called at least once within Configure()`,
		c.target.rt.String(),
		r.String(),
		strings.Title(r.String()),
	)
}