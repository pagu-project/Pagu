package repository

import (
	"time"

	"github.com/pagu-project/Pagu/internal/entity"
)

func (db *DB) AddFaucet(f *entity.PhoenixFaucet) error {
	tx := db.Create(f)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *DB) CanGetFaucet(user *entity.User) bool {
	var lastFaucet entity.PhoenixFaucet
	tx := db.Model(&entity.PhoenixFaucet{}).
		Last(&lastFaucet, "user_id = ?", user.ID)
	if tx.Error != nil {
		return true
	}

	if lastFaucet.ElapsedTime() > 24*time.Hour {
		return true
	}

	return false
}
