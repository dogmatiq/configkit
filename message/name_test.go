package message_test

import (
	. "github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Name", func() {
	Describe("func NameOf()", func() {
		DescribeTable(
			"it returns the expected name",
			func(expect string, v dogma.Message) {
				n := NameOf(v)
				Expect(n.String()).To(Equal(expect))
			},
			Entry(
				"user-defined type",
				"github.com/dogmatiq/dogma/fixtures.MessageA",
				fixtures.MessageA1,
			),
			Entry(
				"basic type",
				"string",
				string(""),
			),
			Entry(
				"pointer",
				"*github.com/dogmatiq/dogma/fixtures.MessageA",
				(*fixtures.MessageA)(nil),
			),
			Entry(
				"slice",
				"[]github.com/dogmatiq/dogma/fixtures.MessageA",
				[]fixtures.MessageA{},
			),
			Entry(
				"array",
				"[5]github.com/dogmatiq/dogma/fixtures.MessageA",
				[5]fixtures.MessageA{},
			),
			Entry(
				"map",
				"map[int]github.com/dogmatiq/dogma/fixtures.MessageA",
				map[int]fixtures.MessageA{},
			),
			Entry(
				"channel",
				"chan github.com/dogmatiq/dogma/fixtures.MessageA",
				(chan fixtures.MessageA)(nil),
			),
			Entry(
				"recv-only channel",
				"<-chan github.com/dogmatiq/dogma/fixtures.MessageA",
				(<-chan fixtures.MessageA)(nil),
			),
			Entry(
				"send-only channel",
				"chan<- github.com/dogmatiq/dogma/fixtures.MessageA",
				(chan<- fixtures.MessageA)(nil),
			),
			Entry(
				"empty struct",
				"struct {}",
				struct{}{},
			),
			Entry(
				"struct w/ single field",
				"struct {F1 github.com/dogmatiq/dogma/fixtures.MessageA}",
				struct {
					F1 fixtures.MessageA
				}{},
			),
			Entry(
				"struct w/ multiple fields",
				"struct {F1 int; F2 github.com/dogmatiq/dogma/fixtures.MessageA}",
				struct {
					F1 int
					F2 fixtures.MessageA
				}{},
			),
			Entry(
				"struct w/ embedded field",
				"struct {MessageA}",
				struct {
					fixtures.MessageA
				}{},
			),
			Entry(
				"nullary func",
				"func()",
				(func())(nil),
			),
			Entry(
				"func w/ multiple intput parameters",
				"func(int, github.com/dogmatiq/dogma/fixtures.MessageA)",
				(func(int, fixtures.MessageA))(nil),
			),
			Entry(
				"func w/ single output parameter",
				"func(github.com/dogmatiq/dogma/fixtures.MessageA) error",
				(func(fixtures.MessageA) error)(nil),
			),
			Entry(
				"func w/ single output parameter",
				"func(github.com/dogmatiq/dogma/fixtures.MessageA) error",
				(func(fixtures.MessageA) error)(nil),
			),
			Entry(
				"func w/ multiple output parameters",
				"func(github.com/dogmatiq/dogma/fixtures.MessageA) (bool, error)",
				(func(fixtures.MessageA) (bool, error))(nil),
			),
		)
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

	Describe("func UnmarshalText()", func() {
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
