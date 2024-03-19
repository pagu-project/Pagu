package grpc

import (
	"context"

	robopac "github.com/robopac-project/RoboPac/grpc/gen/go"
)

type robopacServer struct {
	*Server
}

func newRoboPacServer(server *Server) *robopacServer {
	return &robopacServer{
		Server: server,
	}
}

func (rs *robopacServer) Execute(context.Context, *robopac.ExecuteRequest) (*robopac.ExecuteResponse, error) {
	return nil, nil
}
