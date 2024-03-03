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

	if !db.Migrator().HasTable(&DiscordUser{}) ||
		!db.Migrator().HasTable(&Offer{}) {
		if err := db.AutoMigrate(
			&DiscordUser{},
			&Offer{},
		); err != nil {
			return nil, errors.New("can't auto migrate tables")
		}
	}

	return &DB{
		DB: db,
	}, nil
}

func (db *DB) AddUser(u *DiscordUser) error {
	return db.Create(u).Error
}

func (db *DB) GetUser(dcID string) (*DiscordUser, error) {
	var m DiscordUser

	result := db.First(&m, "discord_id = ?", dcID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &m, nil
}

func (db *DB) HasUser(dcID string) bool {
	var exists bool

	_ = db.Model(&DiscordUser{}).
		Select("count(*) > 0").
		Where("discord_id = ?", dcID).
		Find(&exists).
		Error

	return exists
}

func (db *DB) CreateOffer(o *Offer) error {
	return db.Create(o).Error
}
