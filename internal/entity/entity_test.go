package entity

import (
	"context"
	"errors"

	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type application", func() {
	var app *Application

	BeforeEach(func() {
		app = &Application{}
	})

	Describe("func Accept()", func() {
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
})

var _ = Describe("type handler", func() {
	var hnd *Handler

	BeforeEach(func() {
		hnd = &Handler{}
	})

	Describe("func Accept()", func() {
		var visitor *Visitor

		BeforeEach(func() {
			visitor = &Visitor{}
		})

		When("the handler is an aggregate", func() {
			BeforeEach(func() {
				hnd.HandlerTypeValue = configkit.AggregateHandlerType
			})

			It("calls the correct visitor method", func() {
				visitor.VisitAggregateFunc = func(_ context.Context, h configkit.Aggregate) error {
					Expect(h).To(Equal(hnd))
					return errors.New("<error>")
				}

				err := hnd.AcceptVisitor(context.Background(), visitor)
				Expect(err).To(Equal(err))
			})
		})

		When("the handler is a process", func() {
			BeforeEach(func() {
				hnd.HandlerTypeValue = configkit.ProcessHandlerType
			})

			It("calls the correct visitor method", func() {
				visitor.VisitProcessFunc = func(_ context.Context, h configkit.Process) error {
					Expect(h).To(Equal(hnd))
					return errors.New("<error>")
				}

				err := hnd.AcceptVisitor(context.Background(), visitor)
				Expect(err).To(Equal(err))
			})
		})

		When("the handler is an integration", func() {
			BeforeEach(func() {
				hnd.HandlerTypeValue = configkit.IntegrationHandlerType
			})

			It("calls the correct visitor method", func() {
				visitor.VisitIntegrationFunc = func(_ context.Context, h configkit.Integration) error {
					Expect(h).To(Equal(hnd))
					return errors.New("<error>")
				}

				err := hnd.AcceptVisitor(context.Background(), visitor)
				Expect(err).To(Equal(err))
			})
		})

		When("the handler is a projection", func() {
			BeforeEach(func() {
				hnd.HandlerTypeValue = configkit.ProjectionHandlerType
			})

			It("calls the correct visitor method", func() {
				visitor.VisitProjectionFunc = func(_ context.Context, h configkit.Projection) error {
					Expect(h).To(Equal(hnd))
					return errors.New("<error>")
				}

				err := hnd.AcceptVisitor(context.Background(), visitor)
				Expect(err).To(Equal(err))
			})
		})
	})
})
