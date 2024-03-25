package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
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

	// Set up command handler
	bot.Bot.Handle("/start", func(c tele.Context) error {
		return c.Send("RoboPac has been started. Use the /help command to view all commands.")
	})

	return bot.registerCommands()
}

func (bot *TelegramBot) registerCommands() error {
	token := bot.Config.TelegramBotCfg.Token

	var commands []map[string]string
	beCmds := bot.BotEngine.Commands()
	for _, beCmd := range beCmds {
		if !beCmd.HasAppId(command.AppIdTelegram) {
			continue
		}
		commands = append(commands, map[string]string{
			"command":     "/" + beCmd.Name,
			"description": beCmd.Desc,
		})
		log.Info("Command registered for Telegram:", "name", beCmd.Name, "description", beCmd.Desc)
	}

	// Convert the commands to JSON so we can register them with telegram.
	jsonCommands, err := json.Marshal(commands)
	if err != nil {
		log.Info("Error marshalling commands:", err)
		return err
	}

	// Register the commands with Telegram.
	url := "https://api.telegram.org/bot" + token + "/setMyCommands"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonCommands))
	if err != nil {
		log.Info("Error registering commands:", err)
		return err
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		log.Info("Failed to register commands:", resp.Status)
		return errors.New("failed to register commands")
	}

	log.Info("Commands registered successfully")
	return nil
}

// func (bot *TelegramBot) respondErrMsg(errStr string, c tele.Context) {
// 	// Prepare the error message
// 	errorMsg := "Error: " + errStr

// 	// Send the error message back to the user
// 	chatID := tele.ChatID(c.Chat().ID)
// 	_, err := bot.Bot.Send(chatID, errorMsg)
// 	if err != nil {
// 		log.Error("Failed to send error message", "error", err)
// 	}
// }

func (bot *TelegramBot) respondResultMsg(res command.CommandResult, chatID int64) {
	// Prepare the result message
	var resultMsg string
	if res.Successful {
		resultMsg = "Successful: " + res.Message
	} else {
		resultMsg = "Failed: " + res.Message
	}

	// Send the result message back to the user
	_, err := bot.Bot.Send(tele.ChatID(chatID), resultMsg)
	if err != nil {
		log.Error("Failed to send result message", "error", err)
	}
}

func (bot *TelegramBot) commandHandler(b *tele.Bot, m *tele.Message) {
	if m.Chat.ID != bot.ChatID {
		b.Send(m.Chat, "Please only send messages in the pactus group chat.")
		return
	}

	beInput := []string{}

	commandParts := strings.Split(m.Text, " ")
	cmd := commandParts[0]
	args := commandParts[1:]

	beInput = append(beInput, cmd)
	beInput = append(beInput, args...)

	// Convert m.Sender.ID from int64 to string
	callerID := strconv.FormatInt(m.Sender.ID, 10)

	res := bot.BotEngine.Run(command.AppIdTelegram, callerID, beInput)
	bot.respondResultMsg(res, m.Chat.ID)
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
