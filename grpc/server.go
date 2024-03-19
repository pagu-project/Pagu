package grpc

import (
	"context"
	"net"

	"github.com/robopac-project/RoboPac/config"
	"github.com/robopac-project/RoboPac/engine"
	robopac "github.com/robopac-project/RoboPac/grpc/gen/go"
	"github.com/robopac-project/RoboPac/log"
	"google.golang.org/grpc"
)

type Server struct {
	ctx      context.Context
	cancel   context.CancelFunc
	listener net.Listener
	address  string
	engine   *engine.BotEngine
	grpc     *grpc.Server
	cfg      config.GRPCConfig
}

func NewServer(be *engine.BotEngine, cfg config.GRPCConfig) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		ctx:    ctx,
		cancel: cancel,
		engine: be,
		cfg:    cfg,
	}
}

func (s *Server) StartServer() {
	listener, err := net.Listen("tcp", "")
	if err != nil {
		log.Panic("can't start gRPC server", "err", err)
	}

	s.startListening(listener)
}

func (s *Server) startListening(listener net.Listener) {
	opts := make([]grpc.UnaryServerInterceptor, 0)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(opts...))

	robopacServer := newRoboPacServer(s)

	robopac.RegisterRoboPacServer(grpcServer, robopacServer)

	s.listener = listener
	s.address = listener.Addr().String()
	s.grpc = grpcServer

	log.Info("grpc started listening", "address", listener.Addr().String())
	go func() {
		if err := s.grpc.Serve(listener); err != nil {
			log.Error("error on grpc serve", "error", err)
		}
	}()
}

func (s *Server) StopServer() {
	s.cancel()

	if s.grpc != nil {
		s.grpc.Stop()
		s.listener.Close()
	}
}
