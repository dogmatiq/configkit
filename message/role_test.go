package message_test

import (
	. "github.com/dogmatiq/configkit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Role", func() {
	Describe("func MustValidate()", func() {
		It("does not panic when the role is valid", func() {
			CommandRole.MustValidate()
			EventRole.MustValidate()
			TimeoutRole.MustValidate()
		})

		It("panics when the role is not valid", func() {
			Expect(func() {
				Role("<invalid>").MustValidate()
			}).To(Panic())
		})
	})

	Describe("func Is()", func() {
		It("returns true when the role is in the given set", func() {
			Expect(CommandRole.Is(CommandRole, EventRole)).To(BeTrue())
		})

		It("returns false when the role is not in the given set", func() {
			Expect(TimeoutRole.Is(CommandRole, EventRole)).To(BeFalse())
		})
	})

	Describe("func MustBe()", func() {
		It("does not panic when the role is in the given set", func() {
			CommandRole.MustBe(CommandRole, EventRole)
		})

		It("panics when the role is not in the given set", func() {
			Expect(func() {
				TimeoutRole.MustBe(CommandRole, EventRole)
			}).To(Panic())
		})
	})

	Describe("func MustNotBe()", func() {
		It("does not panic when the role is not in the given set", func() {
			TimeoutRole.MustNotBe(CommandRole, EventRole)
		})

		It("panics when the role is in the given set", func() {
			Expect(func() {
				CommandRole.MustNotBe(CommandRole, EventRole)
			}).To(Panic())
		})
	})

	Describe("func Marker()", func() {
		It("returns the correct marker character", func() {
			Expect(CommandRole.Marker()).To(Equal("?"))
			Expect(EventRole.Marker()).To(Equal("!"))
			Expect(TimeoutRole.Marker()).To(Equal("@"))
		})
	})

	Describe("func ShortString()", func() {
		It("returns the role value as a short string", func() {
			Expect(CommandRole.ShortString()).To(Equal("CMD"))
			Expect(EventRole.ShortString()).To(Equal("EVT"))
			Expect(TimeoutRole.ShortString()).To(Equal("TMO"))
		})
	})

	Describe("func String()", func() {
		It("returns the role value as a string", func() {
			Expect(CommandRole.String()).To(Equal("command"))
			Expect(EventRole.String()).To(Equal("event"))
			Expect(TimeoutRole.String()).To(Equal("timeout"))
			Expect(Role("<invalid>").String()).To(Equal("<invalid message role: <invalid>>"))
		})
	})

	Describe("func MarshalText()", func() {
		It("marshals the role to text", func() {
			Expect(CommandRole.MarshalText()).To(Equal([]byte("command")))
			Expect(EventRole.MarshalText()).To(Equal([]byte("event")))
			Expect(TimeoutRole.MarshalText()).To(Equal([]byte("timeout")))
		})

		It("returns an error if the role is invalid", func() {
			_, err := Role("<invalid>").MarshalText()
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func UnmarshalText()", func() {
		It("unmarshals the role from text", func() {
			var r Role

			err := r.UnmarshalText([]byte("command"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).To(Equal(CommandRole))

			err = r.UnmarshalText([]byte("event"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).To(Equal(EventRole))

			err = r.UnmarshalText([]byte("timeout"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).To(Equal(TimeoutRole))
		})

		It("returns an error if the data is invalid", func() {
			var r Role

			err := r.UnmarshalText([]byte("<invalid>"))
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func MarshalBinary()", func() {
		It("marshals the role to binary", func() {
			Expect(CommandRole.MarshalBinary()).To(Equal([]byte("C")))
			Expect(EventRole.MarshalBinary()).To(Equal([]byte("E")))
			Expect(TimeoutRole.MarshalBinary()).To(Equal([]byte("T")))
		})

		It("returns an error if the role is invalid", func() {
			_, err := Role("<invalid>").MarshalBinary()
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func UnmarshalBinary()", func() {
		It("unmarshals the role from binary", func() {
			var r Role

			err := r.UnmarshalBinary([]byte("C"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).To(Equal(CommandRole))

			err = r.UnmarshalBinary([]byte("E"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).To(Equal(EventRole))

			err = r.UnmarshalBinary([]byte("T"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).To(Equal(TimeoutRole))
		})

		It("returns an error if the data is the wrong length", func() {
			var r Role

			err := r.UnmarshalBinary([]byte("<invalid>"))
			Expect(err).Should(HaveOccurred())
		})

		It("returns an error if the data does not contain a valid role", func() {
			var r Role

			err := r.UnmarshalBinary([]byte("X"))
			Expect(err).Should(HaveOccurred())
		})
	})
})
