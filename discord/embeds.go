package discord

import "github.com/bwmarrin/discordgo"

func helpEmbed(s *discordgo.Session) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "RoboPac Help",
		URL:   "https://pactus.org",
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://pactus.org",
			IconURL: s.State.User.AvatarURL(""),
			Name:    s.State.User.Username,
		},
		Description: "RoboPac is a robot that provides support and information about the Pactus Blockchain.",
	}
}

func claimEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "claim",
		Description: result,
	}
}

func claimerInfoEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "claimer info",
		Description: result,
	}
}

func nodeInfoEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "node info",
		Description: result,
	}
}

func networkHealthEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "network health",
		Description: result,
	}
}

func networkStatusEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "network status",
		Description: result,
	}
}
