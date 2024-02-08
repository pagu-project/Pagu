package discord

import (
	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "help",
		Description: "Help command for RoboPac",
	},
	{
		Name:        "claim",
		Description: "Command to claim the Pactus TestNet rewards coins",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "testnet-addr",
				Description: "Testnet validator address (tpc1p...)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "mainnet-addr",
				Description: "Mainnet validator address (pc1p...)",
				Required:    true,
			},
		},
	},
	{
		Name:        "claimer-info",
		Description: "Get claimer info",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "testnet-addr",
				Description: "Testnet address",
				Required:    true,
			},
		},
	},
	{
		Name:        "node-info",
		Description: "Get node info",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "validator-address",
				Description: "Validator address",
				Required:    true,
			},
		},
	},
	{
		Name:        "reward-calc",
		Description: "calculates how much PAC coins you will earn in a (day/month/year) based on your stake.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "stake",
				Description: "your validator stake amount",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "time",
				Description: "in a day/month/year",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "day",
						Value: "day",
					},
					{
						Name:  "month",
						Value: "month",
					},
					{
						Name:  "year",
						Value: "year",
					},
				},
			},
		},
	},
	{
		Name:        "network-health",
		Description: "network health status",
	},
	{
		Name:        "network-status",
		Description: "status of The Pactus network",
	},
	{
		Name:        "wallet",
		Description: "The RoboPac wallet info",
	},
	{
		Name:        "claim-status",
		Description: "TestNet reward claim status",
	},
	{
		Name:        "get-discount",
		Description: "get twitter campaign discount code",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "twitter-username",
				Description: "your twitter username",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "tweet-link",
				Description: "reposted tweet link",
				Required:    true, 
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "validator-address",
				Description: "your validator address",
				Required:    true, 
			},
	},
}

var commandHandlers = map[string]func(*DiscordBot, *discordgo.Session, *discordgo.InteractionCreate){
	"help":           helpCommandHandler,
	"claim":          claimCommandHandler,
	"claimer-info":   claimerInfoCommandHandler,
	"node-info":      nodeInfoCommandHandler,
	"network-health": networkHealthCommandHandler,
	"network-status": networkStatusCommandHandler,
	"wallet":         walletCommandHandler,
	"claim-status":   claimStatusCommandHandler,
	"reward-calc":    rewardCalcCommandHandler,
	"get-discount": twitterDiscountCampaignCommandHandler,
}
