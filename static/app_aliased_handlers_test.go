package static_test

import (
	"os"

	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("registering handlers using a type alias", func() {
	goDebugBefore := os.Getenv("GODEBUG")

	BeforeEach(func() {
		// Set the GODEBUG environment variable to enable type alias support.
		//
		// TODO: Remove this setting once the package is migrated to Go 1.23
		// that generates `types.Alias` types by default.
		os.Setenv("GODEBUG", "gotypesalias=1")
	})

	AfterEach(func() {
		os.Setenv("GODEBUG", goDebugBefore)
	})

	It("it reports the aliased type name", func() {
		apps := FromDir("testdata/apps/aliased-handlers")
		Expect(apps).To(HaveLen(1))
		Expect(apps[0].Handlers().Aggregates()).To(HaveLen(1))
		Expect(apps[0].Handlers().Processes()).To(HaveLen(1))
		Expect(apps[0].Handlers().Projections()).To(HaveLen(1))
		Expect(apps[0].Handlers().Integrations()).To(HaveLen(1))

		aggregate := apps[0].Handlers().Aggregates()[0]
		Expect(aggregate.Identity()).To(
			Equal(
				configkit.Identity{
					Name: "<aggregate>",
					Key:  "92623de9-c9cf-42f3-8338-33c50eeb06fb",
				},
			),
		)
		Expect(aggregate.TypeName()).To(
			Equal(
				"github.com/dogmatiq/configkit/static/testdata/apps/aliased-handlers.AggregateHandlerAlias",
			),
		)

		process := apps[0].Handlers().Processes()[0]
		Expect(process.Identity()).To(
			Equal(
				configkit.Identity{
					Name: "<process>",
					Key:  "ad9d6955-893a-4d8d-a26e-e25886b113b2",
				},
			),
		)
		Expect(process.TypeName()).To(
			Equal(
				"github.com/dogmatiq/configkit/static/testdata/apps/aliased-handlers.ProcessHandlerAlias",
			),
		)

		projection := apps[0].Handlers().Projections()[0]
		Expect(projection.Identity()).To(
			Equal(
				configkit.Identity{
					Name: "<projection>",
					Key:  "d012b7ed-3c4b-44db-9276-7bbc90fb54fd",
				},
			),
		)
		Expect(projection.TypeName()).To(
			Equal(
				"github.com/dogmatiq/configkit/static/testdata/apps/aliased-handlers.ProjectionHandlerAlias",
			),
		)

		integration := apps[0].Handlers().Integrations()[0]
		Expect(integration.Identity()).To(
			Equal(
				configkit.Identity{
					Name: "<integration>",
					Key:  "4d8cd3f5-21dc-475b-a8dc-80138adde3f2",
				},
			),
		)
		Expect(integration.TypeName()).To(
			Equal(
				"github.com/dogmatiq/configkit/static/testdata/apps/aliased-handlers.IntegrationHandlerAlias",
			),
		)
	})
})
