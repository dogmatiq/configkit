package static_test

import (
	. "github.com/dogmatiq/configkit/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("added nil-valued routes", func() {
	It("does not include the routes", func() {
		apps := FromDir("testdata/handlers/nil-routes")
		Expect(apps).To(HaveLen(1))
		Expect(apps[0].Handlers().Integrations()).To(HaveLen(1))

		integration := apps[0].Handlers().Integrations()[0]
		Expect(integration.MessageNames().All()).To(BeEmpty())
	})
})
