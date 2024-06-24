package grpc

import (
	"context"
	"strings"

	"github.com/pagu-project/Pagu/internal/entity"

	"github.com/pagu-project/Pagu/internal/delivery/grpc/gen/go"
)

type robopacServer struct {
	*Server
}

func newRoboPacServer(server *Server) *robopacServer {
	return &robopacServer{
		Server: server,
	}
}

func (rs *robopacServer) Run(_ context.Context, er *robopac.RunRequest) (*robopac.RunResponse, error) {
	beInput := []string{}

	tokens := strings.Split(er.Command, " ")
	beInput = append(beInput, tokens...)

	res := rs.engine.Run(entity.AppIdgRPC, er.Id, beInput)

	return &robopac.RunResponse{
		Response: res.Message,
	}, nil
}
