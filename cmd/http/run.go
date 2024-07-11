package main

import (
	"os"
	"os/signal"
	"syscall"

	pagucmd "github.com/pagu-project/Pagu/cmd"
	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/delivery/http"
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
		configs, err := config.Load("")
		pagucmd.ExitOnError(cmd, err)

		// Initialize global logger.
		log.InitGlobalLogger(configs.Logger)

		// starting botEngine.
		botEngine, err := engine.NewBotEngine(configs)
		pagucmd.ExitOnError(cmd, err)

		botEngine.RegisterAllCommands()
		botEngine.Start()

		httpServer := http.NewHTTPServer(botEngine, configs.HTTP)

		err = httpServer.Start()
		pagucmd.ExitOnError(cmd, err)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigChan

		if err := httpServer.Stop(); err != nil {
			pagucmd.ExitOnError(cmd, err)
		}

		botEngine.Stop()
	}
}
