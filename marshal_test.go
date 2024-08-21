package configkit

import (
	"github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	"github.com/dogmatiq/interopspec/configspec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func ToProto()", func() {
	var app *unmarshaledApplication

	BeforeEach(func() {
		app = &unmarshaledApplication{
			ident:    MustNewIdentity("<name>", "28c19ec0-a32f-4479-bb1d-02887e90077c"),
			typeName: "<app type>",
			names:    EntityMessageNames{},
			handlers: HandlerSet{},
		}
	})

	It("returns an error if the identity is invalid", func() {
		app.ident.Name = ""
		_, err := ToProto(app)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the type name is empty", func() {
		app.typeName = ""
		_, err := ToProto(app)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if one of the handlers is invalid", func() {
		app.handlers.Add(&unmarshaledHandler{})
		_, err := ToProto(app)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func FromProto()", func() {
	var app *configspec.Application

	BeforeEach(func() {
		app = &configspec.Application{
			Identity: &configspec.Identity{Name: "<name>", Key: "58877f4c-7e29-4428-a38c-7eb052e32cdc"},
			GoType:   "<app type>",
		}
	})

	It("returns an error if the identity is invalid", func() {
		app.Identity.Name = ""
		_, err := FromProto(app)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the type name is empty", func() {
		app.GoType = ""
		_, err := FromProto(app)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if one of the handlers is invalid", func() {
		app.Handlers = append(app.Handlers, &configspec.Handler{})
		_, err := FromProto(app)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func marshalHandler()", func() {
	var handler *unmarshaledHandler

	BeforeEach(func() {
		handler = &unmarshaledHandler{
			ident:       MustNewIdentity("<name>", "26c19bed-f9e8-45b1-8f60-746f7ca6ef36"),
			typeName:    "example.com/somepackage.Message",
			names:       EntityMessageNames{},
			handlerType: AggregateHandlerType,
		}
	})

	It("returns an error if the identity is invalid", func() {
		handler.ident.Name = ""
		_, err := marshalHandler(handler)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the type name is empty", func() {
		handler.typeName = ""
		_, err := marshalHandler(handler)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the handler type is invalid", func() {
		handler.handlerType = "<unknown>"
		_, err := marshalHandler(handler)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the consumed name/roles are invalid", func() {
		handler.names.Consumed = message.NameRoles{
			message.NameFor[CommandStub[TypeA]](): "<unknown>",
		}

		_, err := marshalHandler(handler)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the produced name/roles are invalid", func() {
		handler.names.Produced = message.NameRoles{
			message.NameFor[CommandStub[TypeA]](): "<unknown>",
		}

		_, err := marshalHandler(handler)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalHandler()", func() {
	var handler *configspec.Handler

	BeforeEach(func() {
		handler = &configspec.Handler{
			Identity:         &configspec.Identity{Name: "<name>", Key: "71976ec1-39c6-4f7e-b16f-632ec307e35b"},
			GoType:           "<handler type>",
			Type:             configspec.HandlerType_AGGREGATE,
			ConsumedMessages: map[string]configspec.MessageRole{},
			ProducedMessages: map[string]configspec.MessageRole{},
		}
	})

	It("returns an error if the identity is invalid", func() {
		handler.Identity.Name = ""
		_, err := unmarshalHandler(handler)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the type name is empty", func() {
		handler.GoType = ""
		_, err := unmarshalHandler(handler)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the handler type is invalid", func() {
		handler.Type = configspec.HandlerType_UNKNOWN_HANDLER_TYPE
		_, err := unmarshalHandler(handler)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the consumed messages are invalid", func() {
		handler.ConsumedMessages = map[string]configspec.MessageRole{
			"<name>": configspec.MessageRole_UNKNOWN_MESSAGE_ROLE,
		}

		_, err := unmarshalHandler(handler)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the produced messages are invalid", func() {
		handler.ProducedMessages = map[string]configspec.MessageRole{
			"<name>": configspec.MessageRole_UNKNOWN_MESSAGE_ROLE,
		}

		_, err := unmarshalHandler(handler)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func marshalNameRoles()", func() {
	It("returns an error if the name can not be marshaled", func() {
		in := message.NameRoles{
			message.Name{}: message.CommandRole,
		}
		_, err := marshalNameRoles(in)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the role can not be marshaled", func() {
		in := message.NameRoles{
			message.NameFor[CommandStub[TypeA]](): "<invalid>",
		}
		_, err := marshalNameRoles(in)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalNameRoles()", func() {
	It("returns an error if the name cannot be unmarshaled", func() {
		in := map[string]configspec.MessageRole{
			"": configspec.MessageRole_COMMAND,
		}

		_, err := unmarshalNameRoles(in)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the role cannot be unmarshaled", func() {
		in := map[string]configspec.MessageRole{
			"<name>": configspec.MessageRole_UNKNOWN_MESSAGE_ROLE,
		}

		_, err := unmarshalNameRoles(in)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func marshalIdentity()", func() {
	It("returns the protobuf representation", func() {
		in := MustNewIdentity("<name>", "9c71b756-b0ab-4c97-9ac8-75fae1dc8814")
		out, err := marshalIdentity(in)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).To(Equal(&configspec.Identity{
			Name: "<name>",
			Key:  "9c71b756-b0ab-4c97-9ac8-75fae1dc8814",
		}))
	})

	It("returns an error if the identity is invalid", func() {
		in := Identity{}
		_, err := marshalIdentity(in)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalIdentity()", func() {
	It("returns the native representation", func() {
		in := &configspec.Identity{
			Name: "<name>",
			Key:  "9a63e9ce-40ce-48a7-aa26-88b20a91ec61",
		}
		out, err := unmarshalIdentity(in)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).To(Equal(
			MustNewIdentity("<name>", "9a63e9ce-40ce-48a7-aa26-88b20a91ec61"),
		))
	})

	It("returns an error if the identity is invalid", func() {
		in := &configspec.Identity{}
		_, err := unmarshalIdentity(in)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func marshalHandlerType()", func() {
	DescribeTable(
		"returns the expected enumeration value",
		func(in HandlerType, expect configspec.HandlerType) {
			out, err := marshalHandlerType(in)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).To(Equal(expect))
		},
		Entry("aggregate", AggregateHandlerType, configspec.HandlerType_AGGREGATE),
		Entry("process", ProcessHandlerType, configspec.HandlerType_PROCESS),
		Entry("integration", IntegrationHandlerType, configspec.HandlerType_INTEGRATION),
		Entry("projection", ProjectionHandlerType, configspec.HandlerType_PROJECTION),
	)

	It("returns an error if the message role is invalid", func() {
		_, err := marshalHandlerType("<invalid>")
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalHandlerType()", func() {
	DescribeTable(
		"returns the expected enumeration value",
		func(in configspec.HandlerType, expect HandlerType) {
			out, err := unmarshalHandlerType(in)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).To(Equal(expect))
		},
		Entry("aggregate", configspec.HandlerType_AGGREGATE, AggregateHandlerType),
		Entry("process", configspec.HandlerType_PROCESS, ProcessHandlerType),
		Entry("integration", configspec.HandlerType_INTEGRATION, IntegrationHandlerType),
		Entry("projection", configspec.HandlerType_PROJECTION, ProjectionHandlerType),
	)

	It("returns an error if the message role is invalid", func() {
		_, err := unmarshalHandlerType(configspec.HandlerType_UNKNOWN_HANDLER_TYPE)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func marshalMessageRole()", func() {
	DescribeTable(
		"returns the expected enumeration value",
		func(in message.Role, expect configspec.MessageRole) {
			out, err := marshalMessageRole(in)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).To(Equal(expect))
		},
		Entry("command", message.CommandRole, configspec.MessageRole_COMMAND),
		Entry("event", message.EventRole, configspec.MessageRole_EVENT),
		Entry("timeout", message.TimeoutRole, configspec.MessageRole_TIMEOUT),
	)

	It("returns an error if the message role is invalid", func() {
		_, err := marshalMessageRole("<invalid>")
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalMessageRole()", func() {
	DescribeTable(
		"returns the expected enumeration value",
		func(in configspec.MessageRole, expect message.Role) {
			out, err := unmarshalMessageRole(in)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).To(Equal(expect))
		},
		Entry("command", configspec.MessageRole_COMMAND, message.CommandRole),
		Entry("event", configspec.MessageRole_EVENT, message.EventRole),
		Entry("timeout", configspec.MessageRole_TIMEOUT, message.TimeoutRole),
	)

	It("returns an error if the message role is invalid", func() {
		_, err := unmarshalMessageRole(configspec.MessageRole_UNKNOWN_MESSAGE_ROLE)
		Expect(err).Should(HaveOccurred())
	})
})
