package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/engine"
	"github.com/pagu-project/Pagu/engine/command"
	"github.com/pagu-project/Pagu/log"
	"github.com/pagu-project/Pagu/utils"
)

type DiscordBot struct {
	Session *discordgo.Session
	engine  *engine.BotEngine
	cfg     config.DiscordBot
}

func NewDiscordBot(botEngine *engine.BotEngine, token string, cfg config.DiscordBot) (*DiscordBot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &DiscordBot{
		Session: s,
		engine:  botEngine,
		cfg:     cfg,
	}, nil
}

func (bot *DiscordBot) Start() error {
	log.Info("starting Discord Bot...")

	err := bot.Session.Open()
	if err != nil {
		return err
	}

	bot.deleteAllCommands()
	return bot.registerCommands()
}

func (bot *DiscordBot) deleteAllCommands() {
	cmdsServer, _ := bot.Session.ApplicationCommands(bot.Session.State.User.ID, bot.cfg.GuildID)
	cmdsGlobal, _ := bot.Session.ApplicationCommands(bot.Session.State.User.ID, "")
	cmds := append(cmdsServer, cmdsGlobal...) //nolint

	for _, cmd := range cmds {
		err := bot.Session.ApplicationCommandDelete(cmd.ApplicationID, cmd.GuildID, cmd.ID)
		if err != nil {
			log.Error("unable to delete command", "error", err, "cmd", cmd.Name)
		} else {
			log.Info("discord command unregistered", "name", cmd.Name)
		}
	}
}

func (bot *DiscordBot) registerCommands() error {
	bot.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		bot.commandHandler(bot, s, i)
	})

	beCmds := bot.engine.Commands()
	for i, beCmd := range beCmds {
		if !beCmd.HasAppId(command.AppIdDiscord) {
			continue
		}

		log.Info("registering new command", "name", beCmd.Name, "desc", beCmd.Desc, "index", i, "object", beCmd)

		discordCmd := discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        beCmd.Name,
			Description: beCmd.Desc,
		}

		if beCmd.HasSubCommand() {
			for _, sCmd := range beCmd.SubCommands {
				if sCmd.Name == "" || sCmd.Desc == "" {
					continue
				}

				log.Info("adding command sub-command", "command", beCmd.Name,
					"sub-command", sCmd.Name, "desc", sCmd.Desc)

				subCmd := &discordgo.ApplicationCommandOption{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        sCmd.Name,
					Description: sCmd.Desc,
				}

				for _, arg := range sCmd.Args {
					if arg.Desc == "" || arg.Name == "" {
						continue
					}

					log.Info("adding sub command argument", "command", beCmd.Name,
						"sub-command", sCmd.Name, "argument", arg.Name, "desc", arg.Desc)

					subCmd.Options = append(subCmd.Options, &discordgo.ApplicationCommandOption{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        arg.Name,
						Description: arg.Desc,
						Required:    !arg.Optional,
					})
				}

				discordCmd.Options = append(discordCmd.Options, subCmd)
			}
		} else {
			for _, arg := range beCmd.Args {
				if arg.Desc == "" || arg.Name == "" {
					continue
				}

				log.Info("adding command argument", "command", beCmd.Name,
					"argument", arg.Name, "desc", arg.Desc)

				discordCmd.Options = append(discordCmd.Options, &discordgo.ApplicationCommandOption{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        arg.Name,
					Description: arg.Desc,
					Required:    !arg.Optional,
				})
			}
		}

		cmd, err := bot.Session.ApplicationCommandCreate(bot.Session.State.User.ID, bot.cfg.GuildID, &discordCmd)
		if err != nil {
			log.Error("can not register discord command", "name", discordCmd.Name, "error", err)
			return err
		}
		log.Info("discord command registered", "name", cmd.Name)
	}

	return nil
}

func (bot *DiscordBot) commandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.GuildID != bot.cfg.GuildID {
		bot.respondErrMsg("Please send messages on server chat", s, i)
		return
	}

	beInput := []string{}

	// Get the application command data
	discordCmd := i.ApplicationCommandData()
	beInput = append(beInput, discordCmd.Name)
	for _, opt := range discordCmd.Options {
		beInput = append(beInput, opt.StringValue())
	}

	res := db.engine.Run(command.AppIdDiscord, i.Member.User.ID, beInput)

	bot.respondResultMsg(res, s, i)
}

func (bot *DiscordBot) respondErrMsg(errStr string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	errorEmbed := &discordgo.MessageEmbed{
		Title:       "Error",
		Description: errStr,
		Color:       RED,
	}
	bot.respondEmbed(errorEmbed, s, i)
}

func (bot *DiscordBot) respondResultMsg(res command.CommandResult, s *discordgo.Session, i *discordgo.InteractionCreate) {
	var resEmbed *discordgo.MessageEmbed
	if res.Successful {
		resEmbed = &discordgo.MessageEmbed{
			Title:       "Successful",
			Description: res.Message,
			Color:       GREEN,
		}
	} else {
		resEmbed = &discordgo.MessageEmbed{
			Title:       "Failed",
			Description: res.Message,
			Color:       YELLOW,
		}
	}

	bot.respondEmbed(resEmbed, s, i)
}

func (db *DiscordBot) respondEmbed(embed *discordgo.MessageEmbed, s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Error("InteractionRespond error:", "error", err)
	}
}

func (db *DiscordBot) UpdateStatusInfo() {
	log.Info("info status started")
	for {
		ns, err := db.engine.NetworkStatus()
		if err != nil {
			continue
		}

		err = db.Session.UpdateStatusComplex(newStatus("validators count", utils.FormatNumber(int64(ns.ValidatorsCount))))
		if err != nil {
			log.Error("can't set status", "err", err)
			continue
		}

		time.Sleep(time.Second * 5)

		err = db.Session.UpdateStatusComplex(newStatus("total accounts", utils.FormatNumber(int64(ns.TotalAccounts))))
		if err != nil {
			log.Error("can't set status", "err", err)
			continue
		}

		time.Sleep(time.Second * 5)

		err = db.Session.UpdateStatusComplex(newStatus("height", utils.FormatNumber(int64(ns.CurrentBlockHeight))))
		if err != nil {
			log.Error("can't set status", "err", err)
			continue
		}

		time.Sleep(time.Second * 5)

		circulatingSupplyAmount := amount.Amount(ns.CirculatingSupply)
		formattedCirculatingSupply := circulatingSupplyAmount.Format(amount.UnitPAC) + " PAC"

		err = db.Session.UpdateStatusComplex(newStatus("circ supply", formattedCirculatingSupply))
		if err != nil {
			log.Error("can't set status", "err", err)
			continue
		}

		time.Sleep(time.Second * 5)

		totalNetworkPowerAmount := amount.Amount(ns.TotalNetworkPower)
		formattedTotalNetworkPower := totalNetworkPowerAmount.Format(amount.UnitPAC) + " PAC"

		err = db.Session.UpdateStatusComplex(newStatus("total power", formattedTotalNetworkPower))
		if err != nil {
			log.Error("can't set status", "err", err)
			continue
		}

		time.Sleep(time.Second * 5)
	}
}

func (db *DiscordBot) Stop() error {
	log.Info("Stopping Discord Bot")

	return db.Session.Close()
}
