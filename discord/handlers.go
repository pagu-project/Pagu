package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kehiy/RoboPac/log"
)

func (db *DiscordBot) respondErrMsg(cmdErr error, s *discordgo.Session, i *discordgo.InteractionCreate) {
	errorEmbed := errorEmbedMessage(cmdErr.Error())
	db.respondEmbed(errorEmbed, s, i)
}

func (db *DiscordBot) respondEmbed(embed *discordgo.MessageEmbed, s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Error("InteractionRespond error:", "error", err)
	}
}

// TODO: change it to :
// func (db *DiscordBot) helpCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
// ...
// }
func helpCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	embed := helpEmbed(s)
	db.respondEmbed(embed, s, i)
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
		db.respondErrMsg(err, s, i)

		return
	}

	embed := claimEmbed(s, i, result)
	db.respondEmbed(embed, s, i)
}

func claimerInfoCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	testnetAddr := i.ApplicationCommandData().Options[0].StringValue()
	command := fmt.Sprintf("claimer-info %s", testnetAddr)

	result, err := db.BotEngine.Run(command)
	if err != nil {
		db.respondErrMsg(err, s, i)

		return
	}

	embed := claimerInfoEmbed(s, i, result)
	db.respondEmbed(embed, s, i)
}

func nodeInfoCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	valAddress := i.ApplicationCommandData().Options[0].StringValue()
	command := fmt.Sprintf("node-info %s", valAddress)

	result, err := db.BotEngine.Run(command)
	if err != nil {
		db.respondErrMsg(err, s, i)

		return
	}

	embed := nodeInfoEmbed(s, i, result)
	db.respondEmbed(embed, s, i)
}

func networkHealthCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	command := "network-health"

	result, err := db.BotEngine.Run(command)
	if err != nil {
		db.respondErrMsg(err, s, i)

		return
	}

	var color int
	if strings.Contains(result, "Healthy") {
		color = GREEN
	} else {
		color = RED
	}

	embed := networkHealthEmbed(s, i, result, color)
	db.respondEmbed(embed, s, i)
}

func networkStatusCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	result, err := db.BotEngine.Run("network")
	if err != nil {
		db.respondErrMsg(err, s, i)
		return
	}

	embed := networkStatusEmbed(s, i, result)
	db.respondEmbed(embed, s, i)
}

func walletCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	result, _ := db.BotEngine.Run("wallet")

	embed := botWalletEmbed(s, i, result)
	db.respondEmbed(embed, s, i)
}

func claimStatusCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	result, _ := db.BotEngine.Run("claim-status")

	embed := claimStatusEmbed(s, i, result)
	db.respondEmbed(embed, s, i)
}

func rewardCalcCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	stake := i.ApplicationCommandData().Options[0].StringValue()
	time := i.ApplicationCommandData().Options[1].StringValue()

	result, err := db.BotEngine.Run(fmt.Sprintf("calc-reward %v %v", stake, time))
	if err != nil {
		db.respondErrMsg(err, s, i)

		return
	}

	embed := rewardCalcEmbed(s, i, result)
	db.respondEmbed(embed, s, i)
}

func twitterCampaignCommandHandler(db *DiscordBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !checkMessage(i, s, db.GuildID, i.Member.User.ID) {
		return
	}

	twitterID := i.ApplicationCommandData().Options[0].StringValue()
	valAddr := i.ApplicationCommandData().Options[1].StringValue()

	result, err := db.BotEngine.Run(fmt.Sprintf("twitter-campaign %v %v", twitterID, valAddr))
	if err != nil {
		db.respondErrMsg(err, s, i)
		return
	}

	embed := twitterCampaignEmbed(s, i, result)
	db.respondEmbed(embed, s, i)
}
