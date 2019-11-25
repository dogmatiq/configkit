package message

import (
	"bytes"
	"fmt"
)

// Role is an enumeration of the roles a message can perform within an
// application.
type Role byte

const (
	// CommandRole is the role for command messages.
	CommandRole Role = 'C'

	// EventRole is the role for event messages.
	EventRole Role = 'E'

	// TimeoutRole is the role for timeout messages.
	TimeoutRole Role = 'T'
)

// MessageRoles is a slice of the valid message roles.
var MessageRoles = []Role{
	CommandRole,
	EventRole,
	TimeoutRole,
}

const (
	commandRoleString = "command"
	eventRoleString   = "event"
	timeoutRoleString = "timeout"

	commandRoleShortString = "CMD"
	eventRoleShortString   = "EVT"
	timeoutRoleShortString = "TMO"

	commandRoleMarker = "?"
	eventRoleMarker   = "!"
	timeoutRoleMarker = "@"
)

var (
	commandRoleBytes = []byte(commandRoleString)
	eventRoleBytes   = []byte(eventRoleString)
	timeoutRoleBytes = []byte(timeoutRoleString)
)

// Validate return an error if r is not a valid message role.
func (r Role) Validate() error {
	switch r {
	case CommandRole,
		EventRole,
		TimeoutRole:
		return nil
	default:
		return fmt.Errorf("invalid message role: %#v", r)
	}
}

// MustValidate panics if r is not a valid message role.
func (r Role) MustValidate() {
	if err := r.Validate(); err != nil {
		panic(err)
	}
}

// Is returns true if r is one of the given roles.
func (r Role) Is(roles ...Role) bool {
	r.MustValidate()

	for _, x := range roles {
		x.MustValidate()

		if r == x {
			return true
		}
	}

	return false
}

// MustBe panics if r is not one of the given roles.
func (r Role) MustBe(roles ...Role) {
	if !r.Is(roles...) {
		panic("unexpected message role: " + r.String())
	}
}

// MustNotBe panics if r is one of the given roles.
func (r Role) MustNotBe(roles ...Role) {
	if r.Is(roles...) {
		panic("unexpected message role: " + r.String())
	}
}

// Marker returns a character that identifies the message role when displaying
// message types.
func (r Role) Marker() string {
	r.MustValidate()

	switch r {
	case CommandRole:
		return commandRoleMarker
	case EventRole:
		return eventRoleMarker
	default: // TimeoutRole
		return timeoutRoleMarker
	}
}

// ShortString returns a short (3-character) representation of the handler type.
func (r Role) ShortString() string {
	r.MustValidate()

	switch r {
	case CommandRole:
		return commandRoleShortString
	case EventRole:
		return eventRoleShortString
	default: // TimeoutRole
		return timeoutRoleShortString
	}
}

func (r Role) String() string {
	switch r {
	case CommandRole:
		return commandRoleString
	case EventRole:
		return eventRoleString
	case TimeoutRole:
		return timeoutRoleString
	default:
		return fmt.Sprintf("<invalid message role %#v>", r)
	}
}

// MarshalText returns a UTF-8 representation of the message role.
func (r Role) MarshalText() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	switch r {
	case CommandRole:
		return commandRoleBytes, nil
	case EventRole:
		return eventRoleBytes, nil
	default: // TimeoutRole
		return timeoutRoleBytes, nil
	}
}

// UnmarshalText unmarshals a role from its UTF-8 representation.
func (r *Role) UnmarshalText(text []byte) error {
	if bytes.Equal(text, commandRoleBytes) {
		*r = CommandRole
	} else if bytes.Equal(text, eventRoleBytes) {
		*r = EventRole
	} else if bytes.Equal(text, timeoutRoleBytes) {
		*r = TimeoutRole
	} else {
		return fmt.Errorf("invalid text representation of message role: %s", text)
	}

	return nil
}

// MarshalBinary returns a binary representation of the message role.
func (r Role) MarshalBinary() ([]byte, error) {
	return []byte{byte(r)}, r.Validate()
}

// UnmarshalBinary unmarshals a role from its binary representation.
func (r *Role) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("invalid binary representation of message role, expected exactly 1 byte")
	}

	*r = Role(data[0])
	return r.Validate()
}
