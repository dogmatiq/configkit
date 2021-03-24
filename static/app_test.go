package static_test

import (
	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/packages"
)

var _ = Describe("func FromPackages() (application detection)", func() {
	When("a package contains a single application", func() {
		It("returns the application configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/apps/single-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))

			Expect(apps[0].Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<app>",
						Key:  "8a6baab1-ee64-402e-a081-e43f4bebc243",
					},
				),
			)
			Expect(apps[0].TypeName()).To(Equal("github.com/dogmatiq/configkit/static/testdata/apps/single-app.App"))
			Expect(apps[0].MessageNames()).To(Equal(
				configkit.EntityMessageNames{},
			))
			Expect(apps[0].Handlers()).To(Equal(configkit.HandlerSet{}))
		})
	})

	When("multiple packages contain applications", func() {
		It("returns all of the application configurations", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/apps/multiple-apps-in-pkgs",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
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
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/apps/multiple-apps-in-single-pkg/apps",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
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
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/apps/pointer-receiver-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)

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

	When("none of the packages contain any applications", func() {
		It("returns an empty slice", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/apps/no-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)

			Expect(apps).To(BeEmpty())
		})
	})

	When("a field within the application type is registered as a handler", func() {
		It("includes the handler in the application configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/apps/handler-from-field",
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
						Key:  "195ede4a-3f26-4d19-a8fe-41b2a5f92d06",
					},
				),
			)
		})
	})

	When("a handler with a non-pointer methodset is registered as a pointer", func() {
		It("includes the handler in the application configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/apps/pointer-handler-with-non-pointer-methodset",
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
						Key:  "dad3b670-0852-4711-9efb-af25679734ee",
					},
				),
			)
		})
	})
})
