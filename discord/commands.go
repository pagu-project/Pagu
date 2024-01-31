package discord

import "github.com/bwmarrin/discordgo"

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "help",
		Description: "Help command for RoboPac",
	},
	{
		Name:        "claim",
		Description: "Command to claim Pactus coins",
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
		Name:        "network-health",
		Description: "information on network health",
	},
	{
		Name:        "network-status",
		Description: "status of Pactus network",
	},
	{
		Name:        "bot-wallet",
		Description: "RoboPac wallet address and balance",
	},
}

var commandHandlers = map[string]func(*DiscordBot, *discordgo.Session, *discordgo.InteractionCreate){
	"help":           helpCommandHandler,
	"claim":          claimCommandHandler,
	"claimer-info":   claimerInfoCommandHandler,
	"node-info":      nodeInfoCommandHandler,
	"network-health": networkHealthCommandHandler,
	"network-status": networkStatusCommandHandler,
	"bot-wallet":     botWalletCommandHandler,
}
