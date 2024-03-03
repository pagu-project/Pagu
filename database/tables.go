package database

import (
	"gorm.io/gorm"
)

type DiscordUser struct {
	DiscordID      string `gorm:"primaryKey,unique" json:"discord_id"`
	DepositAddress string `gorm:"unique"            json:"deposit_address"`

	OfferID int64 `gorm:"index"`
	gorm.Model
}

type Offer struct {
	ID          int64   `gorm:"primaryKey,autoIncrement" json:"id"`
	TotalAmount int64   `json:"total_amount"`
	TotalPrice  int64   `json:"total_price"`
	UnitPrice   float64 `gorm:"index"                    json:"unit_price"`
	ChainType   string  `json:"chain_type"`
	Address     string  `json:"address"`

	DiscordUser DiscordUser
	gorm.Model
}
