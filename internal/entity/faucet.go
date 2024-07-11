package entity

import (
	"time"

	"github.com/pagu-project/Pagu/pkg/amount"
	"gorm.io/gorm"
)

type PhoenixFaucet struct {
	ID              uint `gorm:"primaryKey;unique"`
	UserID          uint `gorm:"size:255"`
	Address         string
	Amount          amount.Amount `gorm:"column:amount"`
	TransactionHash string

	gorm.Model
}

func (f *PhoenixFaucet) TableName() string {
	return "phoenix_faucet"
}

func (f *PhoenixFaucet) ElapsedTime() time.Duration {
	return time.Since(f.CreatedAt)
}
