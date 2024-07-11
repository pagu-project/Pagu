package repository

import "github.com/pagu-project/Pagu/internal/entity"

type IZealy interface {
	GetZealyUser(id string) (*entity.ZealyUser, error)
	AddZealyUser(u *entity.ZealyUser) error
	UpdateZealyUser(id string, txHash string) error
	GetAllZealyUser() ([]*entity.ZealyUser, error)
}

func (db *DB) GetZealyUser(id string) (*entity.ZealyUser, error) {
	var u *entity.ZealyUser
	tx := db.Model(&entity.ZealyUser{}).First(&u, "discord_id = ?", id)
	if tx.Error != nil {
		return &entity.ZealyUser{}, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return u, nil
}

func (db *DB) AddZealyUser(u *entity.ZealyUser) error {
	tx := db.Create(u)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *DB) UpdateZealyUser(id, txHash string) error {
	tx := db.Model(&entity.ZealyUser{
		DiscordID: id,
	}).Where("discord_id = ?", id).Update("tx_hash", txHash)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *DB) GetAllZealyUser() ([]*entity.ZealyUser, error) {
	var u []*entity.ZealyUser
	tx := db.Find(&u)
	if tx.Error != nil {
		return nil, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return u, nil
}
