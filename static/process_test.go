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

var _ = Describe("func FromPackages() (process analysis)", func() {
	When("the application contains a single process handler", func() {
		It("returns the process handler configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/processes/single-process-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers().Processes()).To(HaveLen(1))

			a := apps[0].Handlers().Processes()[0]
			Expect(a.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<process>",
						Key:  "1c964b02-623e-4ae1-82e7-1f678ef875c8",
					},
				),
			)
			Expect(a.TypeName()).To(
				Equal(
					"github.com/dogmatiq/configkit/static/testdata/processes/single-process-app.ProcessHandler",
				),
			)
			Expect(a.HandlerType()).To(Equal(configkit.ProcessHandlerType))

			Expect(a.MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.EventRole,
						cfixtures.MessageBTypeName: message.EventRole,
						cfixtures.MessageETypeName: message.TimeoutRole,
						cfixtures.MessageFTypeName: message.TimeoutRole,
					},
					Produced: message.NameRoles{
						cfixtures.MessageCTypeName: message.CommandRole,
						cfixtures.MessageDTypeName: message.CommandRole,
						cfixtures.MessageETypeName: message.TimeoutRole,
						cfixtures.MessageFTypeName: message.TimeoutRole,
					},
				},
			))
		})
	})

	When("the application contains multiple process handlers", func() {
		It("returns all of the process handler configurations", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/processes/multiple-process-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))

			var identities []configkit.Identity
			for _, a := range apps[0].Handlers().Processes() {
				identities = append(identities, a.Identity())
			}

			Expect(identities).To(
				ConsistOf(
					configkit.Identity{
						Name: "<first-process>",
						Key:  "d33198e0-f1f7-4c2d-8ac2-98f68a44414e",
					},
					configkit.Identity{
						Name: "<second-process>",
						Key:  "0311717c-cf51-4292-8ed9-95125302a18e",
					},
				),
			)
		})
	})

	When("a nil value is passed as an process handler", func() {
		It("does not add an process handler to the application configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/processes/nil-process-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers()).To(Equal(configkit.HandlerSet{}))
		})
	})

	When("a nil value is passed as a message", func() {
		It("does not add the message to the process configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/processes/nil-message-process-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers().Processes()).To(HaveLen(1))

			a := apps[0].Handlers().Processes()[0]
			Expect(a.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<nil-message-process>",
						Key:  "d594ac29-eea2-4651-acd4-88c7c34eaab1",
					},
				),
			)
			Expect(a.TypeName()).To(
				Equal(
					"github.com/dogmatiq/configkit/static/testdata/processes/nil-message-process-app.ProcessHandler",
				),
			)
			Expect(a.HandlerType()).To(Equal(configkit.ProcessHandlerType))
			Expect(a.MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.EventRole,
						cfixtures.MessageCTypeName: message.TimeoutRole,
					},
					Produced: message.NameRoles{
						cfixtures.MessageBTypeName: message.CommandRole,
						cfixtures.MessageCTypeName: message.TimeoutRole,
					},
				},
			))
		})
	})
})
