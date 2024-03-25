package telegram

import (
	"strconv"
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
		return nil, err
	}

	commandHandlers := make(map[string]tele.HandlerFunc)

	return &TelegramBot{
		BotEngine:       botEngine,
		ChatID:          chatID,
		Bot:             bot,
		Config:          config, //config
		commandHandlers: commandHandlers,
	}, nil
}

func (bot *TelegramBot) Start() error {
	log.Info("starting Telegram Bot...")

	// Set up command handler for /start
	bot.Bot.Handle("/start", func(c tele.Context) error {
		log.Info("Received /start command from user:", "User ID", c.Sender().ID)
		if err := c.Send("RoboPac has been started. Use the /help command to view all commands."); err != nil {
			log.Error("Failed to send /start response:", err)
		}
		return nil
	})

	// Set up command handler for /help
	bot.Bot.Handle("/help", func(c tele.Context) error {
		log.Info("Received /help command from user:", "User ID", c.Sender().ID)
		beCmds := bot.BotEngine.Commands()
		var helpText string
		for _, beCmd := range beCmds {
			if beCmd.HasAppId(command.AppIdTelegram) {
				helpText += "/" + beCmd.Name + ": " + beCmd.Desc + "\n"
			}
		}
		if err := c.Send(helpText); err != nil {
			log.Error("Failed to send /help response:", err)
		}
		return nil
	})

	beCmds := bot.BotEngine.Commands()
	for _, beCmd := range beCmds {
		log.Info("Fetched command from bot engine:", "Name", beCmd.Name, "Description", beCmd.Desc)
		if beCmd.HasAppId(command.AppIdTelegram) {
			cmd := beCmd
			bot.Bot.Handle("/"+cmd.Name, func(c tele.Context) error {
				log.Info("Received command "+cmd.Name+" from user:", "User ID", c.Sender().ID)
				if c.Chat().ID != bot.ChatID {
					if err := c.Send("Commands are only allowed in the Pactus group chat!"); err != nil {
						log.Error("Failed to send message:", err)
					}
					return nil
				}

				// Extract the entire command, including arguments
				fullCommand := c.Message().Text

				// Pass the full command to the bot engine
				res := bot.BotEngine.Run(command.AppIdTelegram, strconv.FormatInt(c.Sender().ID, 10), []string{fullCommand})
				if err := c.Send(res.Message); err != nil {
					log.Error("Failed to send response:", err)
				}

				// Print the result of the Run method
				log.Info("Result of executing command "+cmd.Name+":", "Message", res.Message)

				return nil
			})
		}
	}

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

	// Check if the bot started successfully, if successful then it sends a message to the chat group/channel
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
