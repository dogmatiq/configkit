package configkit

import (
	"context"
	"errors"
	"fmt"

	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/interopspec/configspec"
)

// ToProto converts an application configuration to its protocol buffers
// representation.
func ToProto(in Application) (*configspec.Application, error) {
	out := &configspec.Application{}

	var err error
	out.Identity, err = marshalIdentity(in.Identity())
	if err != nil {
		return nil, err
	}

	out.GoType = in.TypeName()
	if out.GoType == "" {
		return nil, errors.New("application type name is empty")
	}

	for _, hIn := range in.Handlers() {
		hOut, err := marshalHandler(hIn)
		if err != nil {
			return nil, err
		}

		out.Handlers = append(out.Handlers, hOut)
	}

	return out, nil
}

// FromProto converts an application configuration from its protocol buffers
// representation.
func FromProto(in *configspec.Application) (Application, error) {
	out := &unmarshaledApplication{}

	var err error
	out.ident, err = unmarshalIdentity(in.GetIdentity())
	if err != nil {
		return nil, err
	}

	out.typeName = in.GetGoType()
	if out.typeName == "" {
		return nil, errors.New("application type name is empty")
	}

	for _, hIn := range in.GetHandlers() {
		hOut, err := unmarshalHandler(hIn)
		if err != nil {
			return nil, err
		}

		if out.handlers == nil {
			out.handlers = HandlerSet{}
		}
		out.handlers.Add(hOut)

		for n, r := range hOut.MessageNames().Produced {
			if out.names.Produced == nil {
				out.names.Produced = message.NameRoles{}
			}
			out.names.Produced[n] = r
		}

		for n, r := range hOut.MessageNames().Consumed {
			if out.names.Consumed == nil {
				out.names.Consumed = message.NameRoles{}
			}
			out.names.Consumed[n] = r
		}
	}

	return out, nil
}

// marshalHandler marshals a handler config to its protobuf representation.
func marshalHandler(in Handler) (*configspec.Handler, error) {
	out := &configspec.Handler{
		IsDisabled: in.IsDisabled(),
	}

	var err error
	out.Identity, err = marshalIdentity(in.Identity())
	if err != nil {
		return nil, err
	}

	out.GoType = in.TypeName()
	if out.GoType == "" {
		return nil, errors.New("handler type name is empty")
	}

	out.Type, err = marshalHandlerType(in.HandlerType())
	if err != nil {
		return nil, err
	}

	names := in.MessageNames()
	out.ProducedMessages, err = marshalNameRoles(names.Produced)
	if err != nil {
		return nil, err
	}

	out.ConsumedMessages, err = marshalNameRoles(names.Consumed)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// unmarshalHandler unmarshals a handler configuration from its protocol buffers
// representation.
func unmarshalHandler(in *configspec.Handler) (Handler, error) {
	out := &unmarshaledHandler{
		isDisabled: in.GetIsDisabled(),
	}

	var err error
	out.ident, err = unmarshalIdentity(in.GetIdentity())
	if err != nil {
		return nil, err
	}

	out.typeName = in.GetGoType()
	if out.typeName == "" {
		return nil, errors.New("handler type name is empty")
	}

	out.handlerType, err = unmarshalHandlerType(in.GetType())
	if err != nil {
		return nil, err
	}

	out.names.Produced, err = unmarshalNameRoles(in.GetProducedMessages())
	if err != nil {
		return nil, err
	}

	out.names.Consumed, err = unmarshalNameRoles(in.GetConsumedMessages())
	if err != nil {
		return nil, err
	}

	return out, nil
}

// marshalNameRoles marshals a message.NameRoles collection into
// its protocol buffers representation.
func marshalNameRoles(in message.NameRoles) (map[string]configspec.MessageRole, error) {
	out := map[string]configspec.MessageRole{}

	for nIn, rIn := range in {
		nOut, err := nIn.MarshalText()
		if err != nil {
			return nil, err
		}

		rOut, err := marshalMessageRole(rIn)
		if err != nil {
			return nil, err
		}

		out[string(nOut)] = rOut
	}

	return out, nil
}

// unmarshalNameRoles unmarshals a message.NameRoles collection from
// its protocol buffers representation.
func unmarshalNameRoles(in map[string]configspec.MessageRole) (message.NameRoles, error) {
	out := message.NameRoles{}

	for nIn, rIn := range in {
		var nOut message.Name

		if err := nOut.UnmarshalText([]byte(nIn)); err != nil {
			return nil, err
		}

		rOut, err := unmarshalMessageRole(rIn)
		if err != nil {
			return nil, err
		}

		out[nOut] = rOut
	}

	return out, nil
}

// marshalIdentity marshals a Identity to its protocol buffers
// representation.
func marshalIdentity(in Identity) (*configspec.Identity, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	return &configspec.Identity{
		Name: in.Name,
		Key:  in.Key,
	}, nil
}

// unmarshalIdentity unmarshals a Identity from its protocol buffers
// representation.
func unmarshalIdentity(in *configspec.Identity) (Identity, error) {
	return NewIdentity(
		in.GetName(),
		in.GetKey(),
	)
}

// marshalHandlerType marshals a HandlerType to its protocol buffers
// representation.
func marshalHandlerType(t HandlerType) (configspec.HandlerType, error) {
	if err := t.Validate(); err != nil {
		return configspec.HandlerType_UNKNOWN_HANDLER_TYPE, err
	}

	switch t {
	case AggregateHandlerType:
		return configspec.HandlerType_AGGREGATE, nil
	case ProcessHandlerType:
		return configspec.HandlerType_PROCESS, nil
	case IntegrationHandlerType:
		return configspec.HandlerType_INTEGRATION, nil
	default: // ProjectionHandlerType
		return configspec.HandlerType_PROJECTION, nil
	}
}

// unmarshalHandlerType unmarshals a HandlerType from its protocol
// buffers representation.
func unmarshalHandlerType(t configspec.HandlerType) (HandlerType, error) {
	switch t {
	case configspec.HandlerType_AGGREGATE:
		return AggregateHandlerType, nil
	case configspec.HandlerType_PROCESS:
		return ProcessHandlerType, nil
	case configspec.HandlerType_INTEGRATION:
		return IntegrationHandlerType, nil
	case configspec.HandlerType_PROJECTION:
		return ProjectionHandlerType, nil
	default:
		return "", fmt.Errorf("unknown handler type: %#v", t)
	}
}

// marshalMessageRole marshals a message.Role to its protocol buffers
// representation.
func marshalMessageRole(r message.Role) (configspec.MessageRole, error) {
	if err := r.Validate(); err != nil {
		return configspec.MessageRole_UNKNOWN_MESSAGE_ROLE, err
	}

	switch r {
	case message.CommandRole:
		return configspec.MessageRole_COMMAND, nil
	case message.EventRole:
		return configspec.MessageRole_EVENT, nil
	default: // message.TimeoutRole
		return configspec.MessageRole_TIMEOUT, nil
	}
}

// unmarshalMessageRole unmarshals a message.Role from its protocol buffers
// representation.
func unmarshalMessageRole(r configspec.MessageRole) (message.Role, error) {
	switch r {
	case configspec.MessageRole_COMMAND:
		return message.CommandRole, nil
	case configspec.MessageRole_EVENT:
		return message.EventRole, nil
	case configspec.MessageRole_TIMEOUT:
		return message.TimeoutRole, nil
	default:
		return "", fmt.Errorf("unknown message role: %#v", r)
	}
}

// unmarshaledApplication is an implementation of [Application] that has been
// produced by unmarshaling a configuration.
type unmarshaledApplication struct {
	ident    Identity
	names    EntityMessageNames
	typeName string
	handlers HandlerSet
}

func (a *unmarshaledApplication) Identity() Identity {
	return a.ident
}

func (a *unmarshaledApplication) MessageNames() EntityMessageNames {
	return a.names
}

func (a *unmarshaledApplication) TypeName() string {
	return a.typeName
}

func (a *unmarshaledApplication) AcceptVisitor(ctx context.Context, v Visitor) error {
	return v.VisitApplication(ctx, a)
}

func (a *unmarshaledApplication) Handlers() HandlerSet {
	return a.handlers
}

// unmarshaledHandler is an implementation of [Handler] that has been produced
// by unmarshaling a configuration.
type unmarshaledHandler struct {
	ident       Identity
	names       EntityMessageNames
	typeName    string
	handlerType HandlerType
	isDisabled  bool
}

// Identity returns the identity of the entity.
func (h *unmarshaledHandler) Identity() Identity {
	return h.ident
}

// MessageNames returns information about the messages used by the entity.
func (h *unmarshaledHandler) MessageNames() EntityMessageNames {
	return h.names
}

// TypeName returns the fully-qualified type name of the entity.
func (h *unmarshaledHandler) TypeName() string {
	return h.typeName
}

// HandlerType returns the type of handler.
func (h *unmarshaledHandler) HandlerType() HandlerType {
	return h.handlerType
}

// IsDisabled returns true if the handler is disabled.
func (h *unmarshaledHandler) IsDisabled() bool {
	return h.isDisabled
}

// AcceptVisitor calls the appropriate method on v for this entity type.
func (h *unmarshaledHandler) AcceptVisitor(ctx context.Context, v Visitor) error {
	h.handlerType.MustValidate()

	switch h.handlerType {
	case AggregateHandlerType:
		return v.VisitAggregate(ctx, h)
	case ProcessHandlerType:
		return v.VisitProcess(ctx, h)
	case IntegrationHandlerType:
		return v.VisitIntegration(ctx, h)
	default: // ProjectionHandlerType
		return v.VisitProjection(ctx, h)
	}
}
