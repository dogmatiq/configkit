package configkit_test

import . "github.com/dogmatiq/configkit"

// Following is a set of static assertions that confirms the "inheritance"
// relationships between the various interfaces are as expected.
//
// The expressions can be read as TypeOnLeft (is implemented by) TypeOnRight.
var (
	// Entity
	_ Entity = PortableEntity(nil)
	_ Entity = RichEntity(nil)

	// Application
	_ Entity         = Application(nil)
	_ PortableEntity = PortableApplication(nil)
	_ Application    = PortableApplication(nil)
	_ RichEntity     = RichApplication(nil)
	_ Application    = RichApplication(nil)

	// Handler
	_ Entity         = Handler(nil)
	_ PortableEntity = PortableHandler(nil)
	_ Handler        = PortableHandler(nil)
	_ RichEntity     = RichHandler(nil)
	_ Handler        = RichHandler(nil)

	// Aggregate
	_ Handler         = Aggregate(nil)
	_ PortableHandler = PortableAggregate(nil)
	_ Aggregate       = PortableAggregate(nil)
	_ RichHandler     = RichAggregate(nil)
	_ Aggregate       = RichAggregate(nil)

	// Process
	_ Handler         = Process(nil)
	_ PortableHandler = PortableProcess(nil)
	_ Process         = PortableProcess(nil)
	_ RichHandler     = RichProcess(nil)
	_ Process         = RichProcess(nil)

	// Integration
	_ Handler         = Integration(nil)
	_ PortableHandler = PortableIntegration(nil)
	_ Integration     = PortableIntegration(nil)
	_ RichHandler     = RichIntegration(nil)
	_ Integration     = RichIntegration(nil)

	// Projection
	_ Handler         = Projection(nil)
	_ PortableHandler = PortableProjection(nil)
	_ Projection      = PortableProjection(nil)
	_ RichHandler     = RichProjection(nil)
	_ Projection      = RichProjection(nil)
)
