package api

import (
	"encoding/json"
	"log"
	"transactions/models"
	"transactions/repositories"

	"github.com/nats-io/nats.go"
)

type UserInterface interface {
	CreateUser(user *models.User) (*models.User, error)
	SendUserBalance()
}
type UserApi struct {
	repository repositories.UserRepositoryInterface
	nc         *nats.Conn
}

func GetUserApi(nc *nats.Conn) UserInterface {
	return &UserApi{repository: repositories.GetUserRepository(), nc: nc}
}

func (userApi *UserApi) CreateUser(user *models.User) (*models.User, error) {
	_, err := userApi.repository.Save(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userApi *UserApi) SendUserBalance() {
	userApi.nc.Subscribe("get_balance", func(msg *nats.Msg) {
		email := string(msg.Data)
		response := models.UserBalanceResponse{
			Email: email,
		}
		// Fetch balance for the user
		user, err := userApi.repository.GetByEmail(email)
		if err != nil {
			log.Printf("Failed to get user balance: %v", err)
			response.Error = err.Error()
		} else {
			response.Balance = user.Balance
		}

		// Create response
		balanceResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Failed to marshal balance response: %v", err)
			return
		}

		// Send NATS response
		log.Println("Message sent successfully")
		msg.Respond(balanceResponse)
	})
}
