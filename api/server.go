package api

import (
	"context"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/api/internal/pb"
	"google.golang.org/grpc"
)

// RegisterServer registers a config server for the config applications.
func RegisterServer(
	s *grpc.Server,
	apps ...configkit.Application,
) {
	svr := &server{}

	for _, in := range apps {
		out, err := marshalApplication(in)
		if err != nil {
			panic(err)
		}

		svr.ListApplicationIdentitiesResponse.Identities = append(
			svr.ListApplicationIdentitiesResponse.Identities,
			out.Identity,
		)

		svr.ListApplicationsResponse.Applications = append(
			svr.ListApplicationsResponse.Applications,
			out,
		)
	}

	pb.RegisterConfigServer(s, svr)
}

var _ pb.ConfigServer = (*server)(nil)

type server struct {
	pb.ListApplicationIdentitiesResponse
	pb.ListApplicationsResponse
}

// ListApplicationIdentities returns the identity of all applications.
func (s *server) ListApplicationIdentities(
	ctx context.Context,
	req *pb.ListApplicationIdentitiesRequest,
) (*pb.ListApplicationIdentitiesResponse, error) {
	return &s.ListApplicationIdentitiesResponse, nil
}

// ListApplications returns the full configuration of all applications.
func (s *server) ListApplications(
	ctx context.Context,
	req *pb.ListApplicationsRequest,
) (*pb.ListApplicationsResponse, error) {
	return &s.ListApplicationsResponse, nil
}

// Watch blocks until the calling context is canceled.
func (s *server) Watch(_ *pb.WatchRequest, cs pb.Config_WatchServer) error {
	// Always send a single response, then wait till the client goes away or the
	// server is stopped.
	//
	// At least one response is necessary so that the client can determine
	// whether the config API is implemented by the server at all, after which
	// it will block waiting for the next response which will only return with
	// an error indicating that the server is gone.
	cs.Send(&pb.WatchResponse{})
	<-cs.Context().Done()
	return nil
}
