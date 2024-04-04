package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

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
		ExitOnError(cmd, err)

		// Starting botEngine.
		botEngine, err := engine.NewBotEngine(config)
		ExitOnError(cmd, err)

		log.InitGlobalLogger(config.Logger)

		botEngine.RegisterAllCommands()
		botEngine.Start()

		chatID, err := strconv.ParseInt(config.Telegram.ChatID, 10, 64)
		if err != nil {
			log.Error("Failed to parse ChatId:", err)
			return
		}

		telegramBot, err := telegram.NewTelegramBot(botEngine, config.Telegram.BotToken, chatID, config)
		ExitOnError(cmd, err)

		// register command handlers.
		telegramBot.RegisterStartCommandHandler()

		err = telegramBot.Start()
		ExitOnError(cmd, err)

		log.Info("Telegram Bot started successfully")

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
