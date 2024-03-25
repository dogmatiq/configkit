package configkit_test

import (
	"strings"

	. "github.com/dogmatiq/configkit"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures" // can't dot-import due to conflicts
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func ToString()", func() {
	var cfg Application

	BeforeEach(func() {
		app := &fixtures.Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app>", appKey)

				c.RegisterAggregate(&fixtures.AggregateMessageHandler{
					ConfigureFunc: func(c dogma.AggregateConfigurer) {
						c.Identity("<aggregate>", aggregateKey)
						c.Routes(
							dogma.HandlesCommand[fixtures.MessageC](),
							dogma.RecordsEvent[fixtures.MessageE](),
						)
					},
				})

				c.RegisterProcess(&fixtures.ProcessMessageHandler{
					ConfigureFunc: func(c dogma.ProcessConfigurer) {
						c.Identity("<process>", processKey)
						c.Routes(
							dogma.HandlesEvent[fixtures.MessageE](),
							dogma.ExecutesCommand[fixtures.MessageC](),
							dogma.SchedulesTimeout[fixtures.MessageT](),
						)
					},
				})

				c.RegisterIntegration(&fixtures.IntegrationMessageHandler{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Identity("<integration>", integrationKey)
						c.Routes(
							dogma.HandlesCommand[fixtures.MessageI](),
							dogma.RecordsEvent[fixtures.MessageJ](),
						)
					},
				})

				c.RegisterProjection(&fixtures.ProjectionMessageHandler{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Identity("<projection>", projectionKey)
						c.Routes(
							dogma.HandlesEvent[fixtures.MessageE](),
							dogma.HandlesEvent[fixtures.MessageJ](),
						)
					},
				})
			},
		}

		cfg = FromApplication(app)
	})

	It("returns a human readable string representation", func() {
		expected := "application <app> (59a82a24-a181-41e8-9b93-17a6ce86956e) *github.com/dogmatiq/dogma/fixtures.Application\n"
		expected += "\n"
		expected += "    - aggregate <aggregate> (14769f7f-87fe-48dd-916e-5bcab6ba6aca) *github.com/dogmatiq/dogma/fixtures.AggregateMessageHandler\n"
		expected += "        handles github.com/dogmatiq/dogma/fixtures.MessageC?\n"
		expected += "        records github.com/dogmatiq/dogma/fixtures.MessageE!\n"
		expected += "\n"
		expected += "    - integration <integration> (e28f056e-e5a0-4ee7-aaf1-1d1fe02fb6e3) *github.com/dogmatiq/dogma/fixtures.IntegrationMessageHandler\n"
		expected += "        handles github.com/dogmatiq/dogma/fixtures.MessageI?\n"
		expected += "        records github.com/dogmatiq/dogma/fixtures.MessageJ!\n"
		expected += "\n"
		expected += "    - process <process> (bea52cf4-e403-4b18-819d-88ade7836308) *github.com/dogmatiq/dogma/fixtures.ProcessMessageHandler\n"
		expected += "        handles github.com/dogmatiq/dogma/fixtures.MessageE!\n"
		expected += "        executes github.com/dogmatiq/dogma/fixtures.MessageC?\n"
		expected += "        schedules github.com/dogmatiq/dogma/fixtures.MessageT@\n"
		expected += "\n"
		expected += "    - projection <projection> (70fdf7fa-4b24-448d-bd29-7ecc71d18c56) *github.com/dogmatiq/dogma/fixtures.ProjectionMessageHandler\n"
		expected += "        handles github.com/dogmatiq/dogma/fixtures.MessageE!\n"
		expected += "        handles github.com/dogmatiq/dogma/fixtures.MessageJ!\n"

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
