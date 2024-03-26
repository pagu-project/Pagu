package telegram

import (
	"strconv"
	"strings"
	"time"

	"github.com/pactus-project/pactus/util"
	"github.com/robopac-project/RoboPac/config"
	"github.com/robopac-project/RoboPac/engine"
	"github.com/robopac-project/RoboPac/engine/command"
	"github.com/robopac-project/RoboPac/log"
	"github.com/robopac-project/RoboPac/utils"

	tele "gopkg.in/telebot.v3"
)

type TelegramBot struct {
	BotEngine       *engine.BotEngine
	ChatID          int64
	Bot             *tele.Bot
	Config          *config.Config //config
	commandHandlers map[string]tele.HandlerFunc
}

func NewTelegramBot(botEngine *engine.BotEngine, token string, chatID int64, config *config.Config) (*TelegramBot, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Error("Failed to create Telegram bot:", err)
		return nil, err
	}

	commandHandlers := make(map[string]tele.HandlerFunc)

	return &TelegramBot{
		BotEngine:       botEngine,
		ChatID:          chatID,
		Bot:             bot,
		Config:          config,
		commandHandlers: commandHandlers,
	}, nil
}

func (bot *TelegramBot) Start() error {
	log.Info("Starting Telegram Bot...")

	// Middleware for restricting users to using the bot in only one chat group
	bot.Bot.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Chat().ID != bot.ChatID {
				log.Info("Unauthorized access attempt from chat ID:", c.Chat().ID)
				return nil // Ignore messages from unauthorized chats
			}
			return next(c) // Proceed to the next handler if it's the right chat group.
		}
	})

	// Middleware for error handling
	bot.Bot.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			err := next(c)
			if err != nil {
				log.Error("Unhandled error:", err)
				if err := c.Send("An error occurred while processing your request."); err != nil {
					log.Error("Failed to send error response:", err)
				}
			}
			return err
		}
	})

	// Set up a message handler for text messages in the group
	bot.Bot.Handle(tele.OnText, func(c tele.Context) error {
		// Extract the entire command, including arguments
		fullCommand := c.Message().Text

		// Remove the '/' prefix if present
		fullCommand = strings.TrimPrefix(fullCommand, "/")
		log.Info("Received command from user:", "User ID", c.Sender().ID, "Command", fullCommand)

		// Split the command into an array
		commandParts := strings.Split(fullCommand, " ")
		log.Info("Command parts:", commandParts)

		// Pass the array to the bot engine
		res := bot.BotEngine.Run(command.AppIdTelegram, strconv.FormatInt(c.Sender().ID, 10), commandParts)
		if res.Error != "" {
			log.Error("Failed to execute command:", res.Error)
			if err := c.Send("An error occurred while processing your request."); err != nil {
				log.Error("Failed to send error response:", err)
			}
			return nil
		}

		// Send the response back to the user
		if err := c.Send(res.Message); err != nil {
			log.Error("Failed to send response:", err)
		}

		// Print the result of the Run method
		log.Info("Result of executing command:", "Message", res.Message)

		return nil
	})

	// Attempt to send a startup confirmation message
	msg, err := bot.Bot.Send(tele.ChatID(bot.ChatID), "Telegram Bot started successfully!")
	if err != nil {
		log.Error("Failed to send startup confirmation message:", err)
		return err
	}
	log.Info("Telegram Bot started successfully. Message ID:", msg.ID)

	return nil
}

func (bot *TelegramBot) UpdateStatusInfo(cfg *config.Config) {
	log.Info("info status started")
	for {
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

		// Convert the ChatId string to int64 and then to telebot.ChatID
		chatIdInt64, err := strconv.ParseInt(cfg.TelegramBotCfg.ChatId, 10, 64)
		if err != nil {
			log.Error("Failed to parse ChatId", "error", err)
			continue
		}
		chatID := tele.ChatID(chatIdInt64)

		for _, statusMessage := range statusMessages {
			_, err := bot.Bot.Send(chatID, statusMessage)
			if err != nil {
				log.Error("Failed to send status message", "error", err)
			}
			time.Sleep(time.Minute * 2) // Wait for 2 minutes before sending the next message
		}
	}
}

func (bot *TelegramBot) Stop() {
	log.Info("Shutting down Telegram Bot")
}
