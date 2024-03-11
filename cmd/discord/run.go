package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/discord"
	"github.com/kehiy/RoboPac/engine"
	"github.com/kehiy/RoboPac/log"
	"github.com/spf13/cobra"
)

func RunCommand(parentCmd *cobra.Command) {
	run := &cobra.Command{
		Use:   "run",
		Short: "Runs a mainnet instance of RoboPac",
	}
	log.InitGlobalLogger()

	parentCmd.AddCommand(run)

	run.Run = func(cmd *cobra.Command, _ []string) {
		// load configuration.
		config, err := config.Load()
		ExitOnError(cmd, err)

		// starting botEngine.
		botEngine, err := engine.NewBotEngine(config)
		ExitOnError(cmd, err)

		botEngine.RegisterAllCommands()
		botEngine.Start()

		discordBot, err := discord.NewDiscordBot(botEngine, config.DiscordBotCfg.DiscordToken,
			config.DiscordBotCfg.DiscordGuildID)
		ExitOnError(cmd, err)

		err = discordBot.Start()
		ExitOnError(cmd, err)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigChan

		// gracefully shutdown the bot.
		discordBot.Stop()
		botEngine.Stop()
	}
}
