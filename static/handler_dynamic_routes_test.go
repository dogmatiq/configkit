package static_test

import (
	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/static"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	"github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("dynamic list of routes passed to configurer in a slice", func() {
	It("includes the routes from the slice", func() {
		apps := FromDir("testdata/handlers/dynamic-routes")
		Expect(apps).To(HaveLen(1))
		Expect(apps[0].Handlers().Integrations()).To(HaveLen(1))

		integration := apps[0].Handlers().Integrations()[0]
		Expect(integration.MessageNames()).To(Equal(
			configkit.EntityMessages[message.Name]{
				message.NameOf(CommandA1): {
					Kind:       message.CommandKind,
					IsConsumed: true,
				},
				message.NameOf(EventA1): {
					Kind:       message.EventKind,
					IsProduced: true,
				},
			},
		))
	})
})
