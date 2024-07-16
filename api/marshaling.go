package api

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/internal/entity"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/interopspec/configspec"
)

// marshalApplication marshals an application config to its protobuf
// representation.
func marshalApplication(in configkit.Application) (*configspec.Application, error) {
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

// unmarshalApplication unmarshals an application config from its protobuf
// representation.
func unmarshalApplication(in *configspec.Application) (configkit.Application, error) {
	out := &entity.Application{
		MessageNamesValue: configkit.EntityMessageNames{
			Produced: message.NameRoles{},
			Consumed: message.NameRoles{},
		},
		HandlersValue: configkit.HandlerSet{},
	}

	var err error
	out.IdentityValue, err = unmarshalIdentity(in.GetIdentity())
	if err != nil {
		return nil, err
	}

	out.TypeNameValue = in.GetGoType()
	if out.TypeNameValue == "" {
		return nil, errors.New("application type name is empty")
	}

	for _, hIn := range in.GetHandlers() {
		hOut, err := unmarshalHandler(hIn)
		if err != nil {
			return nil, err
		}

		out.HandlersValue.Add(hOut)

		for n, r := range hOut.MessageNames().Produced {
			out.MessageNamesValue.Produced[n] = r
		}

		for n, r := range hOut.MessageNames().Consumed {
			out.MessageNamesValue.Consumed[n] = r
		}
	}

	return out, nil
}

// marshalHandler marshals a handler config to its protobuf representation.
func marshalHandler(in configkit.Handler) (*configspec.Handler, error) {
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
func unmarshalHandler(in *configspec.Handler) (configkit.Handler, error) {
	out := &entity.Handler{
		MessageNamesValue: configkit.EntityMessageNames{
			Produced: message.NameRoles{},
			Consumed: message.NameRoles{},
		},
		IsDisabledValue: in.GetIsDisabled(),
	}

	var err error
	out.IdentityValue, err = unmarshalIdentity(in.GetIdentity())
	if err != nil {
		return nil, err
	}

	out.TypeNameValue = in.GetGoType()
	if out.TypeNameValue == "" {
		return nil, errors.New("handler type name is empty")
	}

	out.HandlerTypeValue, err = unmarshalHandlerType(in.GetType())
	if err != nil {
		return nil, err
	}

	out.MessageNamesValue.Produced, err = unmarshalNameRoles(in.GetProducedMessages())
	if err != nil {
		return nil, err
	}

	out.MessageNamesValue.Consumed, err = unmarshalNameRoles(in.GetConsumedMessages())
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

// marshalNameRoles unmarshals a message.NameRoles collection from
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

// marshalIdentity marshals a configkit.Identity to its protocol buffers
// representation.
func marshalIdentity(in configkit.Identity) (*configspec.Identity, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	return &configspec.Identity{
		Name: in.Name,
		Key:  in.Key,
	}, nil
}

// unmarshalIdentity unmarshals a configkit.Identity from its protocol buffers
// representation.
func unmarshalIdentity(in *configspec.Identity) (configkit.Identity, error) {
	return configkit.NewIdentity(
		in.GetName(),
		in.GetKey(),
	)
}

// marshalHandlerType marshals a configkit.HandlerType to its protocol buffers
// representation.
func marshalHandlerType(t configkit.HandlerType) (configspec.HandlerType, error) {
	if err := t.Validate(); err != nil {
		return configspec.HandlerType_UNKNOWN_HANDLER_TYPE, err
	}

	switch t {
	case configkit.AggregateHandlerType:
		return configspec.HandlerType_AGGREGATE, nil
	case configkit.ProcessHandlerType:
		return configspec.HandlerType_PROCESS, nil
	case configkit.IntegrationHandlerType:
		return configspec.HandlerType_INTEGRATION, nil
	default: // configkit.ProjectionHandlerType
		return configspec.HandlerType_PROJECTION, nil
	}
}

// unmarshalHandlerType unmarshals a configkit.HandlerType from its protocol
// buffers representation.
func unmarshalHandlerType(t configspec.HandlerType) (configkit.HandlerType, error) {
	switch t {
	case configspec.HandlerType_AGGREGATE:
		return configkit.AggregateHandlerType, nil
	case configspec.HandlerType_PROCESS:
		return configkit.ProcessHandlerType, nil
	case configspec.HandlerType_INTEGRATION:
		return configkit.IntegrationHandlerType, nil
	case configspec.HandlerType_PROJECTION:
		return configkit.ProjectionHandlerType, nil
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
