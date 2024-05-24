package user

import (
	"time"

	"github.com/pagu-project/Pagu/internal/repository/faucet"
)

type User struct {
	ID      string `gorm:"primaryKey;unique"` // The ID that defined and assigned on Pagu.
	Faucets []faucet.Faucet

	CreatedAt time.Time
	UpdatedAt time.Time
}
