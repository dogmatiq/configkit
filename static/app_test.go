package static_test

import (
	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/configkit/static"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func FromPackages() (application detection)", func() {
	When("multiple packages contain applications", func() {
		It("returns all of the application configurations", func() {
			apps := FromDir("testdata/apps/multiple-apps-in-pkgs")
			Expect(apps).To(HaveLen(2))

			Expect(
				[]configkit.Identity{
					apps[0].Identity(),
					apps[1].Identity(),
				},
			).To(
				ConsistOf(
					configkit.Identity{
						Name: "<app-first>",
						Key:  "b754902b-47c8-48fc-84d2-d920c9cbdaec",
					},
					configkit.Identity{
						Name: "<app-second>",
						Key:  "bfaf2a16-23a0-495d-8098-051d77635822",
					},
				),
			)
		})
	})

	When("a single package contains multiple applications", func() {
		It("returns all of the application configurations", func() {
			apps := FromDir("testdata/apps/multiple-apps-in-single-pkg/apps")
			Expect(apps).To(HaveLen(2))

			Expect(
				[]configkit.Identity{
					apps[0].Identity(),
					apps[1].Identity(),
				},
			).To(
				ConsistOf(
					configkit.Identity{
						Name: "<app-first>",
						Key:  "4fec74a1-6ed4-46f4-8417-01e0910be8f1",
					},
					configkit.Identity{
						Name: "<app-second>",
						Key:  "6e97d403-3cb8-4a59-a7ec-74e8e219a7bc",
					},
				),
			)
		})
	})

	When("a package contains an application implemented with pointer receivers", func() {
		It("returns the application configuration", func() {
			apps := FromDir("testdata/apps/pointer-receiver-app")
			Expect(apps).To(HaveLen(1))

			Expect(apps[0].Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<app>",
						Key:  "b754902b-47c8-48fc-84d2-d920c9cbdaec",
					},
				),
			)
		})
	})

	When("an application in the package has multiple handlers", func() {
		It("returns all messages consumed or produced by all handlers", func() {
			apps := FromDir("testdata/apps/app-level-messages")
			Expect(apps).To(HaveLen(1))

			Expect(apps[0].MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Kinds: map[message.Name]message.Kind{
						message.NameOf(CommandA1): message.CommandKind,
						message.NameOf(CommandB1): message.CommandKind,
						message.NameOf(EventA1):   message.EventKind,
						message.NameOf(EventB1):   message.EventKind,
						message.NameOf(EventC1):   message.EventKind,
						message.NameOf(TimeoutA1): message.TimeoutKind,
					},
					Consumed: message.NamesOf(
						CommandA1,
						EventA1,
						EventC1,
						TimeoutA1,
					),
					Produced: message.NamesOf(
						EventA1,
						CommandB1,
						TimeoutA1,
						EventB1,
					),
				},
			))
		})
	})
})
