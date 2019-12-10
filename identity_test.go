package configkit_test

import (
	"fmt"

	. "github.com/dogmatiq/configkit"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Identity", func() {
	Describe("func NewIdentity()", func() {
		It("returns the identity", func() {
			i, err := NewIdentity("<name>", "<key>")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(i).To(Equal(Identity{"<name>", "<key>"}))
		})

		It("returns an error if the name is invalid", func() {
			_, err := NewIdentity("", "<key>")
			Expect(err).Should(HaveOccurred())
		})

		It("returns an error if the key is invalid", func() {
			_, err := NewIdentity("<name>", "")
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func MustNewIdentity()", func() {
		It("returns the identity", func() {
			i := MustNewIdentity("<name>", "<key>")
			Expect(i).To(Equal(Identity{"<name>", "<key>"}))
		})

		It("panics if the name is invalid", func() {
			Expect(func() {
				MustNewIdentity("", "<key>")
			}).To(Panic())
		})

		It("panics if the key is invalid", func() {
			Expect(func() {
				MustNewIdentity("<name>", "")
			}).To(Panic())
		})
	})

	Describe("func IsZero()", func() {
		It("returns true if the identity is empty", func() {
			Expect(Identity{}.IsZero()).To(BeTrue())
		})

		It("returns false if the identity is not empty", func() {
			Expect(Identity{"<name>", "<key>"}.IsZero()).To(BeFalse())
		})
	})

	Describe("func ConflictsWith()", func() {
		DescribeTable(
			"it indicates whether the identities are in conflict",
			func(expect bool, a, b Identity) {
				Expect(a.ConflictsWith(b)).To(Equal(expect))
				Expect(b.ConflictsWith(a)).To(Equal(expect))
			},
			Entry(
				"same identity",
				true,
				Identity{"<name>", "<key>"},
				Identity{"<name>", "<key>"},
			),
			Entry(
				"conflicting name",
				true,
				Identity{"<name>", "<key-1>"},
				Identity{"<name>", "<key-2>"},
			),
			Entry(
				"conflicting key",
				true,
				Identity{"<name-1>", "<key>"},
				Identity{"<name-2>", "<key>"},
			),
			Entry(
				"no conflict",
				false,
				Identity{"<name-1>", "<key-1>"},
				Identity{"<name-2>", "<key-2>"},
			),
		)
	})

	Describe("func Validate()", func() {
		DescribeTable(
			"it returns nil if the name and key are valid",
			func(v string) {
				i := Identity{v, v}
				Expect(i.Validate()).ShouldNot(HaveOccurred())
			},
			Entry("ascii", "foo-bar"),
			Entry("unicode", "😀"),
		)

		invalidEntries := []TableEntry{
			Entry("empty", ""),
			Entry("non-printable ascii character (newline)", "\n"),
			Entry("non-printable ascii character (space)", " "),
			Entry("non-printable unicode character", "\u200B"),
		}

		DescribeTable(
			"it returns an error if the name is invalid",
			func(v string) {
				i := Identity{v, "<key>"}
				Expect(i.Validate()).Should(MatchError(
					fmt.Sprintf(
						"invalid name %#v, names must be non-empty, printable UTF-8 strings with no whitespace",
						v,
					),
				))
			},
			invalidEntries...,
		)

		DescribeTable(
			"it returns an error if the key is invalid",
			func(v string) {
				i := Identity{"<name>", v}
				Expect(i.Validate()).To(MatchError(
					fmt.Sprintf(
						"invalid key %#v, keys must be non-empty, printable UTF-8 strings with no whitespace",
						v,
					),
				))
			},
			invalidEntries...,
		)
	})

	Describe("func String()", func() {
		It("returns a string representation of the identity", func() {
			i := Identity{"<name>", "<key>"}
			Expect(i.String()).To(Equal("<name>/<key>"))
		})
	})

	Describe("func MarshalText()", func() {
		It("marshals the identity to text", func() {
			Expect(Identity{"<name>", "<key>"}.MarshalText()).To(Equal([]byte("<name> <key>")))
		})

		It("returns an error if the identity is invalid", func() {
			_, err := Identity{}.MarshalText()
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func UnmarshalText()", func() {
		It("unmarshals the type from text", func() {
			var i Identity

			err := i.UnmarshalText([]byte("<name> <key>"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(i).To(Equal(Identity{"<name>", "<key>"}))
		})

		It("returns an error if there is no space separator", func() {
			var i Identity

			err := i.UnmarshalText([]byte("<invalid>"))
			Expect(err).Should(HaveOccurred())
		})

		It("returns an error if the data is invalid", func() {
			var i Identity

			err := i.UnmarshalText([]byte("<name> \u200B"))
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func MarshalBinary() and UnmarshalBinary()", func() {
		It("marshals and unmarshals the identity", func() {
			in := Identity{"<name>", "<key>"}

			data, err := in.MarshalBinary()
			Expect(err).ShouldNot(HaveOccurred())

			var out Identity
			err = out.UnmarshalBinary(data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).To(Equal(in))
		})

		It("returns an error if the identity is invalid", func() {
			_, err := Identity{}.MarshalBinary()
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func UnmarshalBinary()", func() {
		It("returns an error if the data is invalid", func() {
			var i Identity

			err := i.UnmarshalBinary([]byte("\u200B"))
			Expect(err).Should(HaveOccurred())
		})
	})
})
