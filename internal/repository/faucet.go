package repository

import (
	"time"

	"github.com/pagu-project/Pagu/internal/entity"
)

func (db *DB) AddFaucet(f *entity.Faucet) error {
	tx := db.Create(f)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *DB) CanGetFaucet(id string) bool {
	var u entity.User
	tx := db.Model(&entity.User{}).Preload("Faucets").First(&u, "id = ?", id)
	if tx.Error != nil {
		return false
	}

	now := time.Now()

	for _, f := range u.Faucets {
		if f.CreatedAt.Year() == now.Year() &&
			f.CreatedAt.Month() == now.Month() &&
			f.CreatedAt.Day() == now.Day() {
			return false
		}
	}

	return true
}
