package database

import (
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(path string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, MigrationError{
			Reason: err.Error(),
		}
	}

	if !db.Migrator().HasTable(&User{}) ||
		!db.Migrator().HasTable(&Faucet{}) ||
		!db.Migrator().HasTable(&ZealyUser{}) {
		if err := db.AutoMigrate(
			&User{},
			&Faucet{},
			&ZealyUser{},
		); err != nil {
			return nil, MigrationError{
				Reason: err.Error(),
			}
		}
	}

	return &DB{
		DB: db,
	}, nil
}

func (db *DB) AddUser(u *User) error {
	tx := db.Create(u)
	if tx.Error != nil {
		return WriteError{
			Reason: tx.Error.Error(),
		}
	}

	return nil
}

func (db *DB) GetUser(id string) (*User, error) {
	var u *User
	tx := db.Model(&User{}).Preload("Faucets").First(&u, "id = ?", id)
	if tx.Error != nil {
		return &User{}, ReadError{
			Reason: tx.Error.Error(),
		}
	}

	return u, nil
}

func (db *DB) AddFaucet(f *Faucet) error {
	tx := db.Create(f)
	if tx.Error != nil {
		return WriteError{
			Reason: tx.Error.Error(),
		}
	}

	return nil
}

func (db *DB) HasUser(id string) bool {
	var exists bool

	_ = db.Model(&User{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).
		Error

	return exists
}

func (db *DB) CanGetFaucet(id string) bool {
	var u User
	tx := db.Model(&User{}).Preload("Faucets").First(&u, "id = ?", id)
	if tx.Error != nil {
		return false
	}

	now := time.Now()

	for _, f := range u.Faucets {
		if f.CreatedAt.Year() == now.Year() &&
			f.CreatedAt.Month() == now.Month() &&
			f.CreatedAt.Day() == now.Day() {
			return false
		}
	}

	return true
}

//! Zealy Database

func (db *DB) GetZealyUser(id string) (*ZealyUser, error) {
	var u *ZealyUser
	tx := db.Model(&ZealyUser{}).First(&u, "discord_id = ?", id)
	if tx.Error != nil {
		return &ZealyUser{}, ReadError{
			Reason: tx.Error.Error(),
		}
	}

	return u, nil
}

func (db *DB) AddZealyUser(u *ZealyUser) error {
	tx := db.Create(u)
	if tx.Error != nil {
		return WriteError{
			Reason: tx.Error.Error(),
		}
	}

	return nil
}

func (db *DB) UpdateZealyUser(id string, txHash string) error {
	tx := db.Model(&ZealyUser{
		DiscordID: id,
	}).Where("discord_id = ?", id).Update("tx_hash", txHash).Update("is_claimed", true)
	if tx.Error != nil {
		return WriteError{
			Reason: tx.Error.Error(),
		}
	}

	return nil
}

func (db *DB) GetAllZealyUser() ([]*ZealyUser, error) {
	var u []*ZealyUser
	tx := db.Find(&u)
	if tx.Error != nil {
		return nil, ReadError{
			Reason: tx.Error.Error(),
		}
	}

	return u, nil
}
