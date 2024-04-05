package main

import (
	"os"
	"os/signal"
	"syscall"

	rpCmd "github.com/robopac-project/RoboPac/cmd"
	"github.com/robopac-project/RoboPac/config"
	"github.com/robopac-project/RoboPac/engine"
	"github.com/robopac-project/RoboPac/log"
	"github.com/robopac-project/RoboPac/telegram"
	"github.com/spf13/cobra"
)

func RunCommand(parentCmd *cobra.Command) {
	run := &cobra.Command{
		Use:   "run",
		Short: "Runs a mainnet instance of RoboPac",
	}

	parentCmd.AddCommand(run)

	run.Run = func(cmd *cobra.Command, _ []string) {
		// Load configuration.
		config, err := config.Load()
		rpCmd.ExitOnError(cmd, err)

		// Starting botEngine.
		botEngine, err := engine.NewBotEngine(config)
		rpCmd.ExitOnError(cmd, err)

		log.InitGlobalLogger(config.Logger)

		botEngine.RegisterAllCommands()
		botEngine.Start()

		chatID := config.Telegram.ChatID
		tgLink := config.Telegram.TgLink

		telegramBot, err := telegram.NewTelegramBot(botEngine, config.Telegram.BotToken, chatID, config)
		rpCmd.ExitOnError(cmd, err)

		// register command handlers.
		telegramBot.RegisterStartCommandHandler(tgLink)

		err = telegramBot.Start()
		rpCmd.ExitOnError(cmd, err)

		// Set up signal handling.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			<-c
			// When a signal is received, stop the bot and perform any other necessary cleanup.
			telegramBot.Stop()
			botEngine.Stop()
			os.Exit(1)
		}()

		// Block the main goroutine until a signal is received.
		select {}
	}
}
