package api

import (
	"context"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/enginekit/grpc/configgrpc"
)

// Server is an implementation of configspec.ConfigAPIServer.
type Server struct {
	response configgrpc.ListApplicationsResponse
}

var _ configgrpc.ConfigAPIServer = (*Server)(nil)

// NewServer returns an API server that serves the configuration of the given
// applications.
func NewServer(apps ...configkit.Application) *Server {
	s := &Server{}

	for _, in := range apps {
		out, err := configkit.ToProto(in)
		if err != nil {
			panic(err)
		}

		s.response.Applications = append(
			s.response.Applications,
			out,
		)
	}

	return s
}

// ListApplications returns the full configuration of all applications.
func (s *Server) ListApplications(
	context.Context,
	*configgrpc.ListApplicationsRequest,
) (*configgrpc.ListApplicationsResponse, error) {
	return &s.response, nil
}
