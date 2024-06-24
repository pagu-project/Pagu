package entity

import (
	"time"

	"gorm.io/gorm"
)

type PhoenixFaucet struct {
	ID              uint `gorm:"primaryKey;unique"`
	UserID          uint `gorm:"size:255"`
	Address         string
	Amount          uint
	TransactionHash string

	gorm.Model
}

func (PhoenixFaucet) TableName() string {
	return "phoenix_faucet"
}

func (f PhoenixFaucet) ElapsedTime() time.Duration {
	return time.Since(f.CreatedAt)
}
