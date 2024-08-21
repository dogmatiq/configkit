package static_test

import (
	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("registering handlers with a non-pointer method-set using a pointer", func() {
	It("includes the handlers in the application configuration", func() {
		apps := FromDir("testdata/apps/non-pointer-handlers-registered-as-pointers")
		Expect(apps).To(HaveLen(1))

		matchIdentities(
			apps[0].Handlers(),
			configkit.Identity{
				Name: "<aggregate>",
				Key:  "ee1814e3-194d-438a-916e-ee7766598646",
			},
			configkit.Identity{
				Name: "<process>",
				Key:  "39af6b34-5fa1-4f3a-b049-40a5e1d9b33b",
			},
			configkit.Identity{
				Name: "<projection>",
				Key:  "3dfcd7cd-1f63-47a1-9be7-3242bd252423",
			},
			configkit.Identity{
				Name: "<integration>",
				Key:  "1425ca64-0448-4bfd-b18d-9fe63a95995f",
			},
		)
	})
})
