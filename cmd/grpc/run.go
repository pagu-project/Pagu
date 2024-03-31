package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/robopac-project/RoboPac/config"
	"github.com/robopac-project/RoboPac/engine"
	"github.com/robopac-project/RoboPac/grpc"
	"github.com/robopac-project/RoboPac/log"
	"github.com/spf13/cobra"
)

func RunCommand(parentCmd *cobra.Command) {
	run := &cobra.Command{
		Use:   "run",
		Short: "Runs a mainnet instance of RoboPac",
	}

	parentCmd.AddCommand(run)

	run.Run = func(cmd *cobra.Command, _ []string) {
		// load configuration.
		config, err := config.Load()
		ExitOnError(cmd, err)

		// Initialize global logger.
		log.InitGlobalLogger(config.LoggerConfig)

		// starting botEngine.
		botEngine, err := engine.NewBotEngine(config)
		ExitOnError(cmd, err)

		botEngine.RegisterAllCommands()
		botEngine.Start()

		grpcServer := grpc.NewServer(botEngine, config.GRPCConfig)

		err = grpcServer.Start()
		ExitOnError(cmd, err)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigChan

		// gracefully shutdown the bot.
		grpcServer.Stop()
		botEngine.Stop()
	}
}
