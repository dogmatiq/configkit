package configkit

import "github.com/dogmatiq/configkit/message"

// HandlerSet is a collection of handlers.
type HandlerSet []Handler

// ByIdentity returns the handler with the given identity.
func (s HandlerSet) ByIdentity(Identity) (Handler, bool) {
	panic("not implemented")
}

// ByName returns the handler with the given name.
func (s HandlerSet) ByName(string) (Handler, bool) {
	panic("not implemented")
}

// ByKey returns the handler with the given key.
func (s HandlerSet) ByKey(string) (Handler, bool) {
	panic("not implemented")
}

// ConsumersOf returns the handlers that consume messages with the given name.
func (s HandlerSet) ConsumersOf(message.Name) HandlerSet {
	panic("not implemented")
}

// ProducersOf returns the handlers that produce messages with the given name.
func (s HandlerSet) ProducersOf(message.Name) HandlerSet {
	panic("not implemented")
}

// RichHandlerSet is a collection of rich handlers.
type RichHandlerSet []RichHandler

// ByIdentity returns the handler with the given identity.
func (s RichHandlerSet) ByIdentity(Identity) (Handler, bool) {
	panic("not implemented")
}

// ByName returns the handler with the given name.
func (s RichHandlerSet) ByName(string) (Handler, bool) {
	panic("not implemented")
}

// ByKey returns the handler with the given key.
func (s RichHandlerSet) ByKey(string) (Handler, bool) {
	panic("not implemented")
}

// ConsumersOf returns the handlers that consume messages of the given type.
func (s RichHandlerSet) ConsumersOf(message.Type) RichHandlerSet {
	panic("not implemented")
}

// ProducersOf returns the handlers that produce messages of the given type.
func (s RichHandlerSet) ProducersOf(message.Type) RichHandlerSet {
	panic("not implemented")
}
