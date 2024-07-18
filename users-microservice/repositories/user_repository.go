package repositories

import (
	"users/conn"
	"users/models"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Save(user *models.User) (*gorm.DB, error)
}

type UserRepository struct {
	db *gorm.DB
}

func GetUserRepository() UserRepositoryInterface {
	return &UserRepository{db: conn.GetDB()}
}

func (userRepository *UserRepository) Save(user *models.User) (*gorm.DB, error) {
	result := userRepository.db.Save(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return result, nil
}
