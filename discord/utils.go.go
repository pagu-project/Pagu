package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	GREEN  = 0x008000
	RED    = 0xFF0000
	YELLOW = 0xFFFF00
	PACTUS = 0x052D5A
)

func newStatus(name string, value interface{}) discordgo.UpdateStatusData {
	return discordgo.UpdateStatusData{
		Status: "online",
		Activities: []*discordgo.Activity{
			{
				Type:     discordgo.ActivityTypeCustom,
				Name:     fmt.Sprintf("%s: %v", name, value),
				URL:      "",
				State:    fmt.Sprintf("%s: %v", name, value),
				Details:  fmt.Sprintf("%s: %v", name, value),
				Instance: true,
			},
		},
	}
}
