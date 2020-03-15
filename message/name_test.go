package message_test

import (
	. "github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func IsEqualSetN()", func() {
	It("returns true for identical sets", func() {
		a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
		Expect(IsEqualSetN(a, b)).To(BeTrue())
	})

	It("returns false for disjoint sets", func() {
		a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := NamesOf(fixtures.MessageC1, fixtures.MessageD1)
		Expect(IsEqualSetN(a, b)).To(BeFalse())
	})

	It("returns false for intersecting sets", func() {
		a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := NamesOf(fixtures.MessageB1, fixtures.MessageC1)
		Expect(IsEqualSetN(a, b)).To(BeFalse())
	})
})

var _ = Describe("func IsIntersectingN()", func() {
	It("returns true for identical sets", func() {
		a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
		Expect(IsIntersectingN(a, b)).To(BeTrue())
	})

	It("returns false for disjoint sets", func() {
		a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := NamesOf(fixtures.MessageC1, fixtures.MessageD1)
		Expect(IsIntersectingN(a, b)).To(BeFalse())
	})

	It("returns true for intersecting sets", func() {
		a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := NamesOf(fixtures.MessageB1, fixtures.MessageC1)
		Expect(IsIntersectingN(a, b)).To(BeTrue())
	})
})

var _ = Describe("func IsSubsetN()", func() {
	It("returns true for identical sets", func() {
		a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
		Expect(IsSubsetN(a, b)).To(BeTrue())
	})

	It("returns true for strict subsets", func() {
		a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := NamesOf(fixtures.MessageA1, fixtures.MessageB1, fixtures.MessageC1)
		Expect(IsSubsetN(a, b)).To(BeTrue())
	})

	It("returns false for supersets", func() {
		a := NamesOf(fixtures.MessageA1, fixtures.MessageB1, fixtures.MessageC1)
		b := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
		Expect(IsSubsetN(a, b)).To(BeFalse())
	})

	It("returns false for disjoint sets", func() {
		a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := NamesOf(fixtures.MessageC1, fixtures.MessageD1)
		Expect(IsSubsetN(a, b)).To(BeFalse())
	})
})

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
