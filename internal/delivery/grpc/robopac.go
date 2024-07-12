package grpc

import (
	"context"
	"strings"

	robopac "github.com/pagu-project/Pagu/internal/delivery/grpc/gen/go"
	"github.com/pagu-project/Pagu/internal/entity"
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
	beInput := make(map[string]any)

	tokens := strings.Split(er.Command, " ")
	for _, t := range tokens {
		beInput[t] = t
	}

	res := rs.engine.Run(entity.AppIDgRPC, er.Id, beInput)

	return &robopac.RunResponse{
		Response: res.Message,
	}, nil
}
