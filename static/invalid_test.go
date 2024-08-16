package static_test

import (
	. "github.com/dogmatiq/configkit/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/packages"
)

var _ = Describe("func FromPackages() (unbuildable packages)", func() {
	When("a package contains a file with invalid Go syntax", func() {
		It("ignores the package", func() {
			cfg := packages.Config{
				Mode: LoadPackagesConfigMode,
				Dir:  "testdata/invalid/invalid-syntax",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)

			Expect(apps).To(BeEmpty())
		})
	})
})
