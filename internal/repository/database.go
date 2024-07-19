package repository

import (
	"github.com/pagu-project/Pagu/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(path string) (Database, error) {
	db, err := gorm.Open(mysql.Open(path), &gorm.Config{})
	if err != nil {
		return nil, ConnectionError{
			Message: err.Error(),
		}
	}

	if !db.Migrator().HasTable(&entity.User{}) ||
		!db.Migrator().HasTable(&entity.PhoenixFaucet{}) ||
		!db.Migrator().HasTable(&entity.Voucher{}) ||
		!db.Migrator().HasTable(&entity.ZealyUser{}) ||
		!db.Migrator().HasTable(&entity.Validator{}) ||
		!db.Migrator().HasTable(&entity.Notification{}) {
		if err := db.AutoMigrate(
			&entity.User{},
			&entity.PhoenixFaucet{},
			&entity.ZealyUser{},
			&entity.Voucher{},
			&entity.Validator{},
			&entity.Notification{},
		); err != nil {
			return nil, MigrationError{
				Message: err.Error(),
			}
		}
	}

	return &DB{
		DB: db,
	}, nil
}
