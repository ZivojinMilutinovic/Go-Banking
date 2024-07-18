package api

import (
	"time"
	"users/actions"
	"users/models"
	"users/repositories"
)

type UserInterface interface {
	CreateUser(userPayload *models.UserPayload) (*models.User, error)
}
type UserApi struct {
	repository repositories.UserRepositoryInterface
}

func GetUserApi() UserInterface {
	return &UserApi{repository: repositories.GetUserRepository()}
}

func (userApi *UserApi) CreateUser(userPayload *models.UserPayload) (*models.User, error) {
	user := models.User{
		Email:     userPayload.Email,
		CreatedAt: time.Now(),
	}

	_, err := userApi.repository.Save(&user)

	if err != nil {
		return nil, err
	}
	actions.PushEventToKafka("user-created-topic", &user)

	return &user, nil
}
