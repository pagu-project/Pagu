package discord

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kehiy/RoboPac/log"
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

	_ = s.InteractionRespond(i.Interaction, response)
}

func claimCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	// Get the application command data
	data := i.ApplicationCommandData()

	// Extract the options
	testnetAddr := data.Options[0].StringValue()
	mainnetAddr := data.Options[1].StringValue()

	log.Info("new claim request", "discordID", i.Member.User.ID, "mainNetAddr", mainnetAddr, "testNetAddr", testnetAddr)

	//! Do we need these two?
	// testnetAddr = strings.TrimPrefix(testnetAddr, "testnet-addr:")
	// mainnetAddr = strings.TrimPrefix(mainnetAddr, "mainnet-addr:")

	command := fmt.Sprintf("claim %s %s %s", i.Member.User.ID, testnetAddr, mainnetAddr)

	result, err := db.BotEngine.Run(command)
	if err != nil {
		errorEmbed := errorEmbedMessage(err.Error())

		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{errorEmbed},
			},
		}

		_ = s.InteractionRespond(i.Interaction, response)

		return
	}

	embed := claimEmbed(s, i, result)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	_ = s.InteractionRespond(i.Interaction, response)
}

func claimerInfoCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	testnetAddr := i.ApplicationCommandData().Options[0].StringValue()
	command := fmt.Sprintf("claimer-info %s", testnetAddr)

	result, err := db.BotEngine.Run(command)
	if err != nil {
		errorEmbed := errorEmbedMessage(err.Error())

		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{errorEmbed},
			},
		}

		_ = s.InteractionRespond(i.Interaction, response)

		return
	}

	embed := claimerInfoEmbed(s, i, result)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	_ = s.InteractionRespond(i.Interaction, response)
}

func nodeInfoCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}
	if i.Member.User.ID != "1081629757793374218" {
		return
	}

	valAddress := i.ApplicationCommandData().Options[0].StringValue()
	command := fmt.Sprintf("node-info %s", valAddress)

	result, err := db.BotEngine.Run(command)
	if err != nil {
		errorEmbed := errorEmbedMessage(err.Error())

		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{errorEmbed},
			},
		}

		_ = s.InteractionRespond(i.Interaction, response)

		return
	}

	embed := nodeInfoEmbed(s, i, result)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	_ = s.InteractionRespond(i.Interaction, response)
}

func networkHealthCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	command := "network-health"

	result, err := db.BotEngine.Run(command)
	if err != nil {
		errorEmbed := errorEmbedMessage(err.Error())

		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{errorEmbed},
			},
		}

		_ = s.InteractionRespond(i.Interaction, response)

		return
	}

	var color int
	if strings.Contains(result, "Healthy") {
		color = GREEN
	} else {
		color = RED
	}

	embed := networkHealthEmbed(s, i, result, color)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	_ = s.InteractionRespond(i.Interaction, response)
}

func networkStatusCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	result, err := db.BotEngine.Run("network")
	if err != nil {
		errorEmbed := errorEmbedMessage(err.Error())

		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{errorEmbed},
			},
		}

		_ = s.InteractionRespond(i.Interaction, response)

		return
	}

	embed := networkStatusEmbed(s, i, result)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	_ = s.InteractionRespond(i.Interaction, response)
}

func botWalletCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	result, _ := db.BotEngine.Run("bot-wallet")

	embed := botWalletEmbed(s, i, result)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	_ = s.InteractionRespond(i.Interaction, response)
}

func claimStatusCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	result, _ := db.BotEngine.Run("claim-status")

	embed := claimStatusEmbed(s, i, result)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	_ = s.InteractionRespond(i.Interaction, response)
}

func unclaimedCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	result := make([]UnClaimed, 300)
	uc := db.BotEngine.Unclaimed()

	for _, claimer := range uc {
		user, err := db.Session.User(claimer.DiscordID)
		if err != nil {
			continue
		}
		result = append(result, UnClaimed{
			DiscordUserName: user.Username,
			DiscordID:       claimer.DiscordID,
		})
	}

	data, err := json.Marshal(result)
	if err != nil {
		errorEmbed := errorEmbedMessage(err.Error())

		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{errorEmbed},
			},
		}

		_ = s.InteractionRespond(i.Interaction, response)

		return
	}

	err = os.WriteFile("unclaimed.json", data, 0o600)
	if err != nil {
		errorEmbed := errorEmbedMessage(err.Error())

		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{errorEmbed},
			},
		}

		_ = s.InteractionRespond(i.Interaction, response)

		return
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "successfully saved on `unclaimed.json`!",
		},
	}
	_ = s.InteractionRespond(i.Interaction, response)
}
