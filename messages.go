package configkit

import "github.com/dogmatiq/configkit/message"

// EntityMessageNames is describes how messages are used within a Dogma entity
// where message is identified by its name.
type EntityMessageNames struct {
	// Roles is a map of message name to its role within the entity.
	Roles map[message.Name]message.Role

	// Produced is a list of message names produced by the entity.
	Produced []message.Name

	// Consumed is a list of message names consumed by the entity.
	Consumed []message.Name
}

// EntityMessageTypes is describes how messages are used within a Dogma entity
// where message is identified by its type.
type EntityMessageTypes struct {
	// Roles is a map of message type to its role within the entity.
	Roles map[message.Name]message.Role

	// Produced is a list of message types produced by the entity.
	Produced []message.Name

	// Consumed is a list of message types consumed by the entity.
	Consumed []message.Name
}
