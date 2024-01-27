package discord

import "github.com/bwmarrin/discordgo"

func helpCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.ChannelMessageSendEmbedReply(i.ChannelID, helpEmbed(s), i.Message.Reference())
}
