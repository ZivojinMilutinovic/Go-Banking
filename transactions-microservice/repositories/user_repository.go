package repositories

import (
	"errors"
	"log"
	"transactions/conn"
	"transactions/models"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Save(user *models.User) (*gorm.DB, error)
	GetByEmail(email string) (*models.User, error)
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

func (userRepository *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	filter := map[string]interface{}{"email": email}
	result := userRepository.db.Where(filter).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	} else if result.Error != nil {
		log.Println(result.Error.Error())
		return nil, result.Error
	}

	return &user, nil
}
