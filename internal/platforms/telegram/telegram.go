package telegram

import (
	"context"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/pactus-project/pactus/util"
	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/engine"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/log"
	tele "gopkg.in/telebot.v3"
	"strconv"
	"strings"
	"time"
)

type Bot struct {
	ctx         context.Context
	cancel      context.CancelFunc
	engine      *engine.BotEngine
	botInstance *tele.Bot
	cfg         *config.Config
	target      string
}

func NewTelegramBot(botEngine *engine.BotEngine, token string, cfg *config.Config) (*Bot, error) {

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Error("Failed to create Telegram bot:", err)
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Bot{
		engine:      botEngine,
		botInstance: bot,
		cfg:         cfg,
		ctx:         ctx,
		cancel:      cancel,
	}, nil
}

func (bot *Bot) Start() error {
	bot.deleteAllCommands()
	if err := bot.registerCommands(); err != nil {
		return err
	}

	go func() {
		bot.botInstance.Start()
	}()

	log.Info("Starting Telegram Bot...")
	return nil
}

func (bot *Bot) deleteAllCommands() {
	//cmdsServer, _ := bot.Session.ApplicationCommands(bot.Session.State.User.ID, bot.cfg.GuildID)
	//cmdsGlobal, _ := bot.Session.ApplicationCommands(bot.Session.State.User.ID, "")
	//cmds := append(cmdsServer, cmdsGlobal...) //nolint
	//
	//for _, cmd := range cmds {
	//	err := bot.Session.ApplicationCommandDelete(cmd.ApplicationID, cmd.GuildID, cmd.ID)
	//	if err != nil {
	//		log.Error("unable to delete command", "error", err, "cmd", cmd.Name)
	//	} else {
	//		log.Info("telegram command unregistered", "name", cmd.Name)
	//	}
	//}
}

func (bot *Bot) registerCommands() error {
	rows := make([]tele.Row, 0)
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	commands := make([]string, 0)
	for i, beCmd := range bot.engine.Commands() {
		if !beCmd.HasAppID(entity.AppIDTelegram) {
			continue
		}

		switch bot.target {
		case config.BotNamePaguMainnet:
			if !util.IsFlagSet(beCmd.TargetFlag, command.TargetMaskMain) {
				continue
			}

		case config.BotNamePaguTestnet:
			if !util.IsFlagSet(beCmd.TargetFlag, command.TargetMaskTest) {
				continue
			}

		case config.BotNamePaguModerator:
			if !util.IsFlagSet(beCmd.TargetFlag, command.TargetMaskModerator) {
				continue
			}
		}

		log.Info("registering new command", "name", beCmd.Name, "desc", beCmd.Help, "index", i, "object", beCmd)

		btn := menu.Data(beCmd.Name, fmt.Sprintf("%s", beCmd.Name))
		commands = append(commands, beCmd.Name)
		rows = append(rows, menu.Row(btn))
		if beCmd.HasSubCommand() {
			subMenu := &tele.ReplyMarkup{ResizeKeyboard: true}
			subRows := make([]tele.Row, 0)
			for _, sCmd := range beCmd.SubCommands {
				switch bot.target {
				case config.BotNamePaguMainnet:
					if !util.IsFlagSet(sCmd.TargetFlag, command.TargetMaskMain) {
						continue
					}

				case config.BotNamePaguTestnet:
					if !util.IsFlagSet(sCmd.TargetFlag, command.TargetMaskTest) {
						continue
					}

				case config.BotNamePaguModerator:
					if !util.IsFlagSet(sCmd.TargetFlag, command.TargetMaskModerator) {
						continue
					}
				}

				log.Info("adding command sub-command", "command", beCmd.Name,
					"sub-command", sCmd.Name, "desc", sCmd.Help)

				commands = append(commands, beCmd.Name)
				subBtn := subMenu.Data(sCmd.Name, fmt.Sprintf("%s", sCmd.Name))
				//bot.botInstance.Handle(&subBtn, func(c tele.Context) error {
				//	return bot.commandHandler(c, c)
				//})
				subRows = append(subRows, subMenu.Row(subBtn))

				//for _, arg := range sCmd.Args {
				//	if arg.Desc == "" || arg.Name == "" {
				//		continue
				//	}
				//
				//	log.Info("adding sub command argument", "command", beCmd.Name,
				//		"sub-command", sCmd.Name, "argument", arg.Name, "desc", arg.Desc)
				//
				//	commands = append(commands, &tele.Command{
				//		Text:        sCmd.Name,
				//		Description: arg.Desc,
				//	})
				//}
			}
			bot.botInstance.Handle(&btn, func(c tele.Context) error {
				return c.Send(btn.Text, subMenu)
			})
			subMenu.Inline(subRows...)
		} else {
			//for _, arg := range beCmd.Args {
			//	if arg.Desc == "" || arg.Name == "" {
			//		continue
			//	}
			//
			//	log.Info("adding command argument", "command", beCmd.Name,
			//		"argument", arg.Name, "desc", arg.Desc)
			//
			//	bot.botInstance.Handle(arg.Name, func(c tele.Context) error {
			//		return c.Send("Hello!")
			//	})
			//
			//	commands = append(commands, &tele.Command{
			//		Text:        arg.Name,
			//		Description: arg.Desc,
			//	})
			//}
		}
	}

	menu.Inline(rows...)
	bot.botInstance.Handle("/start", func(c tele.Context) error {
		return c.Send("select an option", menu)
	})

	return bot.commandHandler(bot.botInstance.sess)
}

func (bot *Bot) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
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
		beInput := make(map[string]string)

		for _, t := range messageParts {
			beInput[t] = t
		}

		// Pass the array to the bot engine.
		callerID := strconv.FormatInt(ctx.EffectiveSender.User.Id, 10)
		res := bot.engine.Run(entity.AppIDTelegram, callerID, []string{}, beInput)

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

func (bot *Bot) commandHandler(c tele.Context) error {
	var commands []string
	args := make(map[string]string)

	cmd, err := bot.botInstance.Commands()
	if err != nil {
		return err
	}
	for _, c := range cmd {
		commands = append(commands, c.Text)
	}

	callerID := strconv.Itoa(int(c.Sender().ID))
	res := bot.engine.Run(entity.AppIDTelegram, callerID, commands, args)
	return c.Send(res.Message)
}

func (bot *Bot) Stop() {
	log.Info("Shutting down Telegram Bot")
	bot.cancel()
	bot.botInstance.Stop()
}
