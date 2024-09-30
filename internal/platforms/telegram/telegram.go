package telegram

import (
	"context"
	"fmt"
	"strconv"
	"time"

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
	if err := bot.registerCommands(); err != nil {
		return err
	}

	go func() {
		bot.botInstance.Start()
	}()

	log.Info("Starting Telegram Bot...")
	return nil
}

func (bot *Bot) Stop() {
	log.Info("Shutting down Telegram Bot")
	bot.cancel()
	bot.botInstance.Stop()
}

func (bot *Bot) registerCommands() error {
	rows := make([]tele.Row, 0)
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

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

		btn := menu.Data(beCmd.Name, beCmd.Name)
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

			// add back to top menu button
			// backButton := subMenu.Text("Back Main Menu")
			// subRows = append(subRows, subMenu.Row(backButton))

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
	// bot.botInstance.SetCommands(commands)

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

/*
	func (bot *Bot) handleUpdate(ctx tele.Context) error {
		// Check if the message is from a private chat (DM)
		if ctx.Update().Message.Chat.Type != "private" {
			return nil // Ignore non-private messages.
		}

		// Handle the update.
		if ctx.Update().Message != nil {
			// Extract the entire message, including commands.
			fullMessage := ctx.Update().Message.Text

			// Remove the slash prefix.
			fullMessage = strings.TrimPrefix(fullMessage, "/")

			// Split the message into an array.
			messageParts := strings.Split(fullMessage, " ")
			beInput := make(map[string]string)

			for _, t := range messageParts {
				beInput[t] = t
			}

			// Pass the array to the bot engine.
			// callerID := strconv.FormatInt(ctx.Message().Sender.ID, 10)
			callerID := strconv.Itoa(int(ctx.Message().Sender.ID))
			res := bot.engine.Run(entity.AppIDTelegram, callerID, []string{}, beInput)

			// Check if the command execution resulted in an error.
			if res.Error != "" {
				log.Error("Failed to execute command:", res.Error)

				_, err := bot.botInstance.Send(ctx.Message().Sender, "An error occurred while processing your request.", nil)
				if err != nil {
					log.Error("Failed to send error response:", err)
				}
				return nil
			}

			// Send the response back to the user.
			_, err := bot.botInstance.Send(ctx.Message().Sender, res.Message, nil)
			return err
			//if err != nil {
			//	log.Error("Failed to send response:", err)
			//}
			//
			//return nil
		}

		return nil
	}
*/
