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

var _ = Describe("func FromPackages() (projection analysis)", func() {
	When("the application contains a single projection handler", func() {
		It("returns the projection handler configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/projections/single-projection-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers().Projections()).To(HaveLen(1))

			a := apps[0].Handlers().Projections()[0]
			Expect(a.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<projection>",
						Key:  "1c964b02-623e-4ae1-82e7-1f678ef875c8",
					},
				),
			)
			Expect(a.TypeName()).To(
				Equal(
					"github.com/dogmatiq/configkit/static/testdata/projections/single-projection-app.ProjectionHandler",
				),
			)
			Expect(a.HandlerType()).To(Equal(configkit.ProjectionHandlerType))

			Expect(a.MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.EventRole,
						cfixtures.MessageBTypeName: message.EventRole,
					},
					Produced: message.NameRoles{},
				},
			))
		})
	})

	When("the application contains multiple projection handlers", func() {
		It("returns all of the projection handler configurations", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/projections/multiple-projection-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))

			var identities []configkit.Identity
			for _, a := range apps[0].Handlers().Projections() {
				identities = append(identities, a.Identity())
			}

			Expect(identities).To(
				ConsistOf(
					configkit.Identity{
						Name: "<first-projection>",
						Key:  "9174783f-4f12-4619-b5c6-c4ab70bd0937",
					},
					configkit.Identity{
						Name: "<second-projection>",
						Key:  "2e22850e-7c84-4b3f-b8b3-25ac743d90f2",
					},
				),
			)
		})
	})

	When("a nil value is passed as a projection handler", func() {
		It("does not add an process handler to the application configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/projections/nil-projection-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers()).To(Equal(configkit.HandlerSet{}))
		})
	})

	When("a nil value is passed as a message", func() {
		It("does not add the message to the projection configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/projections/nil-message-projection-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers().Projections()).To(HaveLen(1))

			a := apps[0].Handlers().Projections()[0]
			Expect(a.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<nil-message-projection>",
						Key:  "d594ac29-eea2-4651-acd4-88c7c34eaab1",
					},
				),
			)
			Expect(a.TypeName()).To(
				Equal(
					"github.com/dogmatiq/configkit/static/testdata/projections/nil-message-projection-app.ProjectionHandler",
				),
			)
			Expect(a.HandlerType()).To(Equal(configkit.ProjectionHandlerType))
			Expect(a.MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.EventRole,
					},
					Produced: message.NameRoles{},
				},
			))
		})
	})
})
