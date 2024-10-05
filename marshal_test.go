package configkit

import (
	//revive:disable:dot-imports
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/protobuf/configpb"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func ToProto()", func() {
	var app *unmarshaledApplication

	BeforeEach(func() {
		app = &unmarshaledApplication{
			ident:    MustNewIdentity("<app>", "28c19ec0-a32f-4479-bb1d-02887e90077c"),
			typeName: "<app type>",
			handlers: HandlerSet{
				MustNewIdentity("<handler>", "3c73fa07-1073-4cf3-a208-644e26b747d7"): &unmarshaledHandler{
					ident: MustNewIdentity("<handler>", "3c73fa07-1073-4cf3-a208-644e26b747d7"),
					names: EntityMessages[message.Name]{
						message.NameOf(CommandA1): {
							Kind:       message.CommandKind,
							IsProduced: true,
						},
						message.NameOf(EventA1): {
							Kind:       message.EventKind,
							IsConsumed: true,
						},
					},
					typeName:    "<handler type>",
					handlerType: IntegrationHandlerType,
				},
			},
		}
	})

	It("produces a value that can be unmarshaled to an equivalent application", func() {
		marshaled, err := ToProto(app)
		Expect(err).ShouldNot(HaveOccurred())

		unmarshaled, err := FromProto(marshaled)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(ToString(unmarshaled)).To(Equal(ToString(app)))
		Expect(unmarshaled).To(Equal(app))
		Expect(IsApplicationEqual(unmarshaled, app)).To(BeTrue())
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
	var app *configpb.Application

	BeforeEach(func() {
		app = &configpb.Application{
			Identity: &identitypb.Identity{
				Name: "<name>",
				Key:  uuidpb.Generate(),
			},
			GoType: "<app type>",
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
		app.Handlers = append(app.Handlers, &configpb.Handler{})
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
			names:       EntityMessages[message.Name]{},
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

	It("returns an error if there is an invalid message name", func() {
		handler.names[""] = EntityMessage{}

		_, err := marshalHandler(handler)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalHandler()", func() {
	var handler *configpb.Handler

	BeforeEach(func() {
		handler = &configpb.Handler{
			Identity: &identitypb.Identity{
				Name: "<name>",
				Key:  uuidpb.Generate(),
			},
			GoType: "<handler type>",
			Type:   configpb.HandlerType_AGGREGATE,
		}
	})

	It("returns an error if the identity is invalid", func() {
		handler.Identity.Name = ""
		_, err := unmarshalHandler(handler, nil)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the type name is empty", func() {
		handler.GoType = ""
		_, err := unmarshalHandler(handler, nil)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the handler type is invalid", func() {
		handler.Type = configpb.HandlerType_UNKNOWN_HANDLER_TYPE
		_, err := unmarshalHandler(handler, nil)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if there is a message with an invalid name", func() {
		handler.Messages = map[string]*configpb.MessageUsage{
			"": {},
		}
		_, err := unmarshalHandler(handler, nil)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if there is a message with no associated kind", func() {
		handler.Messages = map[string]*configpb.MessageUsage{
			"pkg.Command": {},
		}
		_, err := unmarshalHandler(handler, nil)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func marshalIdentity()", func() {
	It("returns the protobuf representation", func() {
		key := uuidpb.Generate()
		in := MustNewIdentity("<name>", key.AsString())

		out, err := marshalIdentity(in)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).To(Equal(&identitypb.Identity{
			Name: "<name>",
			Key:  key,
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
		in := &identitypb.Identity{
			Name: "<name>",
			Key:  uuidpb.Generate(),
		}
		out, err := unmarshalIdentity(in)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).To(Equal(
			MustNewIdentity("<name>", in.Key.AsString()),
		))
	})

	It("returns an error if the identity is invalid", func() {
		in := &identitypb.Identity{}
		_, err := unmarshalIdentity(in)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func marshalHandlerType()", func() {
	DescribeTable(
		"returns the expected enumeration value",
		func(in HandlerType, expect configpb.HandlerType) {
			out, err := marshalHandlerType(in)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).To(Equal(expect))
		},
		Entry("aggregate", AggregateHandlerType, configpb.HandlerType_AGGREGATE),
		Entry("process", ProcessHandlerType, configpb.HandlerType_PROCESS),
		Entry("integration", IntegrationHandlerType, configpb.HandlerType_INTEGRATION),
		Entry("projection", ProjectionHandlerType, configpb.HandlerType_PROJECTION),
	)

	It("returns an error if the handler type is invalid", func() {
		_, err := marshalHandlerType("<invalid>")
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalHandlerType()", func() {
	DescribeTable(
		"returns the expected enumeration value",
		func(in configpb.HandlerType, expect HandlerType) {
			out, err := unmarshalHandlerType(in)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).To(Equal(expect))
		},
		Entry("aggregate", configpb.HandlerType_AGGREGATE, AggregateHandlerType),
		Entry("process", configpb.HandlerType_PROCESS, ProcessHandlerType),
		Entry("integration", configpb.HandlerType_INTEGRATION, IntegrationHandlerType),
		Entry("projection", configpb.HandlerType_PROJECTION, ProjectionHandlerType),
	)

	It("returns an error if the handler type is invalid", func() {
		_, err := unmarshalHandlerType(configpb.HandlerType_UNKNOWN_HANDLER_TYPE)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func marshalMessageKind()", func() {
	DescribeTable(
		"returns the expected enumeration value",
		func(in message.Kind, expect configpb.MessageKind) {
			out, err := marshalMessageKind(in)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).To(Equal(expect))
		},
		Entry("command", message.CommandKind, configpb.MessageKind_COMMAND),
		Entry("event", message.EventKind, configpb.MessageKind_EVENT),
		Entry("timeout", message.TimeoutKind, configpb.MessageKind_TIMEOUT),
	)

	It("returns an error if the message kind is invalid", func() {
		_, err := marshalMessageKind(-1)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalMessageKind()", func() {
	DescribeTable(
		"returns the expected enumeration value",
		func(in configpb.MessageKind, expect message.Kind) {
			out, err := unmarshalMessageKind(in)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).To(Equal(expect))
		},
		Entry("command", configpb.MessageKind_COMMAND, message.CommandKind),
		Entry("event", configpb.MessageKind_EVENT, message.EventKind),
		Entry("timeout", configpb.MessageKind_TIMEOUT, message.TimeoutKind),
	)

	It("returns an error if the message kind is invalid", func() {
		_, err := unmarshalMessageKind(configpb.MessageKind_UNKNOWN_MESSAGE_KIND)
		Expect(err).Should(HaveOccurred())
	})
})
