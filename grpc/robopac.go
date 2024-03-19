package grpc

import (
	"context"
	"strings"

	"github.com/robopac-project/RoboPac/engine/command"
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

func (rs *robopacServer) Execute(_ context.Context, er *robopac.ExecuteRequest) (*robopac.ExecuteResponse, error) {
	beInput := []string{}

	tokens := strings.Split(er.Command, " ")
	beInput = append(beInput, tokens...)

	res := rs.engine.Run(command.AppIdgRPC, "", beInput)

	return &robopac.ExecuteResponse{
		Response: res.Message,
	}, nil
}
