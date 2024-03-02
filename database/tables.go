package database

import "gorm.io/gorm"

type User struct {
	gorm.Model

	DiscordID      string `gorm:"unique"           json:"discord_id"`
	Address        string `json:"address"`
	OpenOffers     int    `json:"open_offers"`
	HasOpenPayment bool   `json:"has_open_payment"`
}
