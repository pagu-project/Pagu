package telegram

import (
	"context"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/robopac-project/RoboPac/config"
	"github.com/robopac-project/RoboPac/engine"
	"github.com/robopac-project/RoboPac/engine/command"
	"github.com/robopac-project/RoboPac/log"
)

type TelegramBot struct {
	BotEngine       *engine.BotEngine
	ChatID          int64
	BotInstance     *gotgbot.Bot
	Config          *config.Config
	commandHandlers map[string]ext.Handler
	ctx             context.Context
	cancel          context.CancelFunc
}

func NewTelegramBot(botEngine *engine.BotEngine, token string, chatID int64, config *config.Config) (*TelegramBot, error) {
	bot, err := gotgbot.NewBot(token, nil)
	if err != nil {
		log.Error("Failed to create Telegram bot:", err)
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	commandHandlers := make(map[string]ext.Handler)

	return &TelegramBot{
		BotEngine:       botEngine,
		ChatID:          chatID,
		BotInstance:     bot,
		Config:          config,
		commandHandlers: commandHandlers,
		ctx:             ctx,
		cancel:          cancel,
	}, nil
}

func (bot *TelegramBot) Start() error {
	log.Info("Starting Telegram Bot...")

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Error("Error handling update:", err)
			bot.cancel()
			return ext.DispatcherActionNoop
		},
		Panic: func(b *gotgbot.Bot, ctx *ext.Context, r interface{}) {
			log.Error("Panic occurred:", r)
			bot.cancel()
		},
	})

	dispatcher.AddHandler(bot)

	updater := ext.NewUpdater(dispatcher, nil)

	go func() {
		log.Info("Starting polling for updates...")
		if err := updater.StartPolling(bot.BotInstance, nil); err != nil {
			log.Error("Failed to start polling:", err)
			bot.cancel()
		}
	}()

	return nil
}

func (bot *TelegramBot) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	// Check if the message is from a private chat (DM)
	if ctx.Update.Message.Chat.Type != "private" {
		return nil // Ignore non-private messages.
	}

	// Handle the update.
	if ctx.Update.Message != nil {
		// Extract the entire message, including commands.
		fullMessage := ctx.Update.Message.Text

		// Remove the slash prefix.
		fullMessage = strings.TrimPrefix(fullMessage, "/")

		// Split the message into an array.
		messageParts := strings.Split(fullMessage, " ")

		// Pass the array to the bot engine.
		res := bot.BotEngine.Run(command.AppIdTelegram, strconv.FormatInt(ctx.EffectiveSender.User.Id, 10), messageParts)

		// Check if the command execution resulted in an error.
		if res.Error != "" {
			log.Error("Failed to execute command:", res.Error)

			_, err := b.SendMessage(ctx.EffectiveChat.Id, "An error occurred while processing your request.", nil)
			if err != nil {
				log.Error("Failed to send error response:", err)
			}
			return nil
		}

		// Send the response back to the user.
		_, err := b.SendMessage(ctx.EffectiveChat.Id, res.Message, nil)
		if err != nil {
			log.Error("Failed to send response:", err)
		}

		return nil
	}

	return nil
}

func (bot *TelegramBot) RegisterCommandHandler(command string, handler CommandFunc) {
	bot.commandHandlers[command] = NewCommandHandler(handler)
}

func (bot *TelegramBot) GetName() string {
	return "TelegramBotHandler"
}

func (bot *TelegramBot) Stop() {
	log.Info("Shutting down Telegram Bot")
	bot.cancel()
}
