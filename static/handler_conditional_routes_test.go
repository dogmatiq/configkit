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
				Kinds: map[message.Name]message.Kind{
					message.NameOf(CommandA1): message.CommandKind,
					message.NameOf(CommandB1): message.CommandKind,
					message.NameOf(EventA1):   message.EventKind,
					message.NameOf(EventB1):   message.EventKind,
				},
				Consumed: message.NamesOf(
					CommandA1,
					CommandB1,
				),
				Produced: message.NamesOf(
					EventA1,
					EventB1,
				),
			},
		))
	})
})
