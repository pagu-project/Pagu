package repository

import "github.com/pagu-project/Pagu/internal/entity"

type IValidator interface {
	AddValidator(*entity.Validator) error
	GetValidator(uint) (entity.Validator, error)
}

func (db *DB) AddValidator(v *entity.Validator) error {
	tx := db.Create(v)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *DB) GetValidator(id uint) (entity.Validator, error) {
	var validator entity.Validator
	err := db.Model(&entity.Validator{}).Where("id = ?", id).First(&validator).Error
	if err != nil {
		return entity.Validator{}, err
	}

	return validator, nil
}
