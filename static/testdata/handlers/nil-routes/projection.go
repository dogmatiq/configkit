package app

import (
	"context"

	"github.com/dogmatiq/dogma"
)

// ProjectionHandler is a test implementation of dogma.ProjectionMessageHandler.
type ProjectionHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (ProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("<projection>", "9208f704-641c-44e4-91dc-5274598b30bd")

	c.Routes(nil)
}

// HandleEvent updates the projection to reflect the occurrence of an event.
func (ProjectionHandler) HandleEvent(
	_ context.Context,
	_, _, _ []byte,
	_ dogma.ProjectionEventScope,
	_ dogma.Event,
) (ok bool, err error) {
	return false, nil
}

// ResourceVersion returns the version of the resource r.
func (ProjectionHandler) ResourceVersion(context.Context, []byte) ([]byte, error) {
	return nil, nil
}

// CloseResource informs the projection that the resource r will not be
// used in any future calls to HandleEvent().
func (ProjectionHandler) CloseResource(context.Context, []byte) error {
	return nil
}

// Compact reduces the size of the projection's data.
func (ProjectionHandler) Compact(context.Context, dogma.ProjectionCompactScope) error {
	return nil
}
