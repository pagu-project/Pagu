package entity

import (
	"github.com/pagu-project/Pagu/pkg/notification"
	"gorm.io/gorm"
)

type NotificationStatus int

const (
	NotificationStatusPending = iota
	NotificationStatusDone
	NotificationStatusFail
)

type Notification struct {
	ID        uint `gorm:"primaryKey;unique"`
	Type      notification.NotificationType
	Recipient string
	Data      []byte `gorm:"Blob"`
	Status    NotificationStatus

	gorm.Model
}
