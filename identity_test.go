package configkit_test

import (
	"fmt"

	. "github.com/dogmatiq/configkit"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

const (
	appKey         = "59a82a24-a181-41e8-9b93-17a6ce86956e"
	aggregateKey   = "14769f7f-87fe-48dd-916e-5bcab6ba6aca"
	processKey     = "bea52cf4-e403-4b18-819d-88ade7836308"
	integrationKey = "e28f056e-e5a0-4ee7-aaf1-1d1fe02fb6e3"
	projectionKey  = "70fdf7fa-4b24-448d-bd29-7ecc71d18c56"
)

var _ = Describe("type Identity", func() {
	Describe("func NewIdentity()", func() {
		It("returns the identity", func() {
			i, err := NewIdentity("<name>", appKey)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(i).To(Equal(Identity{"<name>", appKey}))
		})

		It("returns an error if the name is invalid", func() {
			_, err := NewIdentity("", appKey)
			Expect(err).Should(HaveOccurred())
		})

		It("returns an error if the key is invalid", func() {
			_, err := NewIdentity("<name>", "")
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func MustNewIdentity()", func() {
		It("returns the identity", func() {
			i := MustNewIdentity("<name>", appKey)
			Expect(i).To(Equal(Identity{"<name>", appKey}))
		})

		It("panics if the name is invalid", func() {
			Expect(func() {
				MustNewIdentity("", appKey)
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
			Expect(Identity{"<name>", appKey}.IsZero()).To(BeFalse())
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
				Identity{"<name>", appKey},
				Identity{"<name>", appKey},
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
				Identity{"<name-1>", appKey},
				Identity{"<name-2>", appKey},
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
				i := Identity{
					v,
					uuid.NewString(),
				}
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
				i := Identity{v, appKey}
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
						"invalid key %#v, keys must be RFC 4122 UUIDs",
						v,
					),
				))
			},
			invalidEntries...,
		)
	})

	Describe("func String()", func() {
		It("returns a string representation of the identity", func() {
			i := Identity{"<name>", appKey}
			Expect(i.String()).To(Equal("<name>/59a82a24-a181-41e8-9b93-17a6ce86956e"))
		})
	})

	Describe("func MarshalText()", func() {
		It("marshals the identity to text", func() {
			Expect(
				Identity{"<name>", appKey}.MarshalText(),
			).To(
				Equal([]byte("<name> 59a82a24-a181-41e8-9b93-17a6ce86956e")),
			)
		})

		It("returns an error if the identity is invalid", func() {
			_, err := Identity{}.MarshalText()
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func UnmarshalText()", func() {
		It("unmarshals the type from text", func() {
			var i Identity

			err := i.UnmarshalText([]byte("<name> 59a82a24-a181-41e8-9b93-17a6ce86956e"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(i).To(Equal(Identity{"<name>", appKey}))
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
			in := Identity{"<name>", appKey}

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
