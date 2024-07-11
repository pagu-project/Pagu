package main

import (
	"os"
	"os/signal"
	"syscall"

	pagucmd "github.com/pagu-project/Pagu/cmd"
	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/engine"
	"github.com/pagu-project/Pagu/internal/platforms/discord"
	"github.com/pagu-project/Pagu/pkg/log"
	"github.com/spf13/cobra"
)

func runCommand(parentCmd *cobra.Command) {
	run := &cobra.Command{
		Use:   "run",
		Short: "Runs an instance of Pagu",
	}

	parentCmd.AddCommand(run)

	run.Run = func(cmd *cobra.Command, _ []string) {
		// load configuration.
		configs, err := config.Load(configPath)
		pagucmd.ExitOnError(cmd, err)

		// Initialize global logger.
		log.InitGlobalLogger(configs.Logger)

		// starting botEngine.
		botEngine, err := engine.NewBotEngine(configs)
		pagucmd.ExitOnError(cmd, err)

		botEngine.RegisterAllCommands()
		botEngine.Start()

		discordBot, err := discord.NewDiscordBot(botEngine, configs.DiscordBot, configs.BotName)
		pagucmd.ExitOnError(cmd, err)

		err = discordBot.Start()
		pagucmd.ExitOnError(cmd, err)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigChan

		if err := discordBot.Stop(); err != nil {
			pagucmd.ExitOnError(cmd, err)
		}

		botEngine.Stop()
	}
}
