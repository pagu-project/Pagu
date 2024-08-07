package entity

import (
	"github.com/pagu-project/Pagu/pkg/notification"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type NotificationStatus int

const (
	NotificationStatusPending = iota
	NotificationStatusDone
	NotificationStatusFail
)

type Notification struct {
	ID        uint                          `gorm:"primaryKey;unique"`
	Type      notification.NotificationType `gorm:"size:2"`
	Recipient string                        `gorm:"size:255"`
	Data      datatypes.JSON
	Status    NotificationStatus `gorm:"size:2"`

	gorm.Model
}

type VoucherNotificationData struct {
	Code      string  `json:"code"`
	Amount    float64 `json:"amount"`
	Recipient string  `json:"recipient"`
}
