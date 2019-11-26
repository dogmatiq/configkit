package typename_test

import (
	"reflect"
	"strings"

	. "github.com/dogmatiq/configkit/internal/typename"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type testType struct{}

var _ = Describe("func Of()", func() {
	DescribeTable(
		"it returns the expected name",
		func(expect string, rt reflect.Type) {
			n := Of(rt)
			Expect(n).To(Equal(expect))

			// verify that the result is similar to ReflectType.String()
			s := strings.ReplaceAll(n, "github.com/dogmatiq/configkit/internal/", "")
			Expect(s).To(Equal(rt.String()))
		},
		Entry(
			"user-defined type",
			"github.com/dogmatiq/configkit/internal/typename_test.testType",
			reflect.TypeOf(testType{}),
		),
		Entry(
			"basic type",
			"string",
			reflect.TypeOf(""),
		),
		Entry(
			"pointer",
			"*github.com/dogmatiq/configkit/internal/typename_test.testType",
			reflect.TypeOf(&testType{}),
		),
		Entry(
			"slice",
			"[]github.com/dogmatiq/configkit/internal/typename_test.testType",
			reflect.TypeOf([]testType{}),
		),
		Entry(
			"array",
			"[5]github.com/dogmatiq/configkit/internal/typename_test.testType",
			reflect.TypeOf([5]testType{}),
		),
		Entry(
			"map",
			"map[int]github.com/dogmatiq/configkit/internal/typename_test.testType",
			reflect.TypeOf(map[int]testType{}),
		),
		Entry(
			"channel",
			"chan github.com/dogmatiq/configkit/internal/typename_test.testType",
			reflect.TypeOf((chan testType)(nil)),
		),
		Entry(
			"recv-only channel",
			"<-chan github.com/dogmatiq/configkit/internal/typename_test.testType",
			reflect.TypeOf((<-chan testType)(nil)),
		),
		Entry(
			"send-only channel",
			"chan<- github.com/dogmatiq/configkit/internal/typename_test.testType",
			reflect.TypeOf((chan<- testType)(nil)),
		),
		Entry(
			"empty interface",
			"interface {}",
			reflect.TypeOf((*interface{})(nil)).Elem(),
		),
		Entry(
			"interface w/ single method",
			"interface { Foo() }",
			reflect.TypeOf((*interface {
				Foo()
			})(nil)).Elem(),
		),
		Entry(
			"interface w/ multiple methods",
			"interface { Bar(); Foo() }",
			reflect.TypeOf((*interface {
				Foo()
				Bar()
			})(nil)).Elem(),
		),
		Entry(
			"empty struct",
			"struct {}",
			reflect.TypeOf(struct{}{}),
		),
		Entry(
			"struct w/ single field",
			"struct { F1 github.com/dogmatiq/configkit/internal/typename_test.testType }",
			reflect.TypeOf(struct {
				F1 testType
			}{}),
		),
		Entry(
			"struct w/ multiple fields",
			"struct { F1 int; F2 github.com/dogmatiq/configkit/internal/typename_test.testType }",
			reflect.TypeOf(struct {
				F1 int
				F2 testType
			}{}),
		),
		Entry(
			"struct w/ embedded type",
			"struct { github.com/dogmatiq/configkit/internal/typename_test.testType }",
			reflect.TypeOf(struct {
				testType
			}{}),
		),
		Entry(
			"nullary func",
			"func()",
			reflect.TypeOf(func() {}),
		),
		Entry(
			"func w/ multiple intput parameters",
			"func(int, github.com/dogmatiq/configkit/internal/typename_test.testType)",
			reflect.TypeOf(func(int, testType) {}),
		),
		Entry(
			"func w/ single output parameter",
			"func(github.com/dogmatiq/configkit/internal/typename_test.testType) error",
			reflect.TypeOf(func(testType) error { return nil }),
		),
		Entry(
			"func w/ multiple output parameters",
			"func(github.com/dogmatiq/configkit/internal/typename_test.testType) (bool, error)",
			reflect.TypeOf(func(testType) (bool, error) { return true, nil }),
		),
	)
})
