package entity

import (
	"time"
)

type User struct {
	ID      string `gorm:"primaryKey;unique"` // The ID that defined and assigned on Pagu.
	Faucets []Faucet

	CreatedAt time.Time
	UpdatedAt time.Time
}
