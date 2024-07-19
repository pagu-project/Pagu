package entity

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Validator struct {
	ID    uint `gorm:"primaryKey;unique"`
	Name  string
	Email string `gorm:"size:255;unique;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime

	gorm.Model
}
