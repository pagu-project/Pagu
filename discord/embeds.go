package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func helpEmbed(s *discordgo.Session) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "RoboPac Help🆘",
		URL:   "https://pactus.org",
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://pactus.org",
			IconURL: s.State.User.AvatarURL(""),
			Name:    s.State.User.Username,
		},
		Description: "RoboPac is a robot that provides support and information about the Pactus Blockchain.\n" +
			"Here is a list of commands supported by RoboPac:\n" +
			"```/claim``` Will help you to claim your test-net rewards on main-net.\n" +
			"```/claimer-info``` Shows you status of your claim reward.\n" +
			"```/node-info``` Shows a node and validator info in network and blockchain.\n" +
			"```/network-status``` Shows a brief info about network.\n" +
			"```/network-health``` Check and shows network health status.\n" +
			"```/bot-wallet``` Shows RoboPac wallet address and balance.\n",
		Color: PACTUS,
	}
}

func claimEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Claim Result💸",
		Description: result,
		Color:       GREEN,
	}
}

func claimerInfoEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Claimer Infoℹ️",
		Description: result,
		Color:       PACTUS,
	}
}

func nodeInfoEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Node Info🛟",
		Description: result,
		Color:       PACTUS,
	}
}

func networkHealthEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string, color int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Network Health🧑‍⚕️",
		Description: result,
		Color:       color,
	}
}

func networkStatusEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Network Status🕸️",
		Description: result,
		Color:       PACTUS,
	}
}

func botWalletEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Bot Wallet🪙",
		Description: result,
		Color:       PACTUS,
	}
}

func claimStatusEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Claim's Status📃",
		Description: result,
		Color:       PACTUS,
	}
}

func rewardCalcEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, result string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Validator reward calculation🧮",
		Description: result,
		Color:       PACTUS,
	}
}

func errorEmbedMessage(reason string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Error",
		Description: fmt.Sprintf("An error occurred: %s", reason),
		Color:       RED,
	}
}
