package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/discord"
	"github.com/kehiy/RoboPac/engine"
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
		if err != nil {
			kill(cmd, err)
		}

		// starting botEngine.
		botEngine, err := engine.NewBotEngine(config)
		if err != nil {
			kill(cmd, err)
		}

		botEngine.RegisterCommands()
		botEngine.Start()

		discordBot, err := discord.NewDiscordBot(botEngine, config.DiscordBotCfg.DiscordToken,
			config.DiscordBotCfg.DiscordGuildID)
		if err != nil {
			kill(cmd, err)
		}

		if err = discordBot.Start(); err != nil {
			kill(cmd, err)
		}

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigChan

		// gracefully shutdown the bot.
		discordBot.Stop()
		botEngine.Stop()
	}
}
