package static_test

import (
	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/internal/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/packages"
)

var _ = Describe("func FromPackages()", func() {
	When("multiple Dogma applications exist in the packages", func() {
		It("returns all Dogma application configurations", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/multiple-apps",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(2))

			Expect(apps[0].Identity()).To(Equal(
				configkit.Identity{
					Name: "app01",
					Key:  "b754902b-47c8-48fc-84d2-d920c9cbdaec",
				},
			))
			Expect(apps[0].TypeName()).To(Equal("github.com/dogmatiq/configkit/internal/static/testdata/multiple-apps/app01.App"))

			Expect(apps[1].Identity()).To(Equal(
				configkit.Identity{
					Name: "app02",
					Key:  "bfaf2a16-23a0-495d-8098-051d77635822",
				},
			))
			Expect(apps[1].TypeName()).To(Equal("*github.com/dogmatiq/configkit/internal/static/testdata/multiple-apps/app02.App"))
		})
	})

	When("no Dogma applications exist in the packages", func() {
		It("returns an empty list of Dogma application configurations", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/no-impl",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(0))
		})
	})

	When("the package with a Dogma application contains errors", func() {
		It("does not return the application configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/pkg-with-errors",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(0))
		})
	})
})
