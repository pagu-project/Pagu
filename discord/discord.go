package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kehiy/RoboPac/engine"
	"github.com/kehiy/RoboPac/log"
	"github.com/pactus-project/pactus/util"
)

type DiscordBot struct {
	Session   *discordgo.Session
	BotEngine engine.IEngine
	GuildID   string
}

func NewDiscordBot(botEngine engine.IEngine, token, guildID string) (*DiscordBot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &DiscordBot{
		Session:   s,
		BotEngine: botEngine,
		GuildID:   guildID,
	}, nil
}

func (db *DiscordBot) Start() {
	log.Info("starting Discord Bot...")

	db.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(db, s, i) // support db *DiscordBot
		}
	})

	err := db.Session.Open()
	if err != nil {
		log.Panic("can't open discord session", "err", err)
	}

	// Updating bot status in real-time by network info.
	log.Info("starting info status")
	go db.UpdateStatusInfo()

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := db.Session.ApplicationCommandCreate(db.Session.State.User.ID, db.GuildID, v)
		if err != nil {
			log.Panic("can not register discord command", "name", v.Name, "err", err)
		}
		registeredCommands[i] = cmd
		log.Info("discord command registered", "name", v.Name)
	}
}

func (db *DiscordBot) UpdateStatusInfo() {
	log.Info("info status started")
	for {
		ns, err := db.BotEngine.NetworkStatus()
		if err != nil {
			continue
		}

		err = db.Session.UpdateStatusComplex(newStatus("validators count", ns.ValidatorsCount))
		if err != nil {
			log.Error("can't set status", "err", err)
			continue
		}

		time.Sleep(time.Second * 5)

		err = db.Session.UpdateStatusComplex(newStatus("height", ns.CurrentBlockHeight))
		if err != nil {
			log.Error("can't set status", "err", err)
			continue
		}

		time.Sleep(time.Second * 5)

		err = db.Session.UpdateStatusComplex(newStatus("circulating supply", util.ChangeToString(ns.CirculatingSupply)))
		if err != nil {
			log.Error("can't set status", "err", err)
			continue
		}

		time.Sleep(time.Second * 5)

		err = db.Session.UpdateStatusComplex(newStatus("total accounts", ns.TotalAccounts))
		if err != nil {
			log.Error("can't set status", "err", err)
			continue
		}

		time.Sleep(time.Second * 5)

		err = db.Session.UpdateStatusComplex(newStatus("total power", util.ChangeToString(ns.TotalNetworkPower)))
		if err != nil {
			log.Error("can't set status", "err", err)
			continue
		}

		time.Sleep(time.Second * 5)
	}
}

func (db *DiscordBot) Stop() {
	log.Info("shutting down Discord Bot...")

	_ = db.Session.Close()
}
