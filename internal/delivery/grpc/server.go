package grpc

import (
	"context"
	"net"

	pagu "github.com/pagu-project/Pagu/internal/delivery/grpc/gen/go"
	"github.com/pagu-project/Pagu/internal/engine"
	"github.com/pagu-project/Pagu/pkg/log"

	"github.com/pagu-project/Pagu/config"
	"google.golang.org/grpc"
)

type Server struct {
	ctx      context.Context
	cancel   context.CancelFunc
	listener net.Listener
	address  string
	engine   *engine.BotEngine
	grpc     *grpc.Server
	cfg      config.GRPC
}

func NewServer(be *engine.BotEngine, cfg config.GRPC) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		ctx:    ctx,
		cancel: cancel,
		engine: be,
		cfg:    cfg,
	}
}

func (s *Server) Start() error {
	log.Info("Starting gRPC Server")
	listener, err := net.Listen("tcp", s.cfg.Listen)
	if err != nil {
		return err
	}

	s.startListening(listener)
	return nil
}

func (s *Server) startListening(listener net.Listener) {
	opts := make([]grpc.UnaryServerInterceptor, 0)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(opts...))

	robopacServer := newRoboPacServer(s)

	pagu.RegisterRoboPacServer(grpcServer, robopacServer)

	s.listener = listener
	s.address = listener.Addr().String()
	s.grpc = grpcServer

	log.Info("gRPC Server Started Listening", "address", listener.Addr().String())
	go func() {
		if err := s.grpc.Serve(listener); err != nil {
			log.Error("error on grpc serve", "error", err)
		}
	}()
}

func (s *Server) Stop() error {
	log.Info("Stopping gRPC Server", "addr", s.address)

	s.cancel()

	if s.grpc != nil {
		s.grpc.Stop()
	}

	return s.listener.Close()
}
