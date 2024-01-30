package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func helpCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Message.Author.ID) {
		return
	}

	_, _ = s.ChannelMessageSendEmbedReply(i.ChannelID, helpEmbed(s), i.Message.Reference())
}

func claimCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Message.Author.ID) {
		return
	}

	discordID := i.Message.Author.ID
	testnetAddr := i.ApplicationCommandData().Options[0].StringValue()
	mainnetAddr := i.ApplicationCommandData().Options[1].StringValue()

	command := fmt.Sprintf("claim %s %s %s", discordID, testnetAddr, mainnetAddr)

	result, err := db.BotEngine.Run(command)
	if err != nil {
		msg := fmt.Sprintf("an error occurred while claiming: %v", err)
		_, _ = s.ChannelMessageSend(i.ChannelID, msg)

		return
	}

	_, _ = s.ChannelMessageSendEmbedReply(i.ChannelID, claimEmbed(s, i, result), i.Message.Reference())
}

func claimerInfoCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Message.Author.ID) {
		return
	}

	testnetAddr := i.ApplicationCommandData().Options[0].StringValue()
	command := fmt.Sprintf("claimer-info %s", testnetAddr)

	result, err := db.BotEngine.Run(command)
	if err != nil {
		msg := fmt.Sprintf("an error occured while checking, please try again: %v", err)
		_, _ = s.ChannelMessageSend(i.ChannelID, msg)

		return
	}

	_, _ = s.ChannelMessageSendEmbedsReply(i.ChannelID, []*discordgo.MessageEmbed{claimerInfoEmbed(s, i, result)}, i.Message.Reference())
}

func nodeInfoCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Message.Author.ID) {
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

	_, _ = s.ChannelMessageSendEmbedsReply(i.ChannelID, []*discordgo.MessageEmbed{nodeInfoEmbed(s, i, result)}, i.Message.Reference())
}

func networkHealthCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Message.Author.ID) {
		return
	}

	command := "network-health"

	result, err := db.BotEngine.Run(command)
	if err != nil {
		msg := fmt.Sprintf("an error occured while checking network health: %v", err)
		_, _ = s.ChannelMessageSend(i.ChannelID, msg)

		return
	}

	_, _ = s.ChannelMessageSendEmbedsReply(i.ChannelID, []*discordgo.MessageEmbed{networkHealthEmbed(s, i, result)}, i.Message.Reference())
}

func networkStatusCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Message.Author.ID) {
		return
	}

	command := "network"

	result, err := db.BotEngine.Run(command)
	if err != nil {
		msg := fmt.Sprintf("an error occured while checking network status: %v", err)
		_, _ = s.ChannelMessageSend(i.ChannelID, msg)

		return
	}

	_, _ = s.ChannelMessageSendEmbedsReply(i.ChannelID, []*discordgo.MessageEmbed{networkStatusEmbed(s, i, result)}, i.Message.Reference())
}
