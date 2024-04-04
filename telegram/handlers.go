package telegram

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
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
