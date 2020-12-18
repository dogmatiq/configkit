package typename

// TestType is the type used in "goreflect" and "gotypes" packages for testing.
type TestType struct{}

var (
	Basic   string
	Named   TestType
	Pointer *TestType
	Slice   []TestType
	Array   [5]TestType
	Map     map[int]TestType
)

var (
	Channel            chan TestType
	ReceiveOnlyChannel chan<- TestType
	SendOnlyChannel    <-chan TestType
)

var (
	Interface                    interface{}
	InterfaceWithSingleMethod    interface{ Foo() }
	InterfaceWithMultipleMethods interface {
		Foo()
		Bar()
	}
)

var (
	Struct                   struct{}
	StructWithSingleField    struct{ F1 TestType }
	StructWithMultipleFields struct {
		F1 int
		F2 TestType
	}
	StructWithEmbeddedType struct{ TestType }
)

var (
	Nullary                      func()
	FuncWithMultipleInputParams  func(int, TestType)
	FuncWithSingleOutputParam    func(TestType) error
	FuncWithMultipleOutputParams func(TestType) (bool, error)
)
