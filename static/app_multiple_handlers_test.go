package static_test

import (
	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("registering multiple handlers of each type", func() {
	It("includes all of the handlers", func() {
		apps := FromDir("testdata/apps/multiple-handlers")
		Expect(apps).To(HaveLen(1))

		matchIdentities(
			apps[0].Handlers(),
			configkit.Identity{
				Name: "<first-aggregate>",
				Key:  "e6300d8d-6530-405e-9729-e9ca21df23d3",
			},
			configkit.Identity{
				Name: "<second-aggregate>",
				Key:  "feeb96d0-c56b-4e58-9cd0-d393683c2ec7",
			},
			configkit.Identity{
				Name: "<first-process>",
				Key:  "d33198e0-f1f7-4c2d-8ac2-98f68a44414e",
			},
			configkit.Identity{
				Name: "<second-process>",
				Key:  "0311717c-cf51-4292-8ed9-95125302a18e",
			},
			configkit.Identity{
				Name: "<first-projection>",
				Key:  "9174783f-4f12-4619-b5c6-c4ab70bd0937",
			},
			configkit.Identity{
				Name: "<second-projection>",
				Key:  "2e22850e-7c84-4b3f-b8b3-25ac743d90f2",
			},
			configkit.Identity{
				Name: "<first-integration>",
				Key:  "14cf2812-eead-43b3-9c9c-10db5b469e94",
			},
			configkit.Identity{
				Name: "<second-integration>",
				Key:  "6bed3fbc-30e2-44c7-9a5b-e440ffe370d9",
			},
		)
	})
})
