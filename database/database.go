package database

import (
	"errors"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(path string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, errors.New("can't open database")
	}

	if !db.Migrator().HasTable(&DiscordUser{}) {
		if err := db.AutoMigrate(&DiscordUser{}); err != nil {
			return nil, errors.New("can't auto migrate member table")
		}
	}

	return &DB{
		DB: db,
	}, nil
}

func (db *DB) AddUser(m *DiscordUser) error {
	result := db.Create(m)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) GetUser(dcID string) (*DiscordUser, error) {
	var m DiscordUser

	result := db.First(&m, "discord_id = ?", dcID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &m, nil
}
