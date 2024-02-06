package discord

import "github.com/bwmarrin/discordgo"

const (
	GREEN  = 0x008000
	RED    = 0xFF0000
	PACTUS = 0x052D5A
)

type UnClaimed struct {
	TestNetAddr     string `json:"testnet_validator_address"`
	DiscordUserName string `json:"discord_username"`
	DiscordID       string `json:"discord_id"`
}

func checkMessage(i *discordgo.InteractionCreate, s *discordgo.Session, guildID, userID string) bool {
	if i.GuildID != guildID || s.State.User.ID == userID {
		return false
	}
	return true
}
