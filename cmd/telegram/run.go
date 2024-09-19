package main

import (
	"os"
	"os/signal"
	"syscall"

	pagucmd "github.com/pagu-project/Pagu/cmd"
	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/engine"
	"github.com/pagu-project/Pagu/internal/platforms/telegram"
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
		// Load configuration.
		configs, err := config.Load(configPath)
		pagucmd.ExitOnError(cmd, err)

		// Starting botEngine.
		botEngine, err := engine.NewBotEngine(configs)
		pagucmd.ExitOnError(cmd, err)

		log.InitGlobalLogger(configs.Logger)

		botEngine.RegisterAllCommands()
		botEngine.Start()

		//chatID := configs.Telegram.ChatID
		//groupLink := configs.Telegram.GroupLink
		telegramBot, err := telegram.NewTelegramBot(botEngine, configs.Telegram.BotToken, configs)
		pagucmd.ExitOnError(cmd, err)

		// register command handlers.
		//telegramBot.RegisterStartCommandHandler(groupLink)

		err = telegramBot.Start()
		pagucmd.ExitOnError(cmd, err)

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
