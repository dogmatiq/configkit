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

var _ = Describe("func FromAggregate()", func() {
	var handler *fixtures.AggregateMessageHandler

	BeforeEach(func() {
		handler = &fixtures.AggregateMessageHandler{
			ConfigureFunc: func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ConsumesCommandType(fixtures.MessageB{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		}
	})

	When("the configuration is valid", func() {
		var cfg RichAggregate

		BeforeEach(func() {
			cfg = FromAggregate(handler)
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
							cfixtures.MessageATypeName: message.CommandRole,
							cfixtures.MessageBTypeName: message.CommandRole,
							cfixtures.MessageETypeName: message.EventRole,
						},
						Produced: message.NewNameSet(
							cfixtures.MessageETypeName,
						),
						Consumed: message.NewNameSet(
							cfixtures.MessageATypeName,
							cfixtures.MessageBTypeName,
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
							cfixtures.MessageAType: message.CommandRole,
							cfixtures.MessageBType: message.CommandRole,
							cfixtures.MessageEType: message.EventRole,
						},
						Produced: message.NewTypeSet(
							cfixtures.MessageEType,
						),
						Consumed: message.NewTypeSet(
							cfixtures.MessageAType,
							cfixtures.MessageBType,
						),
					},
				))
			})
		})

		Describe("func TypeName()", func() {
			It("returns the fully-qualified type name of the handler", func() {
				Expect(cfg.TypeName()).To(Equal("*github.com/dogmatiq/dogma/fixtures.AggregateMessageHandler"))
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
					VisitAggregateFunc: func(_ context.Context, c Aggregate) error {
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
					VisitRichAggregateFunc: func(_ context.Context, c RichAggregate) error {
						Expect(c).To(BeIdenticalTo(cfg))
						return errors.New("<error>")
					},
				}

				err := cfg.AcceptRichVisitor(context.Background(), v)
				Expect(err).To(MatchError("<error>"))
			})
		})

		Describe("func HandlerType()", func() {
			It("returns AggregateHandlerType", func() {
				Expect(cfg.HandlerType()).To(Equal(AggregateHandlerType))
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
			fn func(dogma.AggregateConfigurer),
		) {
			handler.ConfigureFunc = fn

			var err error
			func() {
				defer Recover(&err)
				FromAggregate(handler)
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
			`*fixtures.AggregateMessageHandler.Configure() did not call AggregateConfigurer.Identity()`,
			func(c dogma.AggregateConfigurer) {
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		),
		Entry(
			"when the handler configures multiple identities",
			`*fixtures.AggregateMessageHandler.Configure() has already called AggregateConfigurer.Identity("<name>", "<key>")`,
			func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", "<key>")
				c.Identity("<other>", "<key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		),
		Entry(
			"when the handler configures an invalid name",
			`*fixtures.AggregateMessageHandler.Configure() called AggregateConfigurer.Identity() with an invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
			func(c dogma.AggregateConfigurer) {
				c.Identity("\t \n", "<key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		),
		Entry(
			"when the handler configures an invalid key",
			`*fixtures.AggregateMessageHandler.Configure() called AggregateConfigurer.Identity() with an invalid key "\t \n", keys must be non-empty, printable UTF-8 strings with no whitespace`,
			func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", "\t \n")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		),
		Entry(
			"when the handler does not configure any consumed command types",
			`*fixtures.AggregateMessageHandler.Configure() did not call AggregateConfigurer.ConsumesCommandType()`,
			func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", "<key>")
				c.ProducesEventType(fixtures.MessageE{})
			},
		),
		Entry(
			"when the handler configures the same consumed command type multiple times",
			`*fixtures.AggregateMessageHandler.Configure() has already called AggregateConfigurer.ConsumesCommandType(fixtures.MessageA)`,
			func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		),
		Entry(
			"when the handler does not configure any produced events",
			`*fixtures.AggregateMessageHandler.Configure() did not call AggregateConfigurer.ProducesEventType()`,
			func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesCommandType(fixtures.MessageA{})
			},
		),
		Entry(
			"when the handler configures the same produced event type multiple times",
			`*fixtures.AggregateMessageHandler.Configure() has already called AggregateConfigurer.ProducesEventType(fixtures.MessageE)`,
			func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		),
	)
})
