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
				c.Identity("<name>", projectionKey)
				c.Routes(
					dogma.HandlesEvent[fixtures.MessageA](),
					dogma.HandlesEvent[fixtures.MessageB](),
				)
				c.DeliveryPolicy(dogma.BroadcastProjectionDeliveryPolicy{
					PrimaryFirst: true,
				})
			},
		}
	})

	When("the configuration is valid", func() {
		var cfg RichProjection

		JustBeforeEach(func() {
			cfg = FromProjection(handler)
		})

		Describe("func Identity()", func() {
			It("returns the handler identity", func() {
				Expect(cfg.Identity()).To(Equal(
					MustNewIdentity("<name>", projectionKey),
				))
			})
		})

		Describe("func MessageNames()", func() {
			It("returns the expected message names", func() {
				Expect(cfg.MessageNames()).To(Equal(
					EntityMessageNames{
						Produced: nil,
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.EventRole,
							cfixtures.MessageBTypeName: message.EventRole,
						},
					},
				))
			})
		})

		Describe("func MessageTypes()", func() {
			It("returns the expected message types", func() {
				Expect(cfg.MessageTypes()).To(Equal(
					EntityMessageTypes{
						Produced: nil,
						Consumed: message.TypeRoles{
							cfixtures.MessageAType: message.EventRole,
							cfixtures.MessageBType: message.EventRole,
						},
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

		Describe("func DeliveryPolicy()", func() {
			It("returns the delivery policy", func() {
				Expect(cfg.DeliveryPolicy()).To(Equal(
					dogma.BroadcastProjectionDeliveryPolicy{
						PrimaryFirst: true,
					},
				))
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

		When("the handler is disabled", func() {
			BeforeEach(func() {
				configure := handler.ConfigureFunc
				handler.ConfigureFunc = func(c dogma.ProjectionConfigurer) {
					configure(c)
					c.Disable()
				}
			})

			Describe("func IsDisabled()", func() {
				It("returns true", func() {
					Expect(cfg.IsDisabled()).To(BeTrue())
				})
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
			`*fixtures.ProjectionMessageHandler is configured without an identity, Identity() must be called exactly once within Configure()`,
			func(c dogma.ProjectionConfigurer) {
				c.Routes(
					dogma.HandlesEvent[fixtures.MessageA](),
				)
			},
		),
		Entry(
			"when the handler configures multiple identities",
			`*fixtures.ProjectionMessageHandler is configured with multiple identities (<name>/70fdf7fa-4b24-448d-bd29-7ecc71d18c56 and <other>/70fdf7fa-4b24-448d-bd29-7ecc71d18c56), Identity() must be called exactly once within Configure()`,
			func(c dogma.ProjectionConfigurer) {
				c.Identity("<name>", projectionKey)
				c.Identity("<other>", projectionKey)
				c.Routes(
					dogma.HandlesEvent[fixtures.MessageA](),
				)
			},
		),
		Entry(
			"when the handler configures an invalid name",
			`*fixtures.ProjectionMessageHandler is configured with an invalid identity, invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
			func(c dogma.ProjectionConfigurer) {
				c.Identity("\t \n", projectionKey)
				c.Routes(
					dogma.HandlesEvent[fixtures.MessageA](),
				)
			},
		),
		Entry(
			"when the handler configures an invalid key",
			`*fixtures.ProjectionMessageHandler is configured with an invalid identity, invalid key "\t \n", keys must be RFC 4122 UUIDs`,
			func(c dogma.ProjectionConfigurer) {
				c.Identity("<name>", "\t \n")
				c.Routes(
					dogma.HandlesEvent[fixtures.MessageA](),
				)
			},
		),
		Entry(
			"when the handler does not configure any event routes",
			`*fixtures.ProjectionMessageHandler (<name>) is not configured to handle any events, at least one HandlesEvent() route must be added within Configure()`,
			func(c dogma.ProjectionConfigurer) {
				c.Identity("<name>", projectionKey)
			},
		),
		Entry(
			"when the handler configures multiple routes for the same event",
			`*fixtures.ProjectionMessageHandler (<name>) is configured with multiple HandlesEvent() routes for fixtures.MessageA, should these refer to different message types?`,
			func(c dogma.ProjectionConfigurer) {
				c.Identity("<name>", projectionKey)
				c.Routes(
					dogma.HandlesEvent[fixtures.MessageA](),
					dogma.HandlesEvent[fixtures.MessageA](),
				)
			},
		),
		Entry(
			"when an error occurs before the identity is configured it omits the handler name",
			`*fixtures.ProjectionMessageHandler is configured with multiple HandlesEvent() routes for fixtures.MessageA, should these refer to different message types?`,
			func(c dogma.ProjectionConfigurer) {
				c.Routes(
					dogma.HandlesEvent[fixtures.MessageA](),
					dogma.HandlesEvent[fixtures.MessageA](),
				)
			},
		),
		Entry(
			"when the handler configures a nil delivery policy",
			`delivery policy must not be nil`,
			func(c dogma.ProjectionConfigurer) {
				c.Identity("<name>", projectionKey)
				c.Routes(
					dogma.HandlesEvent[fixtures.MessageA](),
				)
				c.DeliveryPolicy(nil)
			},
		),
	)
})
