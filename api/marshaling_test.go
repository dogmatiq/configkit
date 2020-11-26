package api

import (
	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/api/internal/pb"
	"github.com/dogmatiq/configkit/internal/entity"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma/fixtures"
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func marshalApplication()", func() {
	var app *entity.Application

	BeforeEach(func() {
		app = &entity.Application{
			IdentityValue: configkit.MustNewIdentity("<name>", "<key>"),
			TypeNameValue: "<app type>",
			Messages:      configkit.EntityMessageNames{},
			HandlerSet:    configkit.HandlerSet{},
		}
	})

	It("returns an error if the identity is invalid", func() {
		app.IdentityValue.Name = ""
		_, err := marshalApplication(app)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the type name is empty", func() {
		app.TypeNameValue = ""
		_, err := marshalApplication(app)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if one of the handlers is invalid", func() {
		app.HandlerSet.Add(&entity.Handler{})
		_, err := marshalApplication(app)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalApplication()", func() {
	var app *pb.Application

	BeforeEach(func() {
		app = &pb.Application{
			Identity: &pb.Identity{Name: "<name>", Key: "<key>"},
			TypeName: "<app type>",
		}
	})

	It("returns an error if the identity is invalid", func() {
		app.Identity.Name = ""
		_, err := unmarshalApplication(app)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the type name is empty", func() {
		app.TypeName = ""
		_, err := unmarshalApplication(app)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the name roles are invalid", func() {
		app.Messages = append(app.Messages, &pb.NameRole{})
		_, err := unmarshalApplication(app)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if one of the handlers is invalid", func() {
		app.Handlers = append(app.Handlers, &pb.Handler{})
		_, err := unmarshalApplication(app)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func marshalHandler()", func() {
	var (
		indices map[message.Name]uint32
		hnd     *entity.Handler
	)

	BeforeEach(func() {
		indices = nil
		hnd = &entity.Handler{
			IdentityValue:    configkit.MustNewIdentity("<name>", "<key>"),
			TypeNameValue:    "github.com/dogmatiq/dogma/fixtures.MessageA",
			Messages:         configkit.EntityMessageNames{},
			HandlerTypeValue: configkit.AggregateHandlerType,
		}
	})

	It("returns an error if the identity is invalid", func() {
		hnd.IdentityValue.Name = ""
		_, err := marshalHandler(nil, indices, hnd)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the type name is empty", func() {
		hnd.TypeNameValue = ""
		_, err := marshalHandler(nil, indices, hnd)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the handler type is invalid", func() {
		hnd.HandlerTypeValue = "<unknown>"
		_, err := marshalHandler(nil, indices, hnd)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the consumed name/roles are invalid", func() {
		hnd.Messages.Consumed = message.NameRoles{
			message.NameOf(MessageA{}): "<unknown>",
		}

		_, err := marshalHandler(nil, indices, hnd)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the produced name/roles are invalid", func() {
		hnd.Messages.Produced = message.NameRoles{
			message.NameOf(MessageA{}): "<unknown>",
		}

		_, err := marshalHandler(nil, indices, hnd)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalHandler()", func() {
	var (
		indices []nameRole
		hnd     *pb.Handler
	)

	BeforeEach(func() {
		indices = nil
		hnd = &pb.Handler{
			Identity: &pb.Identity{Name: "<name>", Key: "<key>"},
			TypeName: "<handler type>",
			Type:     pb.HandlerType_AGGREGATE,
		}
	})

	It("returns an error if the identity is invalid", func() {
		hnd.Identity.Name = ""
		_, err := unmarshalHandler(indices, hnd)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the type name is empty", func() {
		hnd.TypeName = ""
		_, err := unmarshalHandler(indices, hnd)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the handler type is invalid", func() {
		hnd.Type = pb.HandlerType_UNKNOWN_HANDLER_TYPE
		_, err := unmarshalHandler(indices, hnd)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the consumed name/roles are invalid", func() {
		hnd.Consumed = append(hnd.Consumed, 1)
		_, err := unmarshalHandler(indices, hnd)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the produced name/roles are invalid", func() {
		hnd.Produced = append(hnd.Produced, 1)
		_, err := unmarshalHandler(indices, hnd)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func marshalNameRoles()", func() {
	It("returns an error if the name/role can not be marshaled", func() {
		in := message.NameRoles{
			message.Name{}: message.CommandRole,
		}
		_, err := marshalNameRoles(nil, nil, in)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalNameRoles()", func() {
	It("returns an error if the index is out of range", func() {
		_, err := unmarshalNameRoles(nil, []uint32{1})
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func marshalNameRole()", func() {
	It("returns the protobuf representation", func() {
		out, err := marshalNameRole(
			message.NameOf(MessageA{}),
			message.CommandRole,
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).To(Equal(&pb.NameRole{
			Name: []byte("github.com/dogmatiq/dogma/fixtures.MessageA"),
			Role: pb.MessageRole_COMMAND,
		}))
	})

	It("returns an error if the name cannot be marshaled", func() {
		_, err := marshalNameRole(
			message.Name{},
			message.CommandRole,
		)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the role cannot be marshaled", func() {
		_, err := marshalNameRole(
			message.NameOf(fixtures.MessageA{}),
			"<invalid>",
		)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalNameRole()", func() {
	It("returns the native representation", func() {
		out, err := unmarshalNameRole(&pb.NameRole{
			Name: []byte("github.com/dogmatiq/dogma/fixtures.MessageA"),
			Role: pb.MessageRole_COMMAND,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).To(Equal(nameRole{
			Name: message.NameOf(MessageA{}),
			Role: message.CommandRole,
		}))
	})

	It("returns an error if the name cannot be unmarshaled", func() {
		_, err := unmarshalNameRole(&pb.NameRole{
			Name: []byte{},
			Role: pb.MessageRole_COMMAND,
		})
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the role cannot be unmarshaled", func() {
		_, err := unmarshalNameRole(&pb.NameRole{
			Name: []byte("github.com/dogmatiq/dogma/fixtures.MessageA"),
			Role: pb.MessageRole_UNKNOWN_MESSAGE_ROLE,
		})
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func marshalIdentity()", func() {
	It("returns the protobuf representation", func() {
		in := configkit.MustNewIdentity("<name>", "<key>")
		out, err := marshalIdentity(in)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).To(Equal(&pb.Identity{
			Name: "<name>",
			Key:  "<key>",
		}))
	})

	It("returns an error if the identity is invalid", func() {
		in := configkit.Identity{}
		_, err := marshalIdentity(in)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalIdentity()", func() {
	It("returns the native representation", func() {
		in := &pb.Identity{
			Name: "<name>",
			Key:  "<key>",
		}
		out, err := unmarshalIdentity(in)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).To(Equal(
			configkit.MustNewIdentity("<name>", "<key>"),
		))
	})

	It("returns an error if the identity is invalid", func() {
		in := &pb.Identity{}
		_, err := unmarshalIdentity(in)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func marshalHandlerType()", func() {
	DescribeTable(
		"returns the expected enumeration value",
		func(in configkit.HandlerType, expect pb.HandlerType) {
			out, err := marshalHandlerType(in)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).To(Equal(expect))
		},
		Entry("aggregate", configkit.AggregateHandlerType, pb.HandlerType_AGGREGATE),
		Entry("process", configkit.ProcessHandlerType, pb.HandlerType_PROCESS),
		Entry("integration", configkit.IntegrationHandlerType, pb.HandlerType_INTEGRATION),
		Entry("projection", configkit.ProjectionHandlerType, pb.HandlerType_PROJECTION),
	)

	It("returns an error if the message role is invalid", func() {
		_, err := marshalHandlerType("<invalid>")
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalHandlerType()", func() {
	DescribeTable(
		"returns the expected enumeration value",
		func(in pb.HandlerType, expect configkit.HandlerType) {
			out, err := unmarshalHandlerType(in)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).To(Equal(expect))
		},
		Entry("aggregate", pb.HandlerType_AGGREGATE, configkit.AggregateHandlerType),
		Entry("process", pb.HandlerType_PROCESS, configkit.ProcessHandlerType),
		Entry("integration", pb.HandlerType_INTEGRATION, configkit.IntegrationHandlerType),
		Entry("projection", pb.HandlerType_PROJECTION, configkit.ProjectionHandlerType),
	)

	It("returns an error if the message role is invalid", func() {
		_, err := unmarshalHandlerType(pb.HandlerType_UNKNOWN_HANDLER_TYPE)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func marshalMessageRole()", func() {
	DescribeTable(
		"returns the expected enumeration value",
		func(in message.Role, expect pb.MessageRole) {
			out, err := marshalMessageRole(in)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).To(Equal(expect))
		},
		Entry("command", message.CommandRole, pb.MessageRole_COMMAND),
		Entry("event", message.EventRole, pb.MessageRole_EVENT),
		Entry("timeout", message.TimeoutRole, pb.MessageRole_TIMEOUT),
	)

	It("returns an error if the message role is invalid", func() {
		_, err := marshalMessageRole("<invalid>")
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func unmarshalMessageRole()", func() {
	DescribeTable(
		"returns the expected enumeration value",
		func(in pb.MessageRole, expect message.Role) {
			out, err := unmarshalMessageRole(in)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).To(Equal(expect))
		},
		Entry("command", pb.MessageRole_COMMAND, message.CommandRole),
		Entry("event", pb.MessageRole_EVENT, message.EventRole),
		Entry("timeout", pb.MessageRole_TIMEOUT, message.TimeoutRole),
	)

	It("returns an error if the message role is invalid", func() {
		_, err := unmarshalMessageRole(pb.MessageRole_UNKNOWN_MESSAGE_ROLE)
		Expect(err).Should(HaveOccurred())
	})
})
