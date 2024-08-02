package repository

import (
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/notification"
)

type INotification interface {
	AddNotification(v *entity.Notification) error
	GetPendingMailNotification() (*entity.Notification, error)
	UpdateNotificationStatus(id uint, status entity.NotificationStatus) error
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

func (db *DB) GetPendingMailNotification() (*entity.Notification, error) {
	var n *entity.Notification
	tx := db.Model(&entity.Notification{}).
		Where("status = ?", entity.NotificationStatusPending).
		Where("type = ?", notification.NotificationTypeMail).
		First(&n)

	if tx.Error != nil {
		return nil, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return n, nil
}

func (db *DB) UpdateNotificationStatus(id uint, status entity.NotificationStatus) error {
	tx := db.Model(&entity.Notification{}).Where("id = ?", id).Update("status", status)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}
