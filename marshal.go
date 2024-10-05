package configkit

import (
	"context"
	"errors"
	"fmt"

	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/protobuf/configpb"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

// ToProto converts an application configuration to its protocol buffers
// representation.
func ToProto(app Application) (*configpb.Application, error) {
	out := &configpb.Application{}

	var err error
	out.Identity, err = marshalIdentity(app.Identity())
	if err != nil {
		return nil, err
	}

	out.GoType = app.TypeName()
	if out.GoType == "" {
		return nil, errors.New("application type name is empty")
	}

	for n, em := range app.MessageNames() {
		kOut, err := marshalMessageKind(em.Kind)
		if err != nil {
			return nil, err
		}

		if out.Messages == nil {
			out.Messages = map[string]configpb.MessageKind{}
		}
		out.Messages[string(n)] = kOut
	}

	for _, h := range app.Handlers() {
		handlerOut, err := marshalHandler(h)
		if err != nil {
			return nil, err
		}
		out.Handlers = append(out.Handlers, handlerOut)
	}

	return out, nil
}

// FromProto converts an application configuration from its protocol buffers
// representation.
func FromProto(app *configpb.Application) (Application, error) {
	out := &unmarshaledApplication{}

	var err error
	out.ident, err = unmarshalIdentity(app.GetIdentity())
	if err != nil {
		return nil, err
	}

	out.typeName = app.GetGoType()
	if out.typeName == "" {
		return nil, errors.New("application type name is empty")
	}

	kinds := map[message.Name]message.Kind{}

	for n, k := range app.GetMessages() {
		kOut, err := unmarshalMessageKind(k)
		if err != nil {
			return nil, err
		}

		kinds[message.Name(n)] = kOut
	}

	for _, h := range app.GetHandlers() {
		handlerOut, err := unmarshalHandler(h, kinds)
		if err != nil {
			return nil, err
		}

		if out.handlers == nil {
			out.handlers = HandlerSet{}
		}
		out.handlers.Add(handlerOut)
	}

	return out, nil
}

// marshalHandler marshals a handler config to its protobuf representation.
func marshalHandler(in Handler) (*configpb.Handler, error) {
	out := &configpb.Handler{
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

	for n, em := range in.MessageNames() {
		if n == "" {
			return nil, errors.New("message name is empty")
		}

		if out.Messages == nil {
			out.Messages = map[string]*configpb.MessageUsage{}
		}

		key := string(n)

		usage, ok := out.Messages[key]
		if !ok {
			usage = &configpb.MessageUsage{}
			out.Messages[key] = usage
		}

		if em.IsConsumed {
			usage.IsConsumed = true
		}

		if em.IsProduced {
			usage.IsProduced = true
		}
	}

	return out, nil
}

// unmarshalHandler unmarshals a handler configuration from its protocol buffers
// representation.
func unmarshalHandler(
	in *configpb.Handler,
	kinds map[message.Name]message.Kind,
) (Handler, error) {
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

	for n, usage := range in.GetMessages() {
		nOut := message.Name(n)
		kOut, ok := kinds[nOut]
		if !ok {
			return nil, fmt.Errorf("message name %s as no associated message kind", n)
		}

		if out.names == nil {
			out.names = EntityMessages[message.Name]{}
		}

		out.names.Update(
			nOut,
			func(n message.Name, em *EntityMessage) {
				em.Kind = kOut

				if usage.IsProduced {
					em.IsProduced = true
				}

				if usage.IsConsumed {
					em.IsConsumed = true
				}
			},
		)
	}

	return out, nil
}

// marshalIdentity marshals a Identity to its protocol buffers
// representation.
func marshalIdentity(in Identity) (*identitypb.Identity, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	return &identitypb.Identity{
		Name: in.Name,
		Key:  uuidpb.MustParse(in.Key),
	}, nil
}

// unmarshalIdentity unmarshals a Identity from its protocol buffers
// representation.
func unmarshalIdentity(in *identitypb.Identity) (Identity, error) {
	if err := ValidateIdentityName(in.GetName()); err != nil {
		return Identity{}, err
	}

	return Identity{
		Name: in.GetName(),
		Key:  in.GetKey().AsString(),
	}, nil
}

// marshalHandlerType marshals a HandlerType to its protocol buffers
// representation.
func marshalHandlerType(t HandlerType) (configpb.HandlerType, error) {
	if err := t.Validate(); err != nil {
		return configpb.HandlerType_UNKNOWN_HANDLER_TYPE, err
	}

	switch t {
	case AggregateHandlerType:
		return configpb.HandlerType_AGGREGATE, nil
	case ProcessHandlerType:
		return configpb.HandlerType_PROCESS, nil
	case IntegrationHandlerType:
		return configpb.HandlerType_INTEGRATION, nil
	default: // ProjectionHandlerType
		return configpb.HandlerType_PROJECTION, nil
	}
}

// unmarshalHandlerType unmarshals a HandlerType from its protocol
// buffers representation.
func unmarshalHandlerType(t configpb.HandlerType) (HandlerType, error) {
	switch t {
	case configpb.HandlerType_AGGREGATE:
		return AggregateHandlerType, nil
	case configpb.HandlerType_PROCESS:
		return ProcessHandlerType, nil
	case configpb.HandlerType_INTEGRATION:
		return IntegrationHandlerType, nil
	case configpb.HandlerType_PROJECTION:
		return ProjectionHandlerType, nil
	default:
		return "", fmt.Errorf("unknown handler type: %#v", t)
	}
}

// marshalMessageKind marshals a [message.Kind] to its protocol buffers
// representation.
func marshalMessageKind(k message.Kind) (configpb.MessageKind, error) {
	switch k {
	case message.CommandKind:
		return configpb.MessageKind_COMMAND, nil
	case message.EventKind:
		return configpb.MessageKind_EVENT, nil
	case message.TimeoutKind:
		return configpb.MessageKind_TIMEOUT, nil
	default:
		return configpb.MessageKind_UNKNOWN_MESSAGE_KIND, fmt.Errorf("unknown message kind: %#v", k)
	}
}

// unmarshalMessageKind unmarshals a [message.Kind] from its protocol buffers
// representation.
func unmarshalMessageKind(r configpb.MessageKind) (message.Kind, error) {
	switch r {
	case configpb.MessageKind_COMMAND:
		return message.CommandKind, nil
	case configpb.MessageKind_EVENT:
		return message.EventKind, nil
	case configpb.MessageKind_TIMEOUT:
		return message.TimeoutKind, nil
	default:
		return 0, fmt.Errorf("unknown message kind: %#v", r)
	}
}

// unmarshaledApplication is an implementation of [Application] that has been
// produced by unmarshaling a configuration.
type unmarshaledApplication struct {
	ident    Identity
	typeName string
	handlers HandlerSet
}

func (a *unmarshaledApplication) Identity() Identity {
	return a.ident
}

func (a *unmarshaledApplication) MessageNames() EntityMessages[message.Name] {
	names := EntityMessages[message.Name]{}

	for _, h := range a.handlers {
		names.merge(h.MessageNames())
	}

	return names
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
	names       EntityMessages[message.Name]
	typeName    string
	handlerType HandlerType
	isDisabled  bool
}

// Identity returns the identity of the entity.
func (h *unmarshaledHandler) Identity() Identity {
	return h.ident
}

// MessageNames returns information about the messages used by the entity.
func (h *unmarshaledHandler) MessageNames() EntityMessages[message.Name] {
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
