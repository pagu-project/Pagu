package discord

import "github.com/bwmarrin/discordgo"

func helpEmbed(s *discordgo.Session) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "RoboPac Help",
		URL: "https://pactus.org",
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://pactus.org",
			IconURL: s.State.User.AvatarURL(""),
			Name:    s.State.User.Username,
		},
		Description: "RoboPac is a robot that provides support and information about the Pactus Blockchain.",
	}
}
