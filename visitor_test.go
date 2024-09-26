package configkit_test

import (
	"context"

	. "github.com/dogmatiq/configkit"
)

// visitorStub is a test implementation of [configkit.Visitor].
type visitorStub struct {
	VisitApplicationFunc func(context.Context, Application) error
	VisitAggregateFunc   func(context.Context, Aggregate) error
	VisitProcessFunc     func(context.Context, Process) error
	VisitIntegrationFunc func(context.Context, Integration) error
	VisitProjectionFunc  func(context.Context, Projection) error
}

func (v *visitorStub) VisitApplication(ctx context.Context, cfg Application) error {
	if v.VisitApplicationFunc != nil {
		return v.VisitApplicationFunc(ctx, cfg)
	}
	return cfg.Handlers().AcceptVisitor(ctx, v)
}

func (v *visitorStub) VisitAggregate(ctx context.Context, cfg Aggregate) error {
	if v.VisitAggregateFunc == nil {
		return nil
	}
	return v.VisitAggregateFunc(ctx, cfg)
}

func (v *visitorStub) VisitProcess(ctx context.Context, cfg Process) error {
	if v.VisitProcessFunc == nil {
		return nil
	}
	return v.VisitProcessFunc(ctx, cfg)
}

func (v *visitorStub) VisitIntegration(ctx context.Context, cfg Integration) error {
	if v.VisitIntegrationFunc == nil {
		return nil
	}
	return v.VisitIntegrationFunc(ctx, cfg)
}

func (v *visitorStub) VisitProjection(ctx context.Context, cfg Projection) error {
	if v.VisitProjectionFunc == nil {
		return nil
	}
	return v.VisitProjectionFunc(ctx, cfg)
}

// richVisitorStub is a test implementation of [configkit.RichVisitor].
type richVisitorStub struct {
	VisitRichApplicationFunc func(context.Context, RichApplication) error
	VisitRichAggregateFunc   func(context.Context, RichAggregate) error
	VisitRichProcessFunc     func(context.Context, RichProcess) error
	VisitRichIntegrationFunc func(context.Context, RichIntegration) error
	VisitRichProjectionFunc  func(context.Context, RichProjection) error
}

func (v *richVisitorStub) VisitRichApplication(ctx context.Context, cfg RichApplication) error {
	if v.VisitRichApplicationFunc != nil {
		return v.VisitRichApplicationFunc(ctx, cfg)
	}
	return cfg.RichHandlers().AcceptRichVisitor(ctx, v)
}

func (v *richVisitorStub) VisitRichAggregate(ctx context.Context, cfg RichAggregate) error {
	if v.VisitRichAggregateFunc == nil {
		return nil
	}
	return v.VisitRichAggregateFunc(ctx, cfg)
}

func (v *richVisitorStub) VisitRichProcess(ctx context.Context, cfg RichProcess) error {
	if v.VisitRichProcessFunc == nil {
		return nil
	}
	return v.VisitRichProcessFunc(ctx, cfg)
}

func (v *richVisitorStub) VisitRichIntegration(ctx context.Context, cfg RichIntegration) error {
	if v.VisitRichIntegrationFunc == nil {
		return nil
	}
	return v.VisitRichIntegrationFunc(ctx, cfg)
}

func (v *richVisitorStub) VisitRichProjection(ctx context.Context, cfg RichProjection) error {
	if v.VisitRichProjectionFunc == nil {
		return nil
	}
	return v.VisitRichProjectionFunc(ctx, cfg)
}
