package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscordUser struct {
	DiscordID      string `gorm:"primaryKey,unique" json:"discord_id"`
	DepositAddress string `gorm:"unique"            json:"deposit_address"`

	OfferID uuid.UUID `gorm:"index"`
	gorm.Model
}

type Offer struct {
	ID          uuid.UUID `gorm:"primaryKey,unique" json:"id"`
	TotalAmount int64     `json:"total_amount"`
	TotalPrice  int64     `json:"total_price"`
	ChainType   string    `json:"chain_type"`
	Address     string    `json:"address"`

	DiscordUser DiscordUser
	gorm.Model
}
