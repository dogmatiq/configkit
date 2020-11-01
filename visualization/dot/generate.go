package dot

import (
	"github.com/dogmatiq/configkit"
	"github.com/emicklei/dot"
)

// Generate builds a message-flow graph for the given Dogma applications
// configurations.
//
// It panics if apps is empty.
func Generate(apps ...configkit.Application) (*dot.Graph, error) {
	if len(apps) == 0 {
		panic("at least one application must be provided")
	}

	g, err := (&generator{}).generate(apps)
	if err != nil {
		return nil, err
	}

	return g, nil
}
