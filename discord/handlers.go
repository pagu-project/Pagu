package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func helpCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	embed := helpEmbed(s)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		fmt.Println(err)
	}
}

func claimCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if db, s, or i is nil
	if db == nil || s == nil || i == nil {
		fmt.Println("One or more arguments are nil")
		return
	}

	// Get the application command data
	data := i.ApplicationCommandData()
	if data.Options == nil {
		msg := "Command options are missing or invalid."
		_, _ = s.ChannelMessageSend(i.ChannelID, msg)
		return
	}

	// Check if the required number of options are present
	if data.Options == nil || len(data.Options) < 3 {
		msg := "Not enough options provided for the claim command."
		_, _ = s.ChannelMessageSend(i.ChannelID, msg)
		return
	}

	// Extract the options
	var discordId, testnetAddr, mainnetAddr string
	if data.Options[0] != nil {
		discordId = data.Options[0].StringValue()
	}
	if data.Options[1] != nil {
		testnetAddr = data.Options[1].StringValue()
	}
	if data.Options[2] != nil {
		mainnetAddr = data.Options[2].StringValue()
	}

	// Remove the "discord-id:" prefix from the discordId
	discordId = strings.TrimPrefix(discordId, "discord-id:")

	// Remove the "testnet-addr:" and "mainnet-addr:" prefixes from the addresses
	testnetAddr = strings.TrimPrefix(testnetAddr, "testnet-addr:")
	mainnetAddr = strings.TrimPrefix(mainnetAddr, "mainnet-addr:")

	if testnetAddr != "" && mainnetAddr != "" {
		command := fmt.Sprintf("claim %s %s %s", discordId, testnetAddr, mainnetAddr)

		// Check if db or db.BotEngine is nil
		// if db == nil || db.BotEngine == nil {
		// 	msg := "db or bot engine is nil."
		// 	_, _ = s.ChannelMessageSend(i.ChannelID, msg)
		// 	return
		// }
		if db.BotEngine == nil {
			msg := "bot engine is nil."
			_, _ = s.ChannelMessageSend(i.ChannelID, msg)
			return
		}

		result, err := db.BotEngine.Run(command)
		if err != nil {
			msg := fmt.Sprintf("an error occurred while claiming: %v", err)
			_, _ = s.ChannelMessageSend(i.ChannelID, msg)
			return
		}

		embed := claimEmbed(s, i, result)
		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		}

		err = s.InteractionRespond(i.Interaction, response)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func claimerInfoCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	testnetAddr := i.ApplicationCommandData().Options[0].StringValue()
	command := fmt.Sprintf("claimer-info %s", testnetAddr)

	result, err := db.BotEngine.Run(command)
	if err != nil {
		msg := fmt.Sprintf("Wallet address not found, please try again, %v", err)
		_, _ = s.ChannelMessageSend(i.ChannelID, msg)

		return
	}

	embed := claimerInfoEmbed(s, i, result)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		fmt.Println(err)
	}
}

func nodeInfoCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	valAddress := i.ApplicationCommandData().Options[0].StringValue()
	command := fmt.Sprintf("node-info %s", valAddress)

	result, err := db.BotEngine.Run(command)
	if err != nil {
		msg := fmt.Sprintf("an error occcured : %v", err)
		_, _ = s.ChannelMessageSend(i.ChannelID, msg)

		return
	}

	embed := nodeInfoEmbed(s, i, result)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		fmt.Println(err)
	}
}

func networkHealthCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	command := "network-health"

	result, err := db.BotEngine.Run(command)
	if err != nil {
		msg := fmt.Sprintf("an error occured while checking network health: %v", err)
		_, _ = s.ChannelMessageSend(i.ChannelID, msg)

		return
	}

	embed := networkHealthEmbed(s, i, result)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		fmt.Println(err)
	}
}

func networkStatusCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	command := "network"

	result, err := db.BotEngine.Run(command)
	if err != nil {
		msg := fmt.Sprintf("an error occured while checking network status: %v", err)
		_, _ = s.ChannelMessageSend(i.ChannelID, msg)

		return
	}

	embed := networkStatusEmbed(s, i, result)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		fmt.Println(err)
	}
}
