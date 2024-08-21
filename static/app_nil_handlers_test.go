package static_test

import (
	. "github.com/dogmatiq/configkit/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("registering nil-valued handlers", func() {
	It("does not include the handlers", func() {
		apps := FromDir("testdata/apps/nil-handlers")
		Expect(apps).To(HaveLen(1))
		Expect(apps[0].Handlers()).To(BeEmpty())
	})
})
