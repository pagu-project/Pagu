package main

import (
	"os"
	"os/signal"
	"syscall"

	rpCmd "github.com/robopac-project/RoboPac/cmd"
	"github.com/robopac-project/RoboPac/config"
	"github.com/robopac-project/RoboPac/engine"
	"github.com/robopac-project/RoboPac/grpc"
	"github.com/robopac-project/RoboPac/log"
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
		config, err := config.Load()
		rpCmd.ExitOnError(cmd, err)

		// Initialize global logger.
		log.InitGlobalLogger(config.Logger)

		// starting botEngine.
		botEngine, err := engine.NewBotEngine(config)
		rpCmd.ExitOnError(cmd, err)

		botEngine.RegisterAllCommands()
		botEngine.Start()

		grpcServer := grpc.NewServer(botEngine, config.GRPC)

		err = grpcServer.Start()
		rpCmd.ExitOnError(cmd, err)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigChan

		if err := grpcServer.Stop(); err != nil {
			rpCmd.ExitOnError(cmd, err)
		}

		botEngine.Stop()
	}
}
