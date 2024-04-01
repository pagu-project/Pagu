package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID      string `gorm:"primaryKey;unique"` // The ID that defined and assigned on RoboPac.
	Faucets []Faucet

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Faucet struct {
	Address         string
	Amount          int
	TransactionHash string
	UserID          string

	gorm.Model
}
