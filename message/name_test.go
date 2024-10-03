package message_test

import (
	"go/token"
	"go/types"

	. "github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Name", func() {
	Describe("func NameOf()", func() {
		It("returns the fully-qualified name", func() {
			n := NameOf(CommandA1)
			Expect(n.String()).To(Equal("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"))
		})

		It("panics if the message is nil", func() {
			Expect(func() {
				NameOf(nil)
			}).To(Panic())
		})
	})

	Describe("func NameFor()", func() {
		It("returns the name", func() {
			na := NameFor[CommandStub[TypeA]]()
			nb := NameOf(CommandA1)
			Expect(na).To(Equal(nb))
		})
	})

	Describe("func NameFromStaticType()", func() {
		It("returns the fully-qualified name", func() {
			pkg := types.NewPackage(
				"example.org/somepackage",
				"somepackage",
			)

			named := types.NewNamed(
				types.NewTypeName(
					token.NoPos,
					pkg,
					"Message",
					&types.Struct{},
				),
				nil,
				nil,
			)

			n := NameFromStaticType(named)
			Expect(n.String()).To(Equal("example.org/somepackage.Message"))
		})

		It("panics if the type is nil", func() {
			Expect(func() {
				NameFromStaticType(nil)
			}).To(Panic())
		})
	})

	Describe("func MarshalText()", func() {
		It("marshals to a textual representation", func() {
			n := NameOf(CommandA1)

			buf, err := n.MarshalText()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf).To(Equal([]byte("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]")))
		})

		It("returns an error if the type is the zero-value", func() {
			var n Name

			_, err := n.MarshalText()
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func UnmarshalText()", func() {
		It("unmarshals from a textual representation", func() {
			var n Name

			err := n.UnmarshalText([]byte("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(n).To(Equal(NameFor[CommandStub[TypeA]]()))
		})

		It("returns an error if the data is empty", func() {
			var n Name

			err := n.UnmarshalText([]byte{})
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func MarshalBinary()", func() {
		It("marshals to a textual representation", func() {
			n := NameFor[CommandStub[TypeA]]()

			buf, err := n.MarshalBinary()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf).To(Equal([]byte("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]")))
		})

		It("returns an error if the type is the zero-value", func() {
			var n Name

			_, err := n.MarshalBinary()
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func UnmarshalBinary()", func() {
		It("unmarshals from a textual representation", func() {
			var n Name

			err := n.UnmarshalBinary([]byte("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(n).To(Equal(NameFor[CommandStub[TypeA]]()))
		})

		It("returns an error if the data is empty", func() {
			var n Name

			err := n.UnmarshalBinary([]byte{})
			Expect(err).Should(HaveOccurred())
		})
	})
})
