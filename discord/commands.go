package discord

import "github.com/bwmarrin/discordgo"

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
		Name:        "network-health",
		Description: "network health status",
	},
	{
		Name:        "network-status",
		Description: "status of The Pactus network",
	},
	{
		Name:        "bot-wallet",
		Description: "The RoboPac wallet address and balance",
	},
	{
		Name:        "claim-status",
		Description: "TestNet reward claim status",
	},
	{
		Name:        "not-claimed",
		Description: "Admin Only",
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
	"claim-status":   claimStatusCommandHandler,
}
