package repository

import (
	"github.com/pagu-project/Pagu/internal/entity"
)

type INotification interface {
	AddNotification(v *entity.Notification) error
}

func (db *DB) AddNotification(v *entity.Notification) error {
	tx := db.Create(v)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}
