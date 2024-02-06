package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kehiy/RoboPac/engine"
	"github.com/kehiy/RoboPac/log"
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
	for {
		ns, err := db.BotEngine.NetworkStatus()
		if err != nil {
			continue
		}

		_ = db.Session.UpdateStatusComplex(newStatus("validators count", ns.ValidatorsCount))
		time.Sleep(time.Second * 2)

		_ = db.Session.UpdateStatusComplex(newStatus("height", ns.CurrentBlockHeight))
		time.Sleep(time.Second * 2)

		_ = db.Session.UpdateStatusComplex(newStatus("circulating supply", ns.CirculatingSupply))
		time.Sleep(time.Second * 2)

		_ = db.Session.UpdateStatusComplex(newStatus("total accounts", ns.TotalAccounts))
		time.Sleep(time.Second * 2)

		_ = db.Session.UpdateStatusComplex(newStatus("total power", ns.TotalNetworkPower))
		time.Sleep(time.Second * 2)
	}
}

func (db *DiscordBot) Stop() {
	log.Info("shutting down Discord Bot...")

	_ = db.Session.Close()
}
