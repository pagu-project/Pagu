package entity

import (
	"github.com/pagu-project/Pagu/pkg/amount"
	"gorm.io/gorm"
)

type ZealyUser struct {
	Amount    amount.Amount `gorm:"column:amount"`
	DiscordID string        `gorm:"column:discord_id"`
	TxHash    string

	gorm.Model
}

func (z *ZealyUser) IsClaimed() bool {
	return z.TxHash != ""
}
