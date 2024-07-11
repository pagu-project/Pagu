package telegram

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *TelegramBot) RegisterStartCommandHandler(tgLink string) {
	bot.RegisterCommandHandler("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		chat := ctx.Update.Message.GetChat()
		chatID := chat.Id

		welcomeMessage := fmt.Sprintf("Hello @%s, this is Pagu. "+
			"Kindly use /help command to learn how to interact with the bot. Join our telegram chat %s",
			chat.Username, tgLink)
		_, err := b.SendMessage(chatID, welcomeMessage, nil)
		if err != nil {
			return err
		}

		return nil
	})
}
