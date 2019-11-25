package configkit

import "github.com/dogmatiq/configkit/message"

// MessageNames is describes how messages are used within a Dogma entity where
// message is identified by its name.
type MessageNames struct {
	// Roles is a map of message name to its role within the entity.
	Roles map[message.Name]message.Role

	// Produced is a list of message names produced by the entity.
	Produced []message.Name

	// Consumed is a list of message names consumed by the entity.
	Consumed []message.Name
}

// MessageTypes is describes how messages are used within a Dogma entity where
// message is identified by its type.
type MessageTypes struct {
	// Roles is a map of message type to its role within the entity.
	Roles map[message.Name]message.Role

	// Produced is a list of message types produced by the entity.
	Produced []message.Name

	// Consumed is a list of message types consumed by the entity.
	Consumed []message.Name
}
