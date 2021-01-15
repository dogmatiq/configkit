package static_test

import (
	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/internal/entity"
	. "github.com/dogmatiq/configkit/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/packages"
)

var _ = Describe("func FromPackages() (application identity)", func() {
	When("the identity is specified with non-literal constants", func() {
		It("uses the values from the constants", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/ident/const-value-ident",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)

			Expect(apps).To(ConsistOf(
				&entity.Application{
					IdentityValue: configkit.Identity{
						Name: "<app>",
						Key:  "04e12cf2-3c66-4414-9203-e045ddbe02c7",
					},
					TypeNameValue: "github.com/dogmatiq/configkit/static/testdata/ident/const-value-ident.App",
					MessageNamesValue: configkit.EntityMessageNames{
						Produced: nil,
						Consumed: nil,
					},
					HandlersValue: nil,
				},
			))
		})
	})

	When("the identity is specified with string literals", func() {
		It("uses the literal values", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/ident/literal-value-ident",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)

			Expect(apps).To(ConsistOf(
				&entity.Application{
					IdentityValue: configkit.Identity{
						Name: "<app>",
						Key:  "9d0af85d-f506-4742-b676-ce87730bb1a0",
					},
					TypeNameValue: "github.com/dogmatiq/configkit/static/testdata/ident/literal-value-ident.App",
					MessageNamesValue: configkit.EntityMessageNames{
						Produced: nil,
						Consumed: nil,
					},
					HandlersValue: nil,
				},
			))
		})
	})

	When("the identity is specified with non-constant expressions", func() {
		It("uses a zero-value identity", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/ident/variable-value-ident",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)

			Expect(apps).To(ConsistOf(
				&entity.Application{
					IdentityValue: configkit.Identity{
						Name: "",
						Key:  "",
					},
					TypeNameValue: "github.com/dogmatiq/configkit/static/testdata/ident/variable-value-ident.App",
					MessageNamesValue: configkit.EntityMessageNames{
						Produced: nil,
						Consumed: nil,
					},
					HandlersValue: nil,
				},
			))
		})
	})
})
