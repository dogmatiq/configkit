package static_test

import (
	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/internal/entity"
	. "github.com/dogmatiq/configkit/internal/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/packages"
)

var _ = Describe("func FromPackages()", func() {
	When("a package contains a single Dogma application", func() {
		It("returns the configuration for the application detected", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/apps/single-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)

			Expect(apps).To(Equal(
				[]configkit.Application{
					&entity.Application{
						IdentityValue: configkit.Identity{
							Name: "<app>",
							Key:  "8a6baab1-ee64-402e-a081-e43f4bebc243",
						},
						TypeNameValue: "github.com/dogmatiq/configkit/internal/static/testdata/apps/single-app.App",
						MessageNamesValue: configkit.EntityMessageNames{
							Produced: nil,
							Consumed: nil,
						},
						HandlersValue: nil,
					},
				},
			))
		})
	})

	When("a package contains multiple Dogma applications", func() {
		It("returns the configuration for all applications detected", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/apps/multiple-apps",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)

			Expect(apps).To(Equal(
				[]configkit.Application{
					&entity.Application{
						IdentityValue: configkit.Identity{
							Name: "<app-first>",
							Key:  "b754902b-47c8-48fc-84d2-d920c9cbdaec",
						},
						TypeNameValue: "github.com/dogmatiq/configkit/internal/static/testdata/apps/multiple-apps/first.App",
						MessageNamesValue: configkit.EntityMessageNames{
							Produced: nil,
							Consumed: nil,
						},
						HandlersValue: nil,
					},
					&entity.Application{
						IdentityValue: configkit.Identity{
							Name: "<app-second>",
							Key:  "bfaf2a16-23a0-495d-8098-051d77635822",
						},
						TypeNameValue: "github.com/dogmatiq/configkit/internal/static/testdata/apps/multiple-apps/second.App",
						MessageNamesValue: configkit.EntityMessageNames{
							Produced: nil,
							Consumed: nil,
						},
						HandlersValue: nil,
					},
				},
			))
		})
	})

	When("a package contains a Dogma application implemented with pointer receiver .Configure() method", func() {
		It("returns the configuration for the application detected", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/apps/pointer-receiver-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)

			Expect(apps).To(Equal(
				[]configkit.Application{
					&entity.Application{
						IdentityValue: configkit.Identity{
							Name: "<app>",
							Key:  "b754902b-47c8-48fc-84d2-d920c9cbdaec",
						},
						TypeNameValue: "*github.com/dogmatiq/configkit/internal/static/testdata/apps/pointer-receiver-app.App",
						MessageNamesValue: configkit.EntityMessageNames{
							Produced: nil,
							Consumed: nil,
						},
						HandlersValue: nil,
					},
				},
			))
		})
	})

	When("a package does not contain any Dogma applications", func() {
		It("returns the empty slice of configkit.Application", func() {
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

	When("a package contains a Dogma application with invalid syntax", func() {
		It("does not parse the Dogma application", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/apps/invalid-syntax-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)

			Expect(apps).To(BeEmpty())
		})
	})
})
