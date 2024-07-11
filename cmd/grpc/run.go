package main

import (
	"os"
	"os/signal"
	"syscall"

	pagucmd "github.com/pagu-project/Pagu/cmd"
	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/delivery/grpc"
	"github.com/pagu-project/Pagu/internal/engine"
	"github.com/pagu-project/Pagu/pkg/log"
	"github.com/spf13/cobra"
)

func runCommand(parentCmd *cobra.Command) {
	run := &cobra.Command{
		Use:   "run",
		Short: "Runs a mainnet instance of RoboPac",
	}

	parentCmd.AddCommand(run)

	run.Run = func(cmd *cobra.Command, _ []string) {
		// load configuration.
		config, err := config.Load("")
		pagucmd.ExitOnError(cmd, err)

		// Initialize global logger.
		log.InitGlobalLogger(config.Logger)

		// starting botEngine.
		botEngine, err := engine.NewBotEngine(config)
		pagucmd.ExitOnError(cmd, err)

		botEngine.RegisterAllCommands()
		botEngine.Start()

		grpcServer := grpc.NewServer(botEngine, config.GRPC)

		err = grpcServer.Start()
		pagucmd.ExitOnError(cmd, err)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigChan

		if err := grpcServer.Stop(); err != nil {
			pagucmd.ExitOnError(cmd, err)
		}

		botEngine.Stop()
	}
}
