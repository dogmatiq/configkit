package message_test

import (
	. "github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Name", func() {
	Describe("func NameOf()", func() {
		It("returns the fully-qualified name", func() {
			n := NameOf(fixtures.MessageA1)
			Expect(n.String()).To(Equal("github.com/dogmatiq/dogma/fixtures.MessageA"))
		})

		It("panics if the message is nil", func() {
			Expect(func() {
				NameOf(nil)
			}).To(Panic())
		})
	})

	Describe("func MarshalText()", func() {
		It("marshals to a textual representation", func() {
			n := NameOf(fixtures.MessageA1)

			buf, err := n.MarshalText()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf).To(Equal([]byte("github.com/dogmatiq/dogma/fixtures.MessageA")))
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

			err := n.UnmarshalText([]byte("github.com/dogmatiq/dogma/fixtures.MessageA"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(n).To(Equal(NameOf(fixtures.MessageA1)))
		})

		It("returns an error if the data is empty", func() {
			var n Name

			err := n.UnmarshalText([]byte{})
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func MarshalBinary()", func() {
		It("marshals to a textual representation", func() {
			n := NameOf(fixtures.MessageA1)

			buf, err := n.MarshalBinary()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf).To(Equal([]byte("github.com/dogmatiq/dogma/fixtures.MessageA")))
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

			err := n.UnmarshalBinary([]byte("github.com/dogmatiq/dogma/fixtures.MessageA"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(n).To(Equal(NameOf(fixtures.MessageA1)))
		})

		It("returns an error if the data is empty", func() {
			var n Name

			err := n.UnmarshalBinary([]byte{})
			Expect(err).Should(HaveOccurred())
		})
	})
})
