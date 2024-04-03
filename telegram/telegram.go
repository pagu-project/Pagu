package telegram

import (
	"context"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/pactus-project/pactus/util"
	"github.com/robopac-project/RoboPac/config"
	"github.com/robopac-project/RoboPac/engine"
	"github.com/robopac-project/RoboPac/log"
	"github.com/robopac-project/RoboPac/utils"
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

func (bot *TelegramBot) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	// Check if the message is from a private chat (DM)
	if ctx.Update.Message.Chat.Type != "private" {
		return nil // Ignore messages from non-private chats
	}

	// Handle the update
	if ctx.Update.Message != nil {
		// Check if the message is a command
		if strings.HasPrefix(ctx.Update.Message.Text, "/") {
			// Extract the command from the message text
			command := strings.TrimPrefix(ctx.Update.Message.Text, "/")

			// Iterate through registered command handlers
			for registeredCommand, handler := range bot.commandHandlers {
				// Check if the message is a command and matches the registered command.
				if command == registeredCommand {
					// Execute the command handler
					return handler.HandleUpdate(b, ctx)
				}
			}
		}
	}

	// Handle unknown commands or non-command messages here
	return nil
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

func (bot *TelegramBot) RegisterCommandHandler(command string, handler CommandFunc) {
	bot.commandHandlers[command] = NewCommandHandler(handler)
}

func (bot *TelegramBot) UpdateStatusInfo() {
	log.Info("info status started")
	for {
		select {
		case <-bot.ctx.Done():
			log.Info("Stopping status update due to context cancellation.")
			return
		default:
			ns, err := bot.BotEngine.NetworkStatus()
			if err != nil {
				log.Error("Failed to fetch network status", "error", err)
				continue
			}

			statusMessages := []string{
				"Validators count: " + utils.FormatNumber(int64(ns.ValidatorsCount)),
				"Total accounts: " + utils.FormatNumber(int64(ns.TotalAccounts)),
				"Current block height: " + utils.FormatNumber(int64(ns.CurrentBlockHeight)),
				"Circulating supply: " + utils.FormatNumber(int64(util.ChangeToCoin(ns.CirculatingSupply))) + " PAC",
				"Total network power: " + utils.FormatNumber(int64(util.ChangeToCoin(ns.TotalNetworkPower))) + " PAC",
			}

			for _, statusMessage := range statusMessages {
				_, err := bot.BotInstance.SendMessage(bot.ChatID, statusMessage, nil)
				if err != nil {
					log.Error("Failed to send status message", "error", err)
				}
				time.Sleep(time.Minute * 2) // Wait for 2 minutes before sending the next message
			}
		}
	}
}

func (bot *TelegramBot) GetName() string {
	return "TelegramBotHandler"
}

func (bot *TelegramBot) Stop() {
	log.Info("Shutting down Telegram Bot")
	bot.cancel()
}
