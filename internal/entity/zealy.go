package entity

import "gorm.io/gorm"

type ZealyUser struct {
	Amount    int64
	DiscordID string `gorm:"column:discord_id"`
	TxHash    string

	gorm.Model
}

func (z *ZealyUser) IsClaimed() bool {
	return len(z.TxHash) > 0
}
