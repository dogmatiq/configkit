package configkit

import "github.com/dogmatiq/configkit/message"

// HandlerSet is a collection of handlers.
type HandlerSet []Handler

// HandlerByIdentity by identity returns the handler with the given identity.
func (s HandlerSet) HandlerByIdentity(Identity) (Handler, bool) {
	panic("not implemented")
}

// HandlerByName by name returns the handler with the given name.
func (s HandlerSet) HandlerByName(string) (Handler, bool) {
	panic("not implemented")
}

// HandlerByKey by name returns the handler with the given key.
func (s HandlerSet) HandlerByKey(string) (Handler, bool) {
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

// RichHandlerSet is a collection of rich handlers
type RichHandlerSet []RichHandler

// HandlerByIdentity by identity returns the handler with the given identity.
func (s RichHandlerSet) HandlerByIdentity(Identity) (Handler, bool) {
	panic("not implemented")
}

// HandlerByName by name returns the handler with the given name.
func (s RichHandlerSet) HandlerByName(string) (Handler, bool) {
	panic("not implemented")
}

// HandlerByKey by name returns the handler with the given key.
func (s RichHandlerSet) HandlerByKey(string) (Handler, bool) {
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
