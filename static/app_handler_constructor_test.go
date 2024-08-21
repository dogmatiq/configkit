package static_test

import (
	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("handlers registered via constructor function calls", func() {
	It("builds the configuration based on the return type of the function", func() {
		apps := FromDir("testdata/apps/handler-via-constructor")
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
					Key:  "ef16c9d1-d7b6-4c99-a0e7-a59218e544fc",
				},
			),
		)
		Expect(aggregate.TypeName()).To(
			Equal(
				"github.com/dogmatiq/configkit/static/testdata/apps/handler-via-constructor.AggregateHandler",
			),
		)

		process := apps[0].Handlers().Processes()[0]
		Expect(process.Identity()).To(
			Equal(
				configkit.Identity{
					Name: "<process>",
					Key:  "5e839b73-170b-42c0-bf41-8feee4b5a583",
				},
			),
		)
		Expect(process.TypeName()).To(
			Equal(
				"github.com/dogmatiq/configkit/static/testdata/apps/handler-via-constructor.ProcessHandler",
			),
		)
		Expect(process.HandlerType()).To(Equal(configkit.ProcessHandlerType))

		projection := apps[0].Handlers().Projections()[0]
		Expect(projection.Identity()).To(
			Equal(
				configkit.Identity{
					Name: "<projection>",
					Key:  "823e61d3-ace1-469d-b0a6-778e84c0a508",
				},
			),
		)
		Expect(projection.TypeName()).To(
			Equal(
				"github.com/dogmatiq/configkit/static/testdata/apps/handler-via-constructor.ProjectionHandler",
			),
		)
		Expect(projection.HandlerType()).To(Equal(configkit.ProjectionHandlerType))

		integration := apps[0].Handlers().Integrations()[0]
		Expect(integration.Identity()).To(
			Equal(
				configkit.Identity{
					Name: "<integration>",
					Key:  "099b5b8d-9e04-422f-bcc3-bb0d451158c7",
				},
			),
		)
		Expect(integration.TypeName()).To(
			Equal(
				"github.com/dogmatiq/configkit/static/testdata/apps/handler-via-constructor.IntegrationHandler",
			),
		)
		Expect(integration.HandlerType()).To(Equal(configkit.IntegrationHandlerType))
	})
})
