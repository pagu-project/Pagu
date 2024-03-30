package telegram

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robopac-project/RoboPac/engine/command"
	"github.com/robopac-project/RoboPac/log"
	tele "gopkg.in/telebot.v3"
)

func (bot *TelegramBot) commandHandler(c tele.Context) error {
	// Ensure commands are only processed in the specified group chat
	if c.Chat().ID != bot.ChatID {
		log.Info("Unauthorized access attempt from chat ID:", c.Chat().ID)
		return nil // Ignore commands from unauthorized chats
	}

	// Extract the command text
	fullCommand := c.Message().Text
	log.Info(fmt.Sprintf("Received command from UserID %d in ChatID %d: '%s'", c.Sender().ID, c.Chat().ID, fullCommand))

	// Process the command
	commandParts := strings.Split(strings.TrimPrefix(fullCommand, "/"), " ")
	if len(commandParts) == 0 {
		return c.Send("Invalid command format.")
	}

	// Send command to the bot engine and get the response
	res := bot.BotEngine.Run(command.AppIdTelegram, strconv.FormatInt(c.Sender().ID, 10), commandParts)
	if res.Error != "" {
		log.Error("Failed to execute command:", "error", res.Error)
		return c.Send("An error occurred while processing your request.")
	}

	// Send the response back to the user
	if err := c.Send(res.Message); err != nil {
		log.Error("Failed to send response:", "error", err)
		return err
	}

	log.Info("Successfully processed command:", "response", res.Message)
	return nil
}

func StartCommandHandler(c tele.Context) error {
	tg_link := "https://t.me/pactuschat"
	chat := c.Chat()
	welcomeMessage := fmt.Sprintf("Hello @%s, you're interacting with robopac, Pactus's helper bot. Please join the Pactus group chat at %s and use commands there.", chat.Username, tg_link)
	return c.Send(welcomeMessage)
}
