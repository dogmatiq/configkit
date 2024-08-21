package entity_test

import (
	"context"
	"errors"
	"reflect"

	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/fixtures"
	cfixtures "github.com/dogmatiq/configkit/fixtures"
	. "github.com/dogmatiq/configkit/internal/entity"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures" // can't dot-import due to conflicts
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Application", func() {
	var app *Application

	BeforeEach(func() {
		app = &Application{}
	})

	Describe("func Identity()", func() {
		It("returns the value of .IdentityValue field", func() {
			i := configkit.MustNewIdentity("<name>", "e68e772b-0f6f-49c8-b882-f0dc997bd6a1")
			app.IdentityValue = i
			Expect(app.Identity()).To(Equal(i))
		})
	})

	Describe("func TypeName()", func() {
		It("returns the value of .TypeNameValue field", func() {
			app.TypeNameValue = reflect.TypeOf(app).String()
			Expect(app.TypeName()).To(Equal("*entity.Application"))
		})
	})

	Describe("func MessageNames()", func() {
		It("returns the value of .MessageNamesValue field", func() {
			m := configkit.EntityMessageNames{
				Produced: message.NameRoles{
					cfixtures.MessageCTypeName: message.CommandRole,
				},
				Consumed: message.NameRoles{
					cfixtures.MessageETypeName: message.EventRole,
				},
			}

			app.MessageNamesValue = m
			Expect(app.MessageNames()).To(Equal(m))
		})
	})

	Describe("func AcceptVisitor()", func() {
		It("calls the correct visitor method", func() {
			v := &Visitor{
				VisitApplicationFunc: func(_ context.Context, a configkit.Application) error {
					Expect(a).To(Equal(app))
					return errors.New("<error>")
				},
			}

			err := app.AcceptVisitor(context.Background(), v)
			Expect(err).To(Equal(err))
		})
	})

	Describe("func Handlers()", func() {
		It("returns the value of .HandlersValue field", func() {
			aggregate := configkit.FromAggregate(&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<agg-name>", "63990c32-ecdd-46dd-8e6a-7bb16f3b1730")
					c.Routes(
						dogma.HandlesCommand[fixtures.MessageC](),
						dogma.RecordsEvent[fixtures.MessageE](),
					)
				},
			})

			projection := configkit.FromProjection(&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("<proj-name>", "b34181e8-2930-4b6c-a649-18a001836ec3")
					c.Routes(
						dogma.HandlesEvent[fixtures.MessageE](),
					)
				},
			})

			s := configkit.NewHandlerSet(aggregate, projection)

			app.HandlersValue = s
			Expect(app.Handlers()).To(Equal(s))

		})
	})
})

var _ = Describe("type Handler", func() {
	var handler *Handler

	BeforeEach(func() {
		handler = &Handler{}
	})

	Describe("func Identity()", func() {
		It("returns the value of .IdentityValue field", func() {
			i := configkit.MustNewIdentity("<name>", "7502b306-8e9e-4f23-99ba-00c5bf138de5")
			handler.IdentityValue = i
			Expect(handler.Identity()).To(Equal(i))
		})
	})

	Describe("func TypeName()", func() {
		It("returns the value of .TypeNameValue field", func() {
			handler.TypeNameValue = reflect.TypeOf(handler).String()
			Expect(handler.TypeName()).To(Equal("*entity.Handler"))
		})
	})

	Describe("func MessageNames()", func() {
		It("returns the value of .MessageNamesValue field", func() {
			m := configkit.EntityMessageNames{
				Produced: message.NameRoles{
					cfixtures.MessageCTypeName: message.CommandRole,
				},
				Consumed: message.NameRoles{
					cfixtures.MessageETypeName: message.EventRole,
				},
			}

			handler.MessageNamesValue = m
			Expect(handler.MessageNames()).To(Equal(m))
		})
	})

	Describe("func HandlerType()", func() {
		It("returns the value of .HandlerTypeValue field", func() {
			handler.HandlerTypeValue = configkit.AggregateHandlerType
			Expect(handler.HandlerType()).To(Equal(configkit.AggregateHandlerType))
		})
	})

	Describe("func IsDisabled()", func() {
		It("returns the value of .IsDisabledValue field", func() {
			handler.IsDisabledValue = true
			Expect(handler.IsDisabled()).To(BeTrue())
		})
	})

	Describe("func AcceptVisitor()", func() {
		var visitor *Visitor

		BeforeEach(func() {
			visitor = &Visitor{}
		})

		When("the handler is an aggregate", func() {
			BeforeEach(func() {
				handler.HandlerTypeValue = configkit.AggregateHandlerType
			})

			It("calls the correct visitor method", func() {
				visitor.VisitAggregateFunc = func(_ context.Context, h configkit.Aggregate) error {
					Expect(h).To(Equal(handler))
					return errors.New("<error>")
				}

				err := handler.AcceptVisitor(context.Background(), visitor)
				Expect(err).To(Equal(err))
			})
		})

		When("the handler is a process", func() {
			BeforeEach(func() {
				handler.HandlerTypeValue = configkit.ProcessHandlerType
			})

			It("calls the correct visitor method", func() {
				visitor.VisitProcessFunc = func(_ context.Context, h configkit.Process) error {
					Expect(h).To(Equal(handler))
					return errors.New("<error>")
				}

				err := handler.AcceptVisitor(context.Background(), visitor)
				Expect(err).To(Equal(err))
			})
		})

		When("the handler is an integration", func() {
			BeforeEach(func() {
				handler.HandlerTypeValue = configkit.IntegrationHandlerType
			})

			It("calls the correct visitor method", func() {
				visitor.VisitIntegrationFunc = func(_ context.Context, h configkit.Integration) error {
					Expect(h).To(Equal(handler))
					return errors.New("<error>")
				}

				err := handler.AcceptVisitor(context.Background(), visitor)
				Expect(err).To(Equal(err))
			})
		})

		When("the handler is a projection", func() {
			BeforeEach(func() {
				handler.HandlerTypeValue = configkit.ProjectionHandlerType
			})

			It("calls the correct visitor method", func() {
				visitor.VisitProjectionFunc = func(_ context.Context, h configkit.Projection) error {
					Expect(h).To(Equal(handler))
					return errors.New("<error>")
				}

				err := handler.AcceptVisitor(context.Background(), visitor)
				Expect(err).To(Equal(err))
			})
		})
	})
})
