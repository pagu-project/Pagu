package telegram

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pactus-project/pactus/util"
	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/engine"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/log"
	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	ctx         context.Context
	cancel      context.CancelFunc
	engine      *engine.BotEngine
	botInstance *tele.Bot
	cfg         *config.Config
	target      string
}

type BotContext struct {
	Commands []string
}

var (
	argsContext = make(map[int64]*BotContext)
	argsValue   = make(map[int64]map[string]string)
)

func NewTelegramBot(botEngine *engine.BotEngine, token string, cfg *config.Config) (*Bot, error) {
	pref := tele.Settings{
		Token: token,
	}

	tgb, err := tele.NewBot(pref)
	if err != nil {
		log.Error("Failed to create Telegram bot:", err)
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Bot{
		engine:      botEngine,
		botInstance: tgb,
		cfg:         cfg,
		ctx:         ctx,
		cancel:      cancel,
		target:      cfg.BotName,
	}, nil
}

func (bot *Bot) Start() error {
	bot.deleteAllCommands()
	if err := bot.registerCommands(); err != nil {
		return err
	}

	go bot.botInstance.Start()
	log.Info("Starting Telegram Bot...")

	return nil
}

func (bot *Bot) Stop() {
	log.Info("Shutting down Telegram Bot")
	bot.cancel()
	bot.botInstance.Stop()
}

func (bot *Bot) deleteAllCommands() {
	for _, cmd := range bot.engine.Commands() {
		err := bot.botInstance.DeleteCommands(tele.Command{Text: cmd.Name})
		if err != nil {
			log.Error("unable to delete command", "error", err, "cmd", cmd.Name)
		}
	}
}

func (bot *Bot) registerCommands() error {
	rows := make([]tele.Row, 0)
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	commands := make([]tele.Command, 0)

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
		}

		log.Info("registering new command", "name", beCmd.Name, "desc", beCmd.Help, "index", i, "object", beCmd)

		btn := menu.Data(beCmd.Name, beCmd.Name)
		commands = append(commands, tele.Command{Text: beCmd.Name, Description: beCmd.Help})
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

				subBtn := subMenu.Data(sCmd.Name, sCmd.Name)

				bot.botInstance.Handle(&subBtn, func(c tele.Context) error {
					if len(sCmd.Args) > 0 {
						return bot.handleArgCommand(c, []string{beCmd.Name, sCmd.Name}, sCmd.Args)
					}

					return bot.handleCommand(c, []string{beCmd.Name, sCmd.Name})
				})
				subRows = append(subRows, subMenu.Row(subBtn))
			}

			subMenu.Inline(subRows...)
			bot.botInstance.Handle(&btn, func(c tele.Context) error {
				_ = bot.botInstance.Delete(c.Message())
				return c.Send(beCmd.Name, subMenu)
			})

			bot.botInstance.Handle(fmt.Sprintf("/%s", beCmd.Name), func(c tele.Context) error {
				_ = bot.botInstance.Delete(c.Message())
				return c.Send(beCmd.Name, subMenu)
			})
		} else {
			bot.botInstance.Handle(&btn, func(c tele.Context) error {
				if len(beCmd.Args) > 0 {
					return bot.handleArgCommand(c, []string{beCmd.Name}, beCmd.Args)
				}

				return bot.handleCommand(c, []string{beCmd.Name})
			})
		}
	}

	// initiate menu button
	_ = bot.botInstance.SetCommands(commands)
	menu.Inline(rows...)
	bot.botInstance.Handle("/start", func(c tele.Context) error {
		_ = bot.botInstance.Delete(c.Message())
		return c.Send("Pagu Main Menu", menu)
	})

	bot.botInstance.Handle(tele.OnText, func(c tele.Context) error {
		if argsContext[c.Message().Chat.ID] == nil {
			return c.Send("Pagu Main Menu", menu)
		}

		if argsValue[c.Message().Chat.ID] == nil {
			argsValue[c.Message().Chat.ID] = make(map[string]string)
		}

		return bot.parsTestMessage(c)
	})

	return nil
}

func (bot *Bot) parsTestMessage(c tele.Context) error {
	chatID := c.Message().Chat.ID
	cmd := findCommand(bot.engine.Commands(), argsContext[chatID].Commands[len(argsContext[chatID].Commands)-1])
	if cmd == nil {
		return c.Send("Invalid command")
	}

	currentArgsIndex := len(argsValue[chatID])
	argsValue[chatID][cmd.Args[currentArgsIndex].Name] = c.Message().Text

	if len(argsValue[chatID]) == len(cmd.Args) {
		return bot.handleCommand(c, argsContext[chatID].Commands)
	}

	return c.Send(fmt.Sprintf("Please Enter %s", cmd.Args[currentArgsIndex+1].Name))
}

func (bot *Bot) handleArgCommand(c tele.Context, commands []string, args []command.Args) error {
	msgCtx := &BotContext{Commands: commands}
	argsContext[c.Chat().ID] = msgCtx
	return c.Send(fmt.Sprintf("Please Enter %s", args[0].Name))
}

func (bot *Bot) handleCommand(c tele.Context, commands []string) error {
	callerID := strconv.Itoa(int(c.Sender().ID))
	res := bot.engine.Run(entity.AppIDTelegram, callerID, commands, argsValue[c.Message().Chat.ID])
	_ = bot.botInstance.Delete(c.Message())

	chatID := c.Message().Chat.ID
	argsContext[chatID] = nil
	argsValue[chatID] = nil
	return c.Send(res.Message)
}

func findCommand(commands []*command.Command, c string) *command.Command {
	for _, cmd := range commands {
		if cmd.Name == c {
			return cmd
		}

		for _, sc := range cmd.SubCommands {
			if sc.Name == c {
				return sc
			}
		}
	}

	return nil
}
