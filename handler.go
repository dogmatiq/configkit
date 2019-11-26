package configkit

import (
	"reflect"

	"github.com/dogmatiq/configkit/internal/typename"
)

// Handler is a specialization of the Entity interface for message handlers.
type Handler interface {
	Entity

	// HandlerType returns the type of handler.
	HandlerType() HandlerType
}

// RichHandler is a specialization of the Handler interface that exposes
// information about the Go types used to implement the Dogma application.
type RichHandler interface {
	RichEntity

	// HandlerType returns the type of handler.
	HandlerType() HandlerType
}

// handler is a partial implementation of RichHandler.
type handler struct {
	rt reflect.Type
	ht HandlerType

	ident Identity
	names EntityMessageNames
	types EntityMessageTypes
}

func (h *handler) Identity() Identity {
	return h.ident
}

func (h *handler) MessageNames() EntityMessageNames {
	return h.names
}

func (h *handler) MessageTypes() EntityMessageTypes {
	return h.types
}

func (h *handler) TypeName() string {
	return typename.Of(h.ReflectType())
}

func (h *handler) ReflectType() reflect.Type {
	return h.rt
}

func (h *aggregate) HandlerType() HandlerType {
	return h.ht
}
