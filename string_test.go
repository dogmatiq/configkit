package configkit_test

import (
	"strings"

	. "github.com/dogmatiq/configkit" // can't dot-import due to conflicts
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
				c.Identity("<app>", "<app-key>")

				c.RegisterAggregate(&fixtures.AggregateMessageHandler{
					ConfigureFunc: func(c dogma.AggregateConfigurer) {
						c.Identity("<aggregate>", "<aggregate-key>")
						c.ConsumesCommandType(fixtures.MessageC{})
						c.ProducesEventType(fixtures.MessageE{})
					},
				})

				c.RegisterProcess(&fixtures.ProcessMessageHandler{
					ConfigureFunc: func(c dogma.ProcessConfigurer) {
						c.Identity("<process>", "<process-key>")
						c.ConsumesEventType(fixtures.MessageE{})
						c.ProducesCommandType(fixtures.MessageC{})
						c.SchedulesTimeoutType(fixtures.MessageT{})
					},
				})

				c.RegisterIntegration(&fixtures.IntegrationMessageHandler{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Identity("<integration>", "<integration-key>")
						c.ConsumesCommandType(fixtures.MessageI{})
						c.ProducesEventType(fixtures.MessageJ{})
					},
				})

				c.RegisterProjection(&fixtures.ProjectionMessageHandler{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Identity("<projection>", "<projection-key>")
						c.ConsumesEventType(fixtures.MessageE{})
						c.ConsumesEventType(fixtures.MessageJ{})
					},
				})
			},
		}

		cfg = FromApplication(app)
	})

	It("returns a human readable string representation", func() {
		expected := "application <app> (<app-key>) *github.com/dogmatiq/dogma/fixtures.Application\n"
		expected += "\n"
		expected += "    - aggregate <aggregate> (<aggregate-key>) *github.com/dogmatiq/dogma/fixtures.AggregateMessageHandler\n"
		expected += "        consumes github.com/dogmatiq/dogma/fixtures.MessageC?\n"
		expected += "        produces github.com/dogmatiq/dogma/fixtures.MessageE!\n"
		expected += "\n"
		expected += "    - integration <integration> (<integration-key>) *github.com/dogmatiq/dogma/fixtures.IntegrationMessageHandler\n"
		expected += "        consumes github.com/dogmatiq/dogma/fixtures.MessageI?\n"
		expected += "        produces github.com/dogmatiq/dogma/fixtures.MessageJ!\n"
		expected += "\n"
		expected += "    - process <process> (<process-key>) *github.com/dogmatiq/dogma/fixtures.ProcessMessageHandler\n"
		expected += "        consumes github.com/dogmatiq/dogma/fixtures.MessageE!\n"
		expected += "        produces github.com/dogmatiq/dogma/fixtures.MessageC?\n"
		expected += "        schedules github.com/dogmatiq/dogma/fixtures.MessageT@\n"
		expected += "\n"
		expected += "    - projection <projection> (<projection-key>) *github.com/dogmatiq/dogma/fixtures.ProjectionMessageHandler\n"
		expected += "        consumes github.com/dogmatiq/dogma/fixtures.MessageE!\n"
		expected += "        consumes github.com/dogmatiq/dogma/fixtures.MessageJ!\n"

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
