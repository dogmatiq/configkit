package markdown

import (
	"io"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
)

// generator generates message-flow markdown tables for one or more Dogma
// applications.
type generator struct {
	id        int
	roles     map[message.Name]message.Role
	producers map[message.Name]bool
	consumers map[message.Name]bool
}

// generate builds markdown tables of the given applications and writes it to w.
func (g *generator) generate(w io.Writer, apps []configkit.Application) error {
	return nil
}
