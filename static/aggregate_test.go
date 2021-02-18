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

var _ = Describe("func FromPackages() (aggregate analysis)", func() {
	When("the application contains a single aggregate handler", func() {
		It("returns the aggregate handler configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/aggregates/single-aggregate-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers().Aggregates()).To(HaveLen(1))

			a := apps[0].Handlers().Aggregates()[0]
			Expect(a.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<aggregate>",
						Key:  "ef16c9d1-d7b6-4c99-a0e7-a59218e544fc",
					},
				),
			)
			Expect(a.TypeName()).To(
				Equal(
					"github.com/dogmatiq/configkit/static/testdata/aggregates/single-aggregate-app.AggregateHandler",
				),
			)
			Expect(a.HandlerType()).To(Equal(configkit.AggregateHandlerType))

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

	When("the application contains multiple aggregate handlers", func() {
		It("returns all of the aggregate handler configurations", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/aggregates/multiple-aggregate-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))

			var identities []configkit.Identity
			for _, a := range apps[0].Handlers().Aggregates() {
				identities = append(identities, a.Identity())
			}

			Expect(identities).To(
				ConsistOf(
					configkit.Identity{
						Name: "<first-aggregate>",
						Key:  "e6300d8d-6530-405e-9729-e9ca21df23d3",
					},
					configkit.Identity{
						Name: "<second-aggregate>",
						Key:  "feeb96d0-c56b-4e58-9cd0-d393683c2ec7",
					},
				),
			)
		})
	})

	When("a nil value is passed as an aggregate handler", func() {
		It("does not add an aggregate handler to the application configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/aggregates/nil-aggregate-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers()).To(Equal(configkit.HandlerSet{}))
		})
	})

	When("a nil value is passed as a message", func() {
		It("does not add the message to the aggregate configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/aggregates/nil-message-aggregate-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers().Aggregates()).To(HaveLen(1))

			a := apps[0].Handlers().Aggregates()[0]
			Expect(a.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<nil-message-aggregate>",
						Key:  "99271492-1ec3-475f-b154-3e69cda11155",
					},
				),
			)
			Expect(a.TypeName()).To(
				Equal(
					"github.com/dogmatiq/configkit/static/testdata/aggregates/nil-message-aggregate-app.AggregateHandler",
				),
			)
			Expect(a.HandlerType()).To(Equal(configkit.AggregateHandlerType))
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
