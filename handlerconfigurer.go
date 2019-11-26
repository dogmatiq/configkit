package configkit

import (
	"strings"

	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

// handlerConfigurer is an implementation of the the configurer interfaces for
// each of the Dogma handler types.
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
			"%s.Configure() has already called %s.Identity(%#v, %#v)",
			c.target.rt.String(),
			c.interfaceName,
			c.target.ident.Name,
			c.target.ident.Key,
		)
	}

	var err error
	c.target.ident, err = NewIdentity(name, key)

	if err != nil {
		Panicf(
			"%s.Configure() called %s.Identity() with an %s",
			c.target.rt.String(),
			c.interfaceName,
			err.Error(),
		)
	}
}

func (c *handlerConfigurer) ConsumesCommandType(m dogma.Message) {
	c.consumes(m, message.CommandRole)
}

func (c *handlerConfigurer) ProducesEventType(m dogma.Message) {
	c.produces(m, message.EventRole)
}

func (c *handlerConfigurer) consumes(m dogma.Message, r message.Role) {
	if c.target.names.Roles == nil {
		c.target.names.Roles = message.NameRoles{}
		c.target.types.Roles = message.TypeRoles{}
	}

	if c.target.names.Consumed == nil {
		c.target.names.Consumed = message.NameSet{}
		c.target.types.Consumed = message.TypeSet{}
	}

	mt := message.TypeOf(m)
	n := mt.Name()

	if x, ok := c.target.names.Roles[n]; ok {
		if x == r {
			Panicf(
				"%s.Configure() has already called %s.Consumes%sType(%s)",
				c.target.rt.String(),
				c.interfaceName,
				strings.Title(r.String()),
				mt,
			)
		}
	}

	c.target.names.Roles[n] = r
	c.target.names.Consumed.Add(n)
	c.target.types.Roles[mt] = r
	c.target.types.Consumed.Add(mt)
}

func (c *handlerConfigurer) produces(m dogma.Message, r message.Role) {
	if c.target.names.Roles == nil {
		c.target.names.Roles = message.NameRoles{}
		c.target.types.Roles = message.TypeRoles{}
	}

	if c.target.names.Produced == nil {
		c.target.names.Produced = message.NameSet{}
		c.target.types.Produced = message.TypeSet{}
	}

	mt := message.TypeOf(m)
	n := mt.Name()

	if x, ok := c.target.names.Roles[n]; ok {
		if x == r {
			Panicf(
				"%s.Configure() has already called %s.Produces%sType(%s)",
				c.target.rt.String(),
				c.interfaceName,
				strings.Title(r.String()),
				mt,
			)
		}
	}

	c.target.names.Roles[n] = r
	c.target.names.Produced.Add(n)
	c.target.types.Roles[mt] = r
	c.target.types.Produced.Add(mt)
}

// validate panics if the configuration is invalid.
func (c *handlerConfigurer) validate() {
	if c.target.ident.IsZero() {
		Panicf(
			"%s.Configure() did not call %s.Identity()",
			c.target.rt.String(),
			c.interfaceName,
		)
	}
}

func (c *handlerConfigurer) mustConsume(r message.Role) {
	for mt := range c.target.names.Consumed {
		if r == c.target.names.Roles[mt] {
			return
		}
	}

	Panicf(
		"%s.Configure() did not call %s.Consumes%sType()",
		c.target.rt.String(),
		c.interfaceName,
		strings.Title(r.String()),
	)
}

func (c *handlerConfigurer) mustProduce(r message.Role) {
	for mt := range c.target.names.Produced {
		if r == c.target.names.Roles[mt] {
			return
		}
	}

	Panicf(
		"%s.Configure() did not call %s.Produces%sType()",
		c.target.rt.String(),
		c.interfaceName,
		strings.Title(r.String()),
	)
}
