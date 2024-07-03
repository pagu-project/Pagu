package entity

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Role int

const (
	Admin     Role = 0
	Mod       Role = 1
	BasicUser Role = 2
)

type User struct {
	ID            uint `gorm:"primaryKey;unique"`
	ApplicationID AppID
	CallerID      string
	Role          Role

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime

	gorm.Model
}
