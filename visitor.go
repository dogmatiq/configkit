package configkit

import "context"

// Visitor is a visitor that visits configurations.
type Visitor interface {
	VisitApplication(context.Context, Application) error
	VisitAggregate(context.Context, Aggregate) error
	VisitProcess(context.Context, Process) error
	VisitIntegration(context.Context, Integration) error
	VisitProjection(context.Context, Projection) error
}

// PortableVisitor is a visitor that visits "portable" configurations.
type PortableVisitor interface {
	VisitPortableApplication(context.Context, PortableApplication) error
	VisitPortableAggregate(context.Context, PortableAggregate) error
	VisitPortableProcess(context.Context, PortableProcess) error
	VisitPortableIntegration(context.Context, PortableIntegration) error
	VisitPortableProjection(context.Context, PortableProjection) error
}

// RichVisitor is a visitor that visits "rich" configurations.
type RichVisitor interface {
	VisitRichApplication(context.Context, RichApplication) error
	VisitRichAggregate(context.Context, RichAggregate) error
	VisitRichProcess(context.Context, RichProcess) error
	VisitRichIntegration(context.Context, RichIntegration) error
	VisitRichProjection(context.Context, RichProjection) error
}
