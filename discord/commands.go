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
				Name:        "discordID",
				Description: "Discord username",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "testnetAddr",
				Description: "Testnet address",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "mainnetAddr",
				Description: "Mainnet address",
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
				Name:        "testnetAddr",
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
				Name:        "valAddress",
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
		Description: "status of pactus network",
	},
}

var commandHandlers = map[string]func(*DiscordBot, *discordgo.Session, *discordgo.InteractionCreate){
	"help":           helpCommandHandler,
	"claim":          claimCommandHandler,
	"claimer-info":   claimerInfoCommandHandler,
	"node-info":      nodeInfoCommandHandler,
	"network-health": networkHealthCommandHandler,
	"network-status": networkStatusCommandHandler,
}
