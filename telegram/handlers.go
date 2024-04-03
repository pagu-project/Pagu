package telegram

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/robopac-project/RoboPac/engine/command"
	"github.com/robopac-project/RoboPac/log"
)

func (bot *TelegramBot) RegisterStartCommandHandler() {
	bot.RegisterCommandHandler("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		// log.Info("Start command received")

		chat := ctx.Update.Message.GetChat()
		chatID := chat.Id

		tg_link := "https://t.me/pactuschat"
		welcomeMessage := fmt.Sprintf("Hello @%s, this is RoboPac. Kindly use /help command to learn how to interact with the bot. Join our telegram chat %s", chat.Username, tg_link)
		_, err := b.SendMessage(chatID, welcomeMessage, nil)
		if err != nil {
			// Log any error encountered during message sending
			// log.Info("Failed to send start message: %v", err)
			return err
		}

		// log.Info("Start message sent successfully")
		return nil
	})
}

func (bot *TelegramBot) RegisterBotEngineCommandHandler() {
	// Iterate over each command provided by the bot engine
	for _, cmd := range bot.BotEngine.Commands() {
		// Register a handler for each command
		bot.RegisterCommandHandler(cmd.Name, func(b *gotgbot.Bot, ctx *ext.Context) error {
			// Extract the entire command, including arguments
			fullCommand := ctx.Message.Text

			// Remove the '/' prefix if present
			fullCommand = strings.TrimPrefix(fullCommand, "/")
			log.Info(fmt.Sprintf("Received command from UserID %d: '%s'", ctx.EffectiveSender.User.Id, fullCommand))

			// Split the command into an array
			commandParts := strings.Split(fullCommand, " ")
			log.Info(fmt.Sprintf("Processing command parts: %v", commandParts))

			// Pass the array to the bot engine
			res := bot.BotEngine.Run(command.AppIdTelegram, strconv.FormatInt(ctx.EffectiveSender.User.Id, 10), commandParts)
			var err error // Declare err variable
			if res.Error != "" {
				log.Error("Failed to execute command:", res.Error)
				_, err = b.SendMessage(ctx.EffectiveChat.Id, "An error occurred while processing your request.", nil)
				if err != nil {
					log.Error("Failed to send error response:", err)
				}
				return nil
			}

			// Send the response back to the user
			_, err = b.SendMessage(ctx.EffectiveChat.Id, res.Message, nil)
			if err != nil {
				log.Error("Failed to send response:", err)
			}

			return nil
		})
	}
}
