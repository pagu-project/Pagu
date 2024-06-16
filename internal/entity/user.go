package entity

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint `gorm:"primaryKey;unique"`
	ApplicationID AppID
	CallerID      string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime

	gorm.Model
}
