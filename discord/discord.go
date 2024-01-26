package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kehiy/RoboPac/engine"
	"github.com/kehiy/RoboPac/log"
)

type DiscordBot struct {
	Session   *discordgo.Session
	BotEngine engine.Engine
	GuildID   string
}

func NewDiscordBot(botEngine engine.Engine, token, guildID string) (*DiscordBot, error) {
	s, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}

	return &DiscordBot{
		Session: s,
		BotEngine: botEngine,
		GuildID: guildID,
	}, nil
}

func (db *DiscordBot) Start() {
	log.Info("starting Discord Bot...")

	err := db.Session.Open()
	if err != nil {
		log.Panic("can't open discord session", "err", err)
	}
}

func (db *DiscordBot) Stop() {
	log.Info("shutting down Discord Bot...")

	_ = db.Session.Close()
}
