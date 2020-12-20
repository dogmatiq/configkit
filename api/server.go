package api

import (
	"context"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/interopspec/configspec"
)

// NewServer returns a configspec.ConfigAPIServer for the given applications.
func NewServer(apps ...configkit.Application) configspec.ConfigAPIServer {
	s := &server{}

	for _, in := range apps {
		out, err := marshalApplication(in)
		if err != nil {
			panic(err)
		}

		s.ListApplicationsResponse.Applications = append(
			s.ListApplicationsResponse.Applications,
			out,
		)
	}

	return s
}

type server struct {
	configspec.ListApplicationsResponse
}

// ListApplications returns the full configuration of all applications.
func (s *server) ListApplications(
	context.Context,
	*configspec.ListApplicationsRequest,
) (*configspec.ListApplicationsResponse, error) {
	return &s.ListApplicationsResponse, nil
}
