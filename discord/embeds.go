package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func helpEmbed(s *discordgo.Session) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "RoboPac Help💸",
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
		Title:       "Claim Result💸",
		Description: result,
	}
}

func claimerInfoEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Claimer Infoℹ️",
		Description: result,
	}
}

func nodeInfoEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Node Info🛟",
		Description: result,
	}
}

func networkHealthEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Network Health🧑‍⚕️",
		Description: result,
	}
}

func networkStatusEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Network Status🕸️",
		Description: result,
	}
}

func botWalletEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Bot Wallet🪙",
		Description: result,
	}
}

func errorEmbedMessage(reason string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Error",
		Description: fmt.Sprintf("An error occurred, please try again! : %s", reason),
		Color:       0xFF0000, // Red color
	}
}
