package markdown

import (
	"io"

	"github.com/dogmatiq/configkit"
)

// Generate builds a message-flow markdown tables for the given Dogma
// applications configurations and writes it to w.
//
// It panics if apps is empty.
func Generate(w io.Writer, apps ...configkit.Application) error {
	if len(apps) == 0 {
		panic("at least one application must be provided")
	}

	return (&generator{}).generate(w, apps)
}
