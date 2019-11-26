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

var _ = Describe("func FromProjection()", func() {
	var handler *fixtures.ProjectionMessageHandler

	BeforeEach(func() {
		handler = &fixtures.ProjectionMessageHandler{
			ConfigureFunc: func(c dogma.ProjectionConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ConsumesEventType(fixtures.MessageB{})
			},
		}
	})

	When("the configuration is valid", func() {
		var cfg RichProjection

		BeforeEach(func() {
			cfg = FromProjection(handler)
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
						},
						Produced: nil,
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
							cfixtures.MessageAType: message.EventRole,
							cfixtures.MessageBType: message.EventRole,
						},
						Produced: nil,
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
				Expect(cfg.TypeName()).To(Equal("*github.com/dogmatiq/dogma/fixtures.ProjectionMessageHandler"))
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
					VisitProjectionFunc: func(_ context.Context, c Projection) error {
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
					VisitRichProjectionFunc: func(_ context.Context, c RichProjection) error {
						Expect(c).To(BeIdenticalTo(cfg))
						return errors.New("<error>")
					},
				}

				err := cfg.AcceptRichVisitor(context.Background(), v)
				Expect(err).To(MatchError("<error>"))
			})
		})

		Describe("func HandlerType()", func() {
			It("returns ProjectionHandlerType", func() {
				Expect(cfg.HandlerType()).To(Equal(ProjectionHandlerType))
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
			fn func(dogma.ProjectionConfigurer),
		) {
			handler.ConfigureFunc = fn

			var err error
			func() {
				defer Recover(&err)
				FromProjection(handler)
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
			`*fixtures.ProjectionMessageHandler.Configure() did not call ProjectionConfigurer.Identity()`,
			func(c dogma.ProjectionConfigurer) {
				c.ConsumesEventType(fixtures.MessageA{})
			},
		),
		Entry(
			"when the handler configures multiple identities",
			`*fixtures.ProjectionMessageHandler.Configure() has already called ProjectionConfigurer.Identity("<name>", "<key>")`,
			func(c dogma.ProjectionConfigurer) {
				c.Identity("<name>", "<key>")
				c.Identity("<other>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
			},
		),
		Entry(
			"when the handler configures an invalid name",
			`*fixtures.ProjectionMessageHandler.Configure() called ProjectionConfigurer.Identity() with an invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
			func(c dogma.ProjectionConfigurer) {
				c.Identity("\t \n", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
			},
		),
		Entry(
			"when the handler configures an invalid key",
			`*fixtures.ProjectionMessageHandler.Configure() called ProjectionConfigurer.Identity() with an invalid key "\t \n", keys must be non-empty, printable UTF-8 strings with no whitespace`,
			func(c dogma.ProjectionConfigurer) {
				c.Identity("<name>", "\t \n")
				c.ConsumesEventType(fixtures.MessageA{})
			},
		),
		Entry(
			"when the handler does not configure any consumed event types",
			`*fixtures.ProjectionMessageHandler.Configure() did not call ProjectionConfigurer.ConsumesEventType()`,
			func(c dogma.ProjectionConfigurer) {
				c.Identity("<name>", "<key>")
			},
		),
		Entry(
			"when the handler configures the same consumed event type multiple times",
			`*fixtures.ProjectionMessageHandler.Configure() has already called ProjectionConfigurer.ConsumesEventType(fixtures.MessageA)`,
			func(c dogma.ProjectionConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ConsumesEventType(fixtures.MessageA{})
			},
		),
	)
})
