package static_test

import (
	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/configkit/static"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("routes added in conditional branches", func() {
	It("includes the union of routes from all branches", func() {
		apps := FromDir("testdata/handlers/conditional-routes")
		Expect(apps).To(HaveLen(1))
		Expect(apps[0].Handlers().Integrations()).To(HaveLen(1))

		integration := apps[0].Handlers().Integrations()[0]

		Expect(integration.MessageNames()).To(Equal(
			configkit.EntityMessageNames{
				Consumed: message.NameRoles{
					message.NameFor[CommandStub[TypeA]](): message.CommandRole,
					message.NameFor[CommandStub[TypeB]](): message.CommandRole,
				},
				Produced: message.NameRoles{
					message.NameFor[EventStub[TypeA]](): message.EventRole,
					message.NameFor[EventStub[TypeB]](): message.EventRole,
				},
			},
		))
	})
})
