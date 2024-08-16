package static_test

import (
	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/configkit/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/packages"
)

var _ = Describe("func FromPackages() (application identity)", func() {
	When("the identity is specified with non-literal constants", func() {
		It("uses the values from the constants", func() {
			cfg := packages.Config{
				Mode: LoadPackagesConfigMode,
				Dir:  "testdata/ident/const-value-ident",
			}

			pkgs := loadPackages(cfg)

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))

			Expect(apps[0].Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<app>",
						Key:  "04e12cf2-3c66-4414-9203-e045ddbe02c7",
					},
				),
			)
			Expect(apps[0].TypeName()).To(Equal("github.com/dogmatiq/configkit/static/testdata/ident/const-value-ident.App"))
			Expect(apps[0].MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Produced: message.NameRoles{},
					Consumed: message.NameRoles{},
				},
			))
			Expect(apps[0].Handlers()).To(Equal(configkit.HandlerSet{}))
		})
	})

	When("the identity is specified with string literals", func() {
		It("uses the literal values", func() {
			cfg := packages.Config{
				Mode: LoadPackagesConfigMode,
				Dir:  "testdata/ident/literal-value-ident",
			}

			pkgs := loadPackages(cfg)

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))

			Expect(apps[0].Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<app>",
						Key:  "9d0af85d-f506-4742-b676-ce87730bb1a0",
					},
				),
			)
			Expect(apps[0].TypeName()).To(Equal("github.com/dogmatiq/configkit/static/testdata/ident/literal-value-ident.App"))
			Expect(apps[0].MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Produced: message.NameRoles{},
					Consumed: message.NameRoles{},
				},
			))
			Expect(apps[0].Handlers()).To(Equal(configkit.HandlerSet{}))
		})
	})

	When("the identity is specified with non-constant expressions", func() {
		It("uses a zero-value identity", func() {
			cfg := packages.Config{
				Mode: LoadPackagesConfigMode,
				Dir:  "testdata/ident/variable-value-ident",
			}

			pkgs := loadPackages(cfg)

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))

			Expect(apps[0].Identity()).To(
				Equal(
					configkit.Identity{
						Name: "",
						Key:  "",
					},
				),
			)
			Expect(apps[0].TypeName()).To(Equal("github.com/dogmatiq/configkit/static/testdata/ident/variable-value-ident.App"))
			Expect(apps[0].MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Produced: message.NameRoles{},
					Consumed: message.NameRoles{},
				},
			))
			Expect(apps[0].Handlers()).To(Equal(configkit.HandlerSet{}))
		})
	})
})
