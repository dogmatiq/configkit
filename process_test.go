package configkit_test

import (
	"context"
	"errors"
	"reflect"

	. "github.com/dogmatiq/configkit"
	cfixtures "github.com/dogmatiq/configkit/fixtures" // can't dot-import due to conflicts
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures" // can't dot-import due to conflicts
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func FromProcess()", func() {
	var handler *fixtures.ProcessMessageHandler

	BeforeEach(func() {
		handler = &fixtures.ProcessMessageHandler{
			ConfigureFunc: func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ConsumesEventType(fixtures.MessageB{})
				c.ProducesCommandType(fixtures.MessageC{})
				c.SchedulesTimeoutType(fixtures.MessageT{})
			},
		}
	})

	When("the configuration is valid", func() {
		var cfg RichProcess

		BeforeEach(func() {
			cfg = FromProcess(handler)
		})

		Describe("func Identity()", func() {
			It("returns the handler identity", func() {
				Expect(cfg.Identity()).To(Equal(
					MustNewIdentity("<name>", "<key>"),
				))
			})
		})

		Describe("func MessageNames()", func() {
			It("returns the expected message names", func() {
				Expect(cfg.MessageNames()).To(Equal(
					EntityMessageNames{
						Roles: message.NameRoles{
							cfixtures.MessageATypeName: message.EventRole,
							cfixtures.MessageBTypeName: message.EventRole,
							cfixtures.MessageCTypeName: message.CommandRole,
							cfixtures.MessageTTypeName: message.TimeoutRole,
						},
						Produced: message.NewNameSet(
							cfixtures.MessageCTypeName,
							cfixtures.MessageTTypeName,
						),
						Consumed: message.NewNameSet(
							cfixtures.MessageATypeName,
							cfixtures.MessageBTypeName,
							cfixtures.MessageTTypeName,
						),
					},
				))
			})
		})

		Describe("func MessageTypes()", func() {
			It("returns the expected message types", func() {
				Expect(cfg.MessageTypes()).To(Equal(
					EntityMessageTypes{
						Roles: message.TypeRoles{
							cfixtures.MessageAType: message.EventRole,
							cfixtures.MessageBType: message.EventRole,
							cfixtures.MessageCType: message.CommandRole,
							cfixtures.MessageTType: message.TimeoutRole,
						},
						Produced: message.NewTypeSet(
							cfixtures.MessageCType,
							cfixtures.MessageTType,
						),
						Consumed: message.NewTypeSet(
							cfixtures.MessageAType,
							cfixtures.MessageBType,
							cfixtures.MessageTType,
						),
					},
				))
			})
		})

		Describe("func TypeName()", func() {
			It("returns the fully-qualified type name of the handler", func() {
				Expect(cfg.TypeName()).To(Equal("*github.com/dogmatiq/dogma/fixtures.ProcessMessageHandler"))
			})
		})

		Describe("func ReflectType()", func() {
			It("returns the type of the handler", func() {
				Expect(cfg.ReflectType()).To(Equal(reflect.TypeOf(handler)))
			})
		})

		Describe("func AcceptVisitor()", func() {
			It("calls the appropriate method on the visitor", func() {
				v := &cfixtures.Visitor{
					VisitProcessFunc: func(_ context.Context, c Process) error {
						Expect(c).To(BeIdenticalTo(cfg))
						return errors.New("<error>")
					},
				}

				err := cfg.AcceptVisitor(context.Background(), v)
				Expect(err).To(MatchError("<error>"))
			})
		})

		Describe("func AcceptRichVisitor()", func() {
			It("calls the appropriate method on the visitor", func() {
				v := &cfixtures.RichVisitor{
					VisitRichProcessFunc: func(_ context.Context, c RichProcess) error {
						Expect(c).To(BeIdenticalTo(cfg))
						return errors.New("<error>")
					},
				}

				err := cfg.AcceptRichVisitor(context.Background(), v)
				Expect(err).To(MatchError("<error>"))
			})
		})

		Describe("func HandlerType()", func() {
			It("returns ProcessHandlerType", func() {
				Expect(cfg.HandlerType()).To(Equal(ProcessHandlerType))
			})
		})

		Describe("func Handler()", func() {
			It("returns the underlying handler", func() {
				Expect(cfg.Handler()).To(BeIdenticalTo(handler))
			})
		})
	})

	DescribeTable(
		"when the configuration is invalid",
		func(
			msg string,
			fn func(dogma.ProcessConfigurer),
		) {
			handler.ConfigureFunc = fn

			var err error
			func() {
				defer Recover(&err)
				FromProcess(handler)
			}()

			Expect(err).Should(HaveOccurred())
			if msg != "" {
				Expect(err).To(MatchError(msg))
			}
		},
		Entry(
			"when the handler does not configure anything",
			"", // any error
			nil,
		),
		Entry(
			"when the handler does not configure an identity",
			`*fixtures.ProcessMessageHandler.Configure() did not call ProcessConfigurer.Identity()`,
			func(c dogma.ProcessConfigurer) {
				c.ConsumesEventType(fixtures.MessageA{})
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler configures multiple identities",
			`*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.Identity("<name>", "<key>")`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.Identity("<other>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler configures an invalid name",
			`*fixtures.ProcessMessageHandler.Configure() called ProcessConfigurer.Identity() with an invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("\t \n", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler configures an invalid key",
			`*fixtures.ProcessMessageHandler.Configure() called ProcessConfigurer.Identity() with an invalid key "\t \n", keys must be non-empty, printable UTF-8 strings with no whitespace`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "\t \n")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler does not configure any consumed event types",
			`*fixtures.ProcessMessageHandler.Configure() did not call ProcessConfigurer.ConsumesEventType()`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler configures the same consumed event type multiple times",
			`*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.ConsumesEventType(fixtures.MessageA)`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ConsumesEventType(fixtures.MessageA{})
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler configures an event that was previously configured as a timeout",
			`*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.SchedulesTimeoutType(fixtures.MessageA)`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.SchedulesTimeoutType(fixtures.MessageA{})
				c.ConsumesEventType(fixtures.MessageA{})
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler does not configure any produced commands",
			`*fixtures.ProcessMessageHandler.Configure() did not call ProcessConfigurer.ProducesCommandType()`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
			},
		),
		Entry(
			"when the handler configures the same produced command type multiple times",
			`*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.ProducesCommandType(fixtures.MessageC)`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ProducesCommandType(fixtures.MessageC{})
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler configures a command that was previously configured as a timeout",
			`*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.SchedulesTimeoutType(fixtures.MessageC)`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.SchedulesTimeoutType(fixtures.MessageC{})
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
	)
})
