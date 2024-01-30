package discord

import "github.com/bwmarrin/discordgo"

func checkMessage(i *discordgo.InteractionCreate, s *discordgo.Session, guildID, userID string) bool {
	if i.GuildID != guildID || s.State.User.ID == userID {
		return false
	}
	return true
}
