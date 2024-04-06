package telegram

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/robopac-project/RoboPac/log"
)

// CommandFunc is a function type that matches the signature of command handlers.
type CommandFunc func(b *gotgbot.Bot, ctx *ext.Context) error

// CommandHandler wraps a CommandFunc to satisfy the ext.Handler interface.
type CommandHandler struct {
	handler CommandFunc
}

func (bot *TelegramBot) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	return ctx.Update.Message.Chat.Type == "private"
}

func (bot *TelegramBot) Name() string {
	return "TelegramBotHandler"
}

// NewCommandHandler creates a new CommandHandler.
func NewCommandHandler(handler CommandFunc) *CommandHandler {
	return &CommandHandler{handler: handler}
}

func (ch *CommandHandler) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	// Check if the update is a message and if it's from a private chat.
	if ctx.Update.Message != nil && ctx.Update.Message.Chat.Type == "private" {
		// Manually check if the message is a command.
		if strings.HasPrefix(ctx.Update.Message.Text, "/") {
			// List of allowed commands
			allowedCommands := []string{"start", "help"}

			// Extract command from the update.
			command := strings.TrimPrefix(ctx.Update.Message.Text, "/")

			// if command has arguments.
			command = strings.Split(command, " ")[0]

			// Check if the command is in the list of allowed commands.
			for _, allowedCommand := range allowedCommands {
				if command == allowedCommand {
					return true // Process this update.
				}
			}
		}
	}

	log.Info("Ignoring unauthorized or unknown command:", ctx.Update.Message.Text)
	return false // Do not process this update.
}

// HandleUpdate calls the wrapped CommandFunc.
func (ch *CommandHandler) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return ch.handler(b, ctx)
}

// Name implements the ext.Handler interface.
func (ch *CommandHandler) Name() string {
	return "CommandHandler"
}
