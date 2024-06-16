package repository

import "github.com/pagu-project/Pagu/internal/entity"

func (db *DB) AddUser(u *entity.User) error {
	tx := db.Create(u)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

/*func (db *DB) GetUser(id string) (*entity.User, error) {
	var u *entity.User
	tx := db.Model(&entity.User{}).Preload("Faucets").First(&u, "id = ?", id)
	if tx.Error != nil {
		return &entity.User{}, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return u, nil
}*/

func (db *DB) HasUser(id string) bool {
	var exists bool

	_ = db.Model(&entity.User{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).
		Error

	return exists
}

func (db *DB) HasUserInApp(appID entity.AppID, callerID string) bool {
	var exists bool

	_ = db.Model(&entity.User{}).
		Select("count(*) > 0").
		Where("application_id = ?", appID).
		Where("caller_id = ?", callerID).
		Find(&exists).
		Error

	return exists
}
