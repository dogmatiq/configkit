package static_test

import (
	"github.com/dogmatiq/configkit"
	cfixtures "github.com/dogmatiq/configkit/fixtures"
	"github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/configkit/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/packages"
)

var _ = Describe("func FromPackages() (integration analysis)", func() {
	When("the application contains a single integration handler", func() {
		It("returns the integration handler configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/integrations/single-integration-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers().Integrations()).To(HaveLen(1))

			a := apps[0].Handlers().Integrations()[0]
			Expect(a.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<integration>",
						Key:  "ef16c9d1-d7b6-4c99-a0e7-a59218e544fc",
					},
				),
			)
			Expect(a.TypeName()).To(
				Equal(
					"github.com/dogmatiq/configkit/static/testdata/integrations/single-integration-app.IntegrationHandler",
				),
			)
			Expect(a.HandlerType()).To(Equal(configkit.IntegrationHandlerType))

			Expect(a.MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.CommandRole,
						cfixtures.MessageBTypeName: message.CommandRole,
					},
					Produced: message.NameRoles{
						cfixtures.MessageCTypeName: message.EventRole,
						cfixtures.MessageDTypeName: message.EventRole,
					},
				},
			))
		})
	})

	When("the application contains multiple integrations handlers", func() {
		It("returns all of the integration handler configurations", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/integrations/multiple-integration-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))

			var identities []configkit.Identity
			for _, a := range apps[0].Handlers().Integrations() {
				identities = append(identities, a.Identity())
			}

			Expect(identities).To(
				ConsistOf(
					configkit.Identity{
						Name: "<first-integration>",
						Key:  "14cf2812-eead-43b3-9c9c-10db5b469e94",
					},
					configkit.Identity{
						Name: "<second-integration>",
						Key:  "6bed3fbc-30e2-44c7-9a5b-e440ffe370d9",
					},
				),
			)
		})
	})

	When("a nil value is passed as an integration handler", func() {
		It("does not add an integration handler to the application configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/integrations/nil-integration-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers()).To(Equal(configkit.HandlerSet{}))
		})
	})

	When("a nil value is passed as a message", func() {
		It("does not add the message to the integration configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/integrations/nil-message-integration-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers().Integrations()).To(HaveLen(1))

			a := apps[0].Handlers().Integrations()[0]
			Expect(a.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<nil-message-integration>",
						Key:  "e389bf6d-66fc-4355-b6f4-1e00fca724c5",
					},
				),
			)
			Expect(a.TypeName()).To(
				Equal(
					"github.com/dogmatiq/configkit/static/testdata/integrations/nil-message-integration-app.IntegrationHandler",
				),
			)
			Expect(a.HandlerType()).To(Equal(configkit.IntegrationHandlerType))
			Expect(a.MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.CommandRole,
					},
					Produced: message.NameRoles{
						cfixtures.MessageBTypeName: message.EventRole,
					},
				},
			))
		})
	})
})
