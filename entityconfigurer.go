package configkit

// handlerConfigurer is a partial implementation of the configurer interfaces for
// all of the Dogma entity types.
//
// - dogma.ApplicationConfigurer
// - dogma.AggregateConfigurer
// - dogma.ProcessConfigurer
// - dogma.IntegrationConfigurer
// - dogma.ProjectionConfigurer
type entityConfigurer struct {
	// entity is the target entity to populate with the configuration values.
	entity *entity
}

// Identity sets the handler's identity.
func (c *entityConfigurer) Identity(name string, key string) {
	if !c.entity.ident.IsZero() {
		Panicf(
			"%s is configured with multiple identities (%s and %s/%s), Identity() must be called exactly once within Configure()",
			c.entity.rt,
			c.entity.ident,
			name,
			key,
		)
	}

	var err error
	c.entity.ident, err = NewIdentity(name, key)

	if err != nil {
		Panicf(
			"%s is configured with an invalid identity, %s",
			c.entity.rt,
			err,
		)
	}
}

// validate panics if the configuration is invalid.
func (c *entityConfigurer) validate() {
	if c.entity.ident.IsZero() {
		Panicf(
			"%s is configured without an identity, Identity() must be called exactly once within Configure()",
			c.entity.rt,
		)
	}
}

// displayName returns a human-readable string used to refer to the entity in
// error messages.
func (c *entityConfigurer) displayName() string {
	s := c.entity.rt.String()

	if !c.entity.ident.IsZero() {
		s += " (" + c.entity.ident.Name + ")"
	}

	return s
}
