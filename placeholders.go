package configkit

// TypeName is a fully-qualified name for a Go type.
type TypeName string

// PLACEHOLDER TYPES
type MessageType interface{}
type HandlerType interface{} // enum AGGREGATE, PROCESS, INTEGRATION, PROJECTION
type MessageRole interface{} // enum COMMAND, EVENT, TIMEOUT
