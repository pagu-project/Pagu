package repository

import (
	"time"

	"github.com/pagu-project/Pagu/internal/repository/faucet"
	"github.com/pagu-project/Pagu/internal/repository/user"
	"github.com/pagu-project/Pagu/internal/repository/zealy"

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

	if !db.Migrator().HasTable(&user.User{}) ||
		!db.Migrator().HasTable(&faucet.Faucet{}) ||
		!db.Migrator().HasTable(&zealy.ZealyUser{}) {
		if err := db.AutoMigrate(
			&user.User{},
			&faucet.Faucet{},
			&zealy.ZealyUser{},
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

func (db *DB) AddUser(u *user.User) error {
	tx := db.Create(u)
	if tx.Error != nil {
		return WriteError{
			Reason: tx.Error.Error(),
		}
	}

	return nil
}

func (db *DB) GetUser(id string) (*user.User, error) {
	var u *user.User
	tx := db.Model(&user.User{}).Preload("Faucets").First(&u, "id = ?", id)
	if tx.Error != nil {
		return &user.User{}, ReadError{
			Reason: tx.Error.Error(),
		}
	}

	return u, nil
}

func (db *DB) AddFaucet(f *faucet.Faucet) error {
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

	_ = db.Model(&user.User{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).
		Error

	return exists
}

func (db *DB) CanGetFaucet(id string) bool {
	var u user.User
	tx := db.Model(&user.User{}).Preload("Faucets").First(&u, "id = ?", id)
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

func (db *DB) GetZealyUser(id string) (*zealy.ZealyUser, error) {
	var u *zealy.ZealyUser
	tx := db.Model(&zealy.ZealyUser{}).First(&u, "discord_id = ?", id)
	if tx.Error != nil {
		return &zealy.ZealyUser{}, ReadError{
			Reason: tx.Error.Error(),
		}
	}

	return u, nil
}

func (db *DB) AddZealyUser(u *zealy.ZealyUser) error {
	tx := db.Create(u)
	if tx.Error != nil {
		return WriteError{
			Reason: tx.Error.Error(),
		}
	}

	return nil
}

func (db *DB) UpdateZealyUser(id string, txHash string) error {
	tx := db.Model(&zealy.ZealyUser{
		DiscordID: id,
	}).Where("discord_id = ?", id).Update("tx_hash", txHash)
	if tx.Error != nil {
		return WriteError{
			Reason: tx.Error.Error(),
		}
	}

	return nil
}

func (db *DB) GetAllZealyUser() ([]*zealy.ZealyUser, error) {
	var u []*zealy.ZealyUser
	tx := db.Find(&u)
	if tx.Error != nil {
		return nil, ReadError{
			Reason: tx.Error.Error(),
		}
	}

	return u, nil
}
