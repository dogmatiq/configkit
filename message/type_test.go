package message_test

import (
	"reflect"

	"github.com/dogmatiq/dogma"

	. "github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Type", func() {
	Describe("func TypeOf()", func() {
		It("returns values that compare as equal for messages of the same type", func() {
			tb := TypeOf(fixtures.MessageA1)
			ta := TypeOf(fixtures.MessageA1)

			Expect(ta).To(Equal(tb))
			Expect(ta == tb).To(BeTrue()) // we're testing == here specifically, hence not using To(Equal())
		})

		It("returns values that do not compare as equal for messages of different types", func() {
			ta := TypeOf(fixtures.MessageA1)
			tb := TypeOf(fixtures.MessageB1)

			Expect(ta).NotTo(Equal(tb))
			Expect(ta != tb).To(BeTrue()) // we're testing != here specifically, hence not using NotTo(Equal())
		})
	})

	Describe("func Name()", func() {
		It("returns the fully qualified type name", func() {
			mt := TypeOf(fixtures.MessageA1)
			Expect(mt.Name()).To(Equal(NameOf(fixtures.MessageA1)))
		})
	})

	Describe("func ReflectType()", func() {
		It("returns the reflect.Type for the message", func() {
			mt := TypeOf(fixtures.MessageA1)
			rt := reflect.TypeOf(fixtures.MessageA1)

			Expect(mt.ReflectType()).To(BeIdenticalTo(rt))
		})
	})

	Describe("func String()", func() {
		DescribeTable(
			"it returns the expected short representation",
			func(expect string, v dogma.Message) {
				t := TypeOf(v)
				Expect(t.String()).To(Equal(expect))
			},
			Entry(
				"user-defined type",
				"fixtures.MessageA",
				fixtures.MessageA1,
			),
			Entry(
				"basic type",
				"string",
				"",
			),
			Entry(
				"pointer",
				"*fixtures.MessageA",
				(*fixtures.MessageA)(nil),
			),
			Entry(
				"slice",
				"[]fixtures.MessageA",
				[]fixtures.MessageA{},
			),
			Entry(
				"array",
				"[5]fixtures.MessageA",
				[5]fixtures.MessageA{},
			),
			Entry(
				"map",
				"map[int]fixtures.MessageA",
				map[int]fixtures.MessageA{},
			),
			Entry(
				"channel",
				"chan fixtures.MessageA",
				(chan fixtures.MessageA)(nil),
			),
			Entry(
				"recv-only channel",
				"<-chan fixtures.MessageA",
				(<-chan fixtures.MessageA)(nil),
			),
			Entry(
				"send-only channel",
				"chan<- fixtures.MessageA",
				(chan<- fixtures.MessageA)(nil),
			),
			Entry(
				"empty struct",
				"struct {}",
				struct{}{},
			),
			Entry(
				"struct w/ single field",
				"struct {F1 fixtures.MessageA}",
				struct {
					F1 fixtures.MessageA
				}{},
			),
			Entry(
				"struct w/ multiple fields",
				"struct {F1 int; F2 fixtures.MessageA}",
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
				"func(int, fixtures.MessageA)",
				(func(int, fixtures.MessageA))(nil),
			),
			Entry(
				"func w/ single output parameter",
				"func(fixtures.MessageA) error",
				(func(fixtures.MessageA) error)(nil),
			),
			Entry(
				"func w/ single output parameter",
				"func(fixtures.MessageA) error",
				(func(fixtures.MessageA) error)(nil),
			),
			Entry(
				"func w/ multiple output parameters",
				"func(fixtures.MessageA) (bool, error)",
				(func(fixtures.MessageA) (bool, error))(nil),
			),
		)
	})
})
