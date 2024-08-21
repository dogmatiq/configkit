package static_test

import (
	"github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/configkit/static"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("constructing routes that are never passed to the configurer's Routes() method", func() {
	It("does not include the unregistered routes", func() {
		apps := FromDir("testdata/handlers/unregistered-routes")
		Expect(apps).To(HaveLen(1))
		Expect(apps[0].Handlers().Integrations()).To(HaveLen(1))

		integration := apps[0].Handlers().Integrations()[0]
		Expect(integration.MessageNames().All()).ShouldNot(
			HaveKey(message.NameFor[CommandStub[TypeX]]()),
		)
	})
})
