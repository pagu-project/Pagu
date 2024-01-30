package discord

import (
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

func (db *DiscordBot) Stop() {
	log.Info("shutting down Discord Bot...")

	_ = db.Session.Close()
}
