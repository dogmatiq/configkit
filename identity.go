package configkit

import (
	"bytes"
	"errors"
	"fmt"
	"unicode"

	"github.com/dogmatiq/configkit/internal/validation"
)

// Identity is the application-defined identity of a Dogma entity.
type Identity struct {
	// Name is the name component of the identity.
	//
	// For handlers, it is unique within an application at any given version,
	// but may be changed over time.
	//
	// It is allowed, but not recommended to use the same name for an
	// application as one of its constituent handlers.
	Name string

	// Key is the key component of the identity.
	//
	// It is not only unique within an application, but forever immutable. It is
	// not permitted for an application and one of its constituent handlers to
	// share the same key.
	Key string
}

// NewIdentity returns a new identity.
//
// It returns a non-nil error if either of the name or key components is invalid.
func NewIdentity(n, k string) (Identity, error) {
	i := Identity{n, k}
	return i, i.Validate()
}

// MustNewIdentity returns a new identity.
//
// It panics if either of the name or key components is invalid.
func MustNewIdentity(n, k string) Identity {
	i, err := NewIdentity(n, k)
	if err != nil {
		panic(err)
	}

	return i
}

// IsZero returns true if the identity is the zero-value.
func (i Identity) IsZero() bool {
	return i.Name == "" && i.Key == ""
}

// ConflictsWith returns true if i has the same name or key as ident.
func (i Identity) ConflictsWith(ident Identity) bool {
	if i.Name == ident.Name {
		return true
	}

	if i.Key == ident.Key {
		return true
	}

	return false
}

// Validate returns an error if i is not a valid identity.
func (i Identity) Validate() error {
	if !isValidIdentityComponent(i.Name) {
		return validation.Errorf(
			"invalid name %#v, names must be non-empty, printable UTF-8 strings with no whitespace",
			i.Name,
		)
	}

	if !isValidIdentityComponent(i.Key) {
		return validation.Errorf(
			"invalid key %#v, keys must be non-empty, printable UTF-8 strings with no whitespace",
			i.Key,
		)
	}

	return nil
}

func (i Identity) String() string {
	return fmt.Sprintf("%s/%s", i.Name, i.Key)
}

// MarshalText returns a UTF-8 representation of the identity.
func (i Identity) MarshalText() ([]byte, error) {
	if err := i.Validate(); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString(i.Name)
	buf.WriteRune(' ')
	buf.WriteString(i.Key)

	return buf.Bytes(), nil
}

// UnmarshalText unmarshals an identity from its UTF-8 representation.
func (i *Identity) UnmarshalText(text []byte) error {
	n := bytes.IndexRune(text, ' ')
	if n == -1 {
		return errors.New("could not decode identity, no name/key separator found")
	}

	i.Name = string(text[:n])
	i.Key = string(text[n+1:])

	return i.Validate()
}

// MarshalBinary returns a binary representation of the identity.
func (i Identity) MarshalBinary() ([]byte, error) {
	return i.MarshalText()
}

// UnmarshalBinary unmarshals an identity from its binary representation.
func (i *Identity) UnmarshalBinary(data []byte) error {
	return i.UnmarshalText(data)
}

// isValidIdentityComponent returns true if n is a valid application or handler
// name or key.
//
// A valid name/key is a non-empty string consisting of Unicode printable
// characters, except whitespace.
func isValidIdentityComponent(n string) bool {
	if len(n) == 0 {
		return false
	}

	for _, r := range n {
		if unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return false
		}
	}

	return true
}
