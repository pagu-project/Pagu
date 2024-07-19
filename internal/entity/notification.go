package entity

import "gorm.io/gorm"

type NotificationType int

const (
	NotificationTypeEmail = 0
)

type NotificationStatus int

const (
	NotificationStatusPending = iota
	NotificationStatusInProgress
	NotificationStatusDone
	NotificationStatusFail
)

type Notification struct {
	ID     uint `gorm:"primaryKey;unique"`
	Type   NotificationType
	Email  string
	Status NotificationStatus

	gorm.Model
}
