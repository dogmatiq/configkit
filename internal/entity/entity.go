package entity

import (
	"context"

	"github.com/dogmatiq/configkit"
)

// Application is an implementation of config.Application.
type Application struct {
	IdentityValue     configkit.Identity
	TypeNameValue     string
	MessageNamesValue configkit.EntityMessageNames
	HandlersValue     configkit.HandlerSet
}

// Identity returns the identity of the entity.
func (a *Application) Identity() configkit.Identity {
	return a.IdentityValue
}

// TypeName returns the fully-qualified type name of the entity.
func (a *Application) TypeName() string {
	return a.TypeNameValue
}

// MessageNames returns information about the messages used by the entity.
func (a *Application) MessageNames() configkit.EntityMessageNames {
	return a.MessageNamesValue
}

// Handlers returns the handlers within this application.
func (a *Application) Handlers() configkit.HandlerSet {
	return a.HandlersValue
}

// AcceptVisitor calls the appropriate method on v for this entity type.
func (a *Application) AcceptVisitor(ctx context.Context, v configkit.Visitor) error {
	return v.VisitApplication(ctx, a)
}

// Handler is an implementation of config.Handler.
type Handler struct {
	IdentityValue     configkit.Identity
	TypeNameValue     string
	MessageNamesValue configkit.EntityMessageNames
	HandlerTypeValue  configkit.HandlerType
	IsDisabledValue   bool
}

// Identity returns the identity of the entity.
func (h *Handler) Identity() configkit.Identity {
	return h.IdentityValue
}

// TypeName returns the fully-qualified type name of the entity.
func (h *Handler) TypeName() string {
	return h.TypeNameValue
}

// MessageNames returns information about the messages used by the entity.
func (h *Handler) MessageNames() configkit.EntityMessageNames {
	return h.MessageNamesValue
}

// HandlerType returns the type of handler.
func (h *Handler) HandlerType() configkit.HandlerType {
	return h.HandlerTypeValue
}

// IsDisabled returns true if the handler is disabled.
func (h *Handler) IsDisabled() bool {
	return h.IsDisabledValue
}

// AcceptVisitor calls the appropriate method on v for this entity type.
func (h *Handler) AcceptVisitor(ctx context.Context, v configkit.Visitor) error {
	h.HandlerTypeValue.MustValidate()

	switch h.HandlerTypeValue {
	case configkit.AggregateHandlerType:
		return v.VisitAggregate(ctx, h)
	case configkit.ProcessHandlerType:
		return v.VisitProcess(ctx, h)
	case configkit.IntegrationHandlerType:
		return v.VisitIntegration(ctx, h)
	default: // configkit.ProjectionHandlerType
		return v.VisitProjection(ctx, h)
	}
}
