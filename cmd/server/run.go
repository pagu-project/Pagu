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
	parentCmd.AddCommand(run)

	run.Run = func(_ *cobra.Command, _ []string) {
		// initializing logger global instance.
		log.InitGlobalLogger()

		// load configuration.
		config, err := config.Load()
		if err != nil {
			log.Panic("error loading configuration", "err", err)
		}

		// starting botEngine.
		botEngine, err := engine.NewBotEngine(config)
		if err != nil {
			log.Panic("could not start discord bot", "err", err)
		}

		botEngine.Start()

		discordBot, err := discord.NewDiscordBot(botEngine, config.DiscordBotCfg.DiscordToken,
			config.DiscordBotCfg.DiscordGuildID)
		if err != nil {
			log.Panic("could not start discord bot", "err", err)
		}
		discordBot.Start()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigChan

		// gracefully shutdown the bot.
		discordBot.Stop()
		botEngine.Stop()
	}
}
