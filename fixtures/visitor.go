package fixtures

import (
	"context"

	"github.com/dogmatiq/configkit"
)

var _ configkit.Visitor = (*Visitor)(nil)

// Visitor is an mock implementation of configkit.Visitor.
type Visitor struct {
	VisitApplicationFunc func(context.Context, configkit.Application) error
	VisitAggregateFunc   func(context.Context, configkit.Aggregate) error
	VisitProcessFunc     func(context.Context, configkit.Process) error
	VisitIntegrationFunc func(context.Context, configkit.Integration) error
	VisitProjectionFunc  func(context.Context, configkit.Projection) error
}

// VisitApplication calls v.VisitApplicationFunc(ctx, cfg) if it is non-nil.
//
// If v.VisitApplicationFunc is nil, each of the handlers in cfg is visited.
func (v *Visitor) VisitApplication(ctx context.Context, cfg configkit.Application) error {
	if v.VisitApplicationFunc != nil {
		return v.VisitApplicationFunc(ctx, cfg)
	}

	for _, h := range cfg.Handlers() {
		if err := h.AcceptVisitor(ctx, v); err != nil {
			return err
		}
	}

	return nil
}

// VisitAggregate calls v.VisitAggregateFunc(ctx, cfg) if it is non-nil.
func (v *Visitor) VisitAggregate(ctx context.Context, cfg configkit.Aggregate) error {
	if v.VisitAggregateFunc == nil {
		return nil
	}

	return v.VisitAggregateFunc(ctx, cfg)
}

// VisitProcess calls v.VisitProcessFunc(ctx, cfg) if it is non-nil.
func (v *Visitor) VisitProcess(ctx context.Context, cfg configkit.Process) error {
	if v.VisitProcessFunc == nil {
		return nil
	}

	return v.VisitProcessFunc(ctx, cfg)
}

// VisitIntegration calls v.VisitIntegrationFunc(ctx, cfg) if it is non-nil.
func (v *Visitor) VisitIntegration(ctx context.Context, cfg configkit.Integration) error {
	if v.VisitIntegrationFunc == nil {
		return nil
	}

	return v.VisitIntegrationFunc(ctx, cfg)
}

// VisitProjection calls v.VisitProjectionFunc(ctx, cfg) if it is non-nil.
func (v *Visitor) VisitProjection(ctx context.Context, cfg configkit.Projection) error {
	if v.VisitProjectionFunc == nil {
		return nil
	}

	return v.VisitProjectionFunc(ctx, cfg)
}

var _ configkit.RichVisitor = (*RichVisitor)(nil)

// RichVisitor is an mock implementation of configkit.RichVisitor.
type RichVisitor struct {
	VisitRichApplicationFunc func(context.Context, *configkit.RichApplication) error
	VisitRichAggregateFunc   func(context.Context, *configkit.RichAggregate) error
	VisitRichProcessFunc     func(context.Context, *configkit.RichProcess) error
	VisitRichIntegrationFunc func(context.Context, *configkit.RichIntegration) error
	VisitRichProjectionFunc  func(context.Context, *configkit.RichProjection) error
}

// VisitRichApplication calls v.VisitApplicationFunc(ctx, cfg) if it is non-nil.
//
// If v.VisitRichApplicationFunc is nil, each of the handlers in cfg is visited.
func (v *RichVisitor) VisitRichApplication(ctx context.Context, cfg *configkit.RichApplication) error {
	if v.VisitRichApplicationFunc != nil {
		return v.VisitRichApplicationFunc(ctx, cfg)
	}

	for _, h := range cfg.RichHandlers() {
		if err := h.AcceptRichVisitor(ctx, v); err != nil {
			return err
		}
	}

	return nil
}

// VisitRichAggregate calls v.VisitAggregateFunc(ctx, cfg) if it is non-nil.
func (v *RichVisitor) VisitRichAggregate(ctx context.Context, cfg *configkit.RichAggregate) error {
	if v.VisitRichAggregateFunc == nil {
		return nil
	}

	return v.VisitRichAggregateFunc(ctx, cfg)
}

// VisitRichProcess calls v.VisitProcessFunc(ctx, cfg) if it is non-nil.
func (v *RichVisitor) VisitRichProcess(ctx context.Context, cfg *configkit.RichProcess) error {
	if v.VisitRichProcessFunc == nil {
		return nil
	}

	return v.VisitRichProcessFunc(ctx, cfg)
}

// VisitRichIntegration calls v.VisitIntegrationFunc(ctx, cfg) if it is non-nil.
func (v *RichVisitor) VisitRichIntegration(ctx context.Context, cfg *configkit.RichIntegration) error {
	if v.VisitRichIntegrationFunc == nil {
		return nil
	}

	return v.VisitRichIntegrationFunc(ctx, cfg)
}

// VisitRichProjection calls v.VisitProjectionFunc(ctx, cfg) if it is non-nil.
func (v *RichVisitor) VisitRichProjection(ctx context.Context, cfg *configkit.RichProjection) error {
	if v.VisitRichProjectionFunc == nil {
		return nil
	}

	return v.VisitRichProjectionFunc(ctx, cfg)
}
