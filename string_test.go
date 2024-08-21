package configkit_test

import (
	"strings"

	. "github.com/dogmatiq/configkit"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func ToString()", func() {
	var cfg Application

	BeforeEach(func() {
		app := &ApplicationStub{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app>", appKey)

				c.RegisterAggregate(&AggregateMessageHandlerStub{
					ConfigureFunc: func(c dogma.AggregateConfigurer) {
						c.Identity("<aggregate>", aggregateKey)
						c.Routes(
							dogma.HandlesCommand[CommandStub[TypeA]](),
							dogma.RecordsEvent[EventStub[TypeA]](),
						)
					},
				})

				c.RegisterProcess(&ProcessMessageHandlerStub{
					ConfigureFunc: func(c dogma.ProcessConfigurer) {
						c.Identity("<process>", processKey)
						c.Routes(
							dogma.HandlesEvent[EventStub[TypeA]](),
							dogma.ExecutesCommand[CommandStub[TypeA]](),
							dogma.SchedulesTimeout[TimeoutStub[TypeA]](),
						)
					},
				})

				c.RegisterIntegration(&IntegrationMessageHandlerStub{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Identity("<integration>", integrationKey)
						c.Routes(
							dogma.HandlesCommand[CommandStub[TypeB]](),
							dogma.RecordsEvent[EventStub[TypeB]](),
						)
					},
				})

				c.RegisterProjection(&ProjectionMessageHandlerStub{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Identity("<projection>", projectionKey)
						c.Routes(
							dogma.HandlesEvent[EventStub[TypeA]](),
							dogma.HandlesEvent[EventStub[TypeB]](),
						)
						c.Disable()
					},
				})
			},
		}

		cfg = FromApplication(app)
	})

	It("returns a human readable string representation", func() {
		expected := "application <app> (59a82a24-a181-41e8-9b93-17a6ce86956e) *github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub\n"
		expected += "\n"
		expected += "    - aggregate <aggregate> (14769f7f-87fe-48dd-916e-5bcab6ba6aca) *github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub\n"
		expected += "        handles github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]?\n"
		expected += "        records github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]!\n"
		expected += "\n"
		expected += "    - integration <integration> (e28f056e-e5a0-4ee7-aaf1-1d1fe02fb6e3) *github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub\n"
		expected += "        handles github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeB]?\n"
		expected += "        records github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeB]!\n"
		expected += "\n"
		expected += "    - process <process> (bea52cf4-e403-4b18-819d-88ade7836308) *github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub\n"
		expected += "        handles github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]!\n"
		expected += "        executes github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]?\n"
		expected += "        schedules github.com/dogmatiq/enginekit/enginetest/stubs.TimeoutStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]@\n"
		expected += "\n"
		expected += "    - projection <projection> (70fdf7fa-4b24-448d-bd29-7ecc71d18c56) *github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub [disabled]\n"
		expected += "        handles github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]!\n"
		expected += "        handles github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeB]!\n"

		s := ToString(cfg)

		// compare as slices, as the failure output for slices is better than
		// for multiline strings :(
		Expect(
			strings.Split(s, "\n"),
		).To(
			Equal(strings.Split(expected, "\n")),
		)
	})
})
